package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils/errs"
)

func TestSignUpUser(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		payload        *loginPayload
		expectedStatus int
		expectedEmail  string
		expectedError  error
	}{
		{
			name: "successful signup",
			payload: &loginPayload{
				Email:    "test@example.com",
				Password: "password",
			},
			expectedStatus: http.StatusCreated,
			expectedEmail:  "test@example.com",
			expectedError:  nil,
		},
		{
			name: "invalid email",
			payload: &loginPayload{
				Email:    "invalid",
				Password: "password",
			},
			expectedStatus: http.StatusBadRequest,
			expectedEmail:  "",
			expectedError:  errs.ErrWrongEmail,
		},
		{
			name: "email already exists",
			payload: &loginPayload{
				Email:    "exists@example.com",
				Password: "password",
			},
			expectedStatus: http.StatusBadRequest,
			expectedEmail:  "",
			expectedError:  errs.ErrEmailAlreadyExists,
		},
		{
			name: "internal error",
			payload: &loginPayload{
				Email:    "test@example.com",
				Password: "password",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedEmail:  "",
			expectedError:  errors.New("internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError != nil {
				mockUseCase.EXPECT().SignUpUser(
					gomock.Any(),
					tt.payload.Email,
					tt.payload.Password,
				).Return(entity.User{}, tt.expectedError)
			} else {
				mockUseCase.EXPECT().SignUpUser(
					gomock.Any(),
					tt.payload.Email,
					tt.payload.Password,
				).Return(entity.User{Email: tt.payload.Email}, nil)
			}

			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, authRegister, bytes.NewReader(body))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.SignUpUser).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				assert.Equal(t, tt.expectedEmail, response["email"])
			}
		})
	}
}

func TestSignInUser(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		payload        *loginPayload
		expectedStatus int
		expectedBody   string
		expectedError  error
	}{
		{
			name: "successful signin",
			payload: &loginPayload{
				Email:    "test@example.com",
				Password: "password",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"access_token":"access-token","refresh_token":"refresh-token"}`,
			expectedError:  nil,
		},
		{
			name: "wrong credentials",
			payload: &loginPayload{
				Email:    "wrong@example.com",
				Password: "wrong password",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"wrong credentials"}`,
			expectedError:  errs.ErrWrongCredentials,
		},
		{
			name: "internal error",
			payload: &loginPayload{
				Email:    "test@example.com",
				Password: "password",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal error"}`,
			expectedError:  errors.New("internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError != nil {
				mockUseCase.EXPECT().SignInUser(
					gomock.Any(),
					tt.payload.Email,
					tt.payload.Password,
				).Return(entity.JWT{}, tt.expectedError)
			} else {
				jwtToken := entity.JWT{
					AccessToken:        "access-token",
					RefreshToken:       "refresh-token",
					AccessTokenMaxAge:  3600,
					RefreshTokenMaxAge: 7200,
					Domain:             "example.com",
				}
				mockUseCase.EXPECT().SignInUser(
					gomock.Any(),
					tt.payload.Email,
					tt.payload.Password,
				).Return(jwtToken, nil)
			}

			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, authLogin, bytes.NewReader(body))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.SignInUser).ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				cookies := rr.Result().Cookies()
				assert.Len(t, cookies, 3)
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				assert.Contains(t, response, "access_token")
				assert.Contains(t, response, "refresh_token")
			}
		})
	}
}

func TestRefreshAccessToken(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		cookieValue    string
		mockReturn     entity.JWT
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful refresh",
			cookieValue:    "valid_refresh_token",
			mockReturn:     entity.JWT{AccessToken: "new_access_token", RefreshToken: "new_refresh_token", AccessTokenMaxAge: 3600, Domain: "example.com"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"access_token":"new_access_token","refresh_token":"new_refresh_token"}`,
		},
		{
			name:           "error refreshing token",
			cookieValue:    "valid_refresh_token",
			mockReturn:     entity.JWT{},
			mockError:      errors.New("unexpected error"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"unexpected error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, authRefresh, nil)
			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tt.cookieValue})
			}

			rr := httptest.NewRecorder()

			if tt.mockError != nil {
				mockUseCase.EXPECT().RefreshAccessToken(gomock.Any(), tt.cookieValue).Return(entity.JWT{}, tt.mockError)
			} else {
				mockUseCase.EXPECT().RefreshAccessToken(gomock.Any(), tt.cookieValue).Return(tt.mockReturn, nil)
			}

			http.HandlerFunc(c.RefreshAccessToken).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Remove trailing newline characters for comparison
			actualBody := strings.TrimSpace(rr.Body.String())
			expectedBody := strings.TrimSpace(tt.expectedBody)

			assert.Equal(t, expectedBody, actualBody)
		})
	}
}

func TestLogoutUser(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	mockUseCase.EXPECT().GetDomainName().Return("example.com")

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success logout",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"logout success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, authLogout, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			http.HandlerFunc(c.LogoutUser).ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}
