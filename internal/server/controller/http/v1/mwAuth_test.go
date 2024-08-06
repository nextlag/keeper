package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/nextlag/keeper/internal/entity"
)

func TestMwAuth(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	token := "token"
	expectedUser := entity.User{ID: uuid.New()}

	tests := []struct {
		name           string
		accessToken    string
		expectedStatus int
		mockReturn     interface{}
	}{
		{
			name:           "successful authorization with Bearer token",
			accessToken:    token,
			expectedStatus: http.StatusOK,
			mockReturn:     expectedUser,
		},
		{
			name:           "unauthorized - missing token",
			accessToken:    "",
			expectedStatus: http.StatusUnauthorized,
			mockReturn:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockReturn != nil {
				mockUseCase.EXPECT().
					CheckAccessToken(gomock.Any(), tt.accessToken).
					Return(tt.mockReturn, nil).
					Times(1)
			}

			req := httptest.NewRequest(http.MethodGet, userAuth, nil)
			if tt.accessToken != "" {
				req.Header.Set("Authorization", "Bearer "+tt.accessToken)
			}
			if tt.accessToken == token {
				req.AddCookie(&http.Cookie{Name: "access_token", Value: tt.accessToken})
			}
			rr := httptest.NewRecorder()

			handler := c.MwAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
