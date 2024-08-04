package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nextlag/keeper/internal/entity"
)

func TestAddLogin(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	validUUID := uuid.New()
	expectedUser := entity.User{ID: validUUID}

	login := entity.Login{
		Name:     "testuser",
		Password: "password123",
	}

	tests := []struct {
		name           string
		mockReturn     error
		expectedStatus int
		expectedBody   string
		reqBody        interface{}
	}{
		{
			name:           "successful add login",
			mockReturn:     nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   `{"uuid":"00000000-0000-0000-0000-000000000000","name":"testuser","login":"","password":"password123","uri":"","meta":null}` + "\n",
			reqBody:        login,
		},
		{
			name:           "error from use case",
			mockReturn:     errors.New("add login failed"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "add login failed\n",
			reqBody:        login,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockReturn != nil {
				mockUseCase.EXPECT().
					AddLogin(gomock.Any(), gomock.Any(), validUUID).
					Return(tt.mockReturn).Times(1)
			} else {
				mockUseCase.EXPECT().
					AddLogin(gomock.Any(), gomock.Any(), validUUID).
					Return(nil).Times(1)
			}

			var reqBody []byte
			var err error

			if str, ok := tt.reqBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.reqBody)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, userLogins, bytes.NewBuffer(reqBody))
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.AddLogin).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestGetLogins(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	validUUID := uuid.New()
	expectedUser := entity.User{ID: validUUID}

	logins := []entity.Login{
		{
			Name:     "example",
			Login:    "example",
			Password: "12345",
			URI:      "https://example.com",
		},
	}

	tests := []struct {
		name           string
		mockReturn     interface{}
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful get logins",
			mockReturn:     logins,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"uuid":"00000000-0000-0000-0000-000000000000","name":"example","login":"example","password":"12345","uri":"https://example.com","meta":null}]` + "\n",
		},
		{
			name:           "error from use case",
			mockReturn:     nil,
			mockError:      errors.New("error fetching logins"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "error fetching logins\n",
		},
		{
			name:           "no logins found",
			mockReturn:     []entity.Login{},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase.EXPECT().
				GetLogins(gomock.Any(), expectedUser).
				Return(tt.mockReturn, tt.mockError).Times(1)

			req := httptest.NewRequest(http.MethodGet, userLogins, nil)
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.GetLogins).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestUpdateLogin(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	validUUID := uuid.New()
	expectedUser := entity.User{ID: validUUID}

	login := entity.Login{
		ID:       validUUID,
		Name:     "example",
		Login:    "example",
		Password: "12345",
		URI:      "https://example.com",
	}

	tests := []struct {
		name           string
		mockReturn     error
		expectedStatus int
		expectedBody   string
		loginID        string
		reqBody        interface{}
		expectCall     bool
	}{
		{
			name:           "successful update login",
			mockReturn:     nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   "Update accepted",
			loginID:        validUUID.String(),
			reqBody:        login,
			expectCall:     true,
		},
		{
			name:           "error from use case",
			mockReturn:     errors.New("update failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "update failed\n",
			loginID:        validUUID.String(),
			reqBody:        login,
			expectCall:     true,
		},
		{
			name:           "invalid UUID in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 0\n",
			loginID:        "",
			reqBody:        login,
			expectCall:     false,
		},
		{
			name:           "non-UUID string in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 10\n",
			loginID:        "123a45test",
			reqBody:        login,
			expectCall:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectCall {
				mockUseCase.EXPECT().
					UpdateLogin(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.mockReturn).Times(1)
			}

			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, userLogins+tt.loginID, bytes.NewBuffer(reqBody))
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.loginID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.UpdateLogin).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestDelLogin(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	validUUID := uuid.New()
	expectedUser := entity.User{ID: validUUID}

	tests := []struct {
		name           string
		mockReturn     error
		expectedStatus int
		expectedBody   string
		loginID        string
		expectCall     bool
	}{
		{
			name:           "successful delete login",
			mockReturn:     nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   "Delete accepted",
			loginID:        validUUID.String(),
			expectCall:     true,
		},
		{
			name:           "error from use case",
			mockReturn:     errors.New("delete failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "delete failed\n",
			loginID:        validUUID.String(),
			expectCall:     true,
		},
		{
			name:           "invalid UUID in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 0\n",
			loginID:        "",
			expectCall:     false,
		},
		{
			name:           "non-UUID string in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 10\n",
			loginID:        "123a45test",
			expectCall:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectCall {
				mockUseCase.EXPECT().
					DelLogin(gomock.Any(), validUUID, expectedUser.ID).
					Return(tt.mockReturn).Times(1)
			}

			req := httptest.NewRequest(http.MethodDelete, userLogins+tt.loginID, nil)
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.loginID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.DelLogin).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
