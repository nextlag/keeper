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

func TestAddNote(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	validUUID := uuid.New()
	expectedUser := entity.User{ID: validUUID}

	note := entity.SecretNote{
		Name: "Test Note",
		Note: "This is a test note.",
		Meta: nil,
	}

	tests := []struct {
		name           string
		mockReturn     error
		expectedStatus int
		expectedBody   string
		reqBody        interface{}
	}{
		{
			name:           "successful add note",
			mockReturn:     nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   `{"uuid":"00000000-0000-0000-0000-000000000000","name":"Test Note","note":"This is a test note.","meta":null}` + "\n",
			reqBody:        note,
		},
		{
			name:           "error from use case",
			mockReturn:     errors.New("failed to add note"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "failed to add note\n",
			reqBody:        note,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatus == http.StatusAccepted || tt.expectedStatus == http.StatusBadRequest {
				mockUseCase.EXPECT().
					AddNote(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.mockReturn).
					Times(1)
			} else {
				mockUseCase.EXPECT().
					AddNote(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			}

			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, userNotes, bytes.NewBuffer(reqBody))
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.AddNote).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestGetNotes(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	validUUID := uuid.New()
	expectedUser := entity.User{ID: validUUID}

	notes := []entity.SecretNote{
		{
			Name: "Note 1",
			Note: "This is the first note.",
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
			name:           "successful get notes",
			mockReturn:     notes,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"uuid":"00000000-0000-0000-0000-000000000000","name":"Note 1","note":"This is the first note.","meta":null}]` + "\n",
		},
		{
			name:           "no notes found",
			mockReturn:     []entity.SecretNote{},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:           "error from use case",
			mockReturn:     nil,
			mockError:      errors.New("failed to get notes"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to get notes\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockError != nil {
				mockUseCase.EXPECT().
					GetNotes(gomock.Any(), gomock.Any()).
					Return(tt.mockReturn, tt.mockError).
					Times(1)
			} else {
				mockUseCase.EXPECT().
					GetNotes(gomock.Any(), gomock.Any()).
					Return(tt.mockReturn, nil).
					Times(1)
			}

			req := httptest.NewRequest(http.MethodGet, userNotes, nil)
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.GetNotes).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestUpdateNote(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	validUUID := uuid.New()
	expectedUser := entity.User{ID: validUUID}

	note := entity.SecretNote{
		ID:   validUUID,
		Name: "Updated Note",
		Note: "This is an updated note.",
	}

	tests := []struct {
		name           string
		mockReturn     error
		expectedStatus int
		expectedBody   string
		noteID         string
		reqBody        interface{}
		expectCall     bool
	}{
		{
			name:           "successful update note",
			mockReturn:     nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   "Update accepted",
			noteID:         validUUID.String(),
			reqBody:        note,
			expectCall:     true,
		},
		{
			name:           "error from use case",
			mockReturn:     errors.New("update failed"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "update failed\n",
			noteID:         validUUID.String(),
			reqBody:        note,
			expectCall:     true,
		},
		{
			name:           "invalid UUID in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 0\n",
			noteID:         "",
			reqBody:        note,
			expectCall:     false,
		},
		{
			name:           "non-UUID string in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 10\n",
			noteID:         "123a45test",
			reqBody:        note,
			expectCall:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectCall {
				mockUseCase.EXPECT().
					UpdateNote(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.mockReturn).Times(1)
			}

			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, userNotes+tt.noteID, bytes.NewBuffer(reqBody))
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.noteID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.UpdateNote).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestDelNote(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	validUUID := uuid.New()
	expectedUser := entity.User{ID: validUUID}

	tests := []struct {
		name           string
		mockReturn     error
		expectedStatus int
		expectedBody   string
		noteID         string
		expectCall     bool
	}{
		{
			name:           "successful delete note",
			mockReturn:     nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   "Delete accepted",
			noteID:         validUUID.String(),
			expectCall:     true,
		},
		{
			name:           "error from use case",
			mockReturn:     errors.New("delete failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "delete failed\n",
			noteID:         validUUID.String(),
			expectCall:     true,
		},
		{
			name:           "invalid UUID in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 0\n",
			noteID:         "",
			expectCall:     false,
		},
		{
			name:           "non-UUID string in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 10\n",
			noteID:         "123a45test",
			expectCall:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectCall {
				mockUseCase.EXPECT().
					DelNote(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.mockReturn).Times(1)
			}

			req := httptest.NewRequest(http.MethodDelete, userNotes+tt.noteID, nil)
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.noteID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.DelNote).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
