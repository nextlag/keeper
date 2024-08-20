package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils/errs"
)

func TestUserInfo(t *testing.T) {
	c, _, ctrl := loadTest(t)
	defer ctrl.Finish()

	expectedUser := entity.User{
		Email: "example@mail.ru",
	}

	tests := []struct {
		name           string
		expectedStatus int
		expectedEmail  string
		expectedError  error
	}{
		{
			name:           "successful user info retrieval",
			expectedStatus: http.StatusOK,
			expectedEmail:  "example@mail.ru",
			expectedError:  nil,
		},
		{
			name:           "error retrieving user from context",
			expectedStatus: http.StatusInternalServerError,
			expectedEmail:  "",
			expectedError:  errs.ErrUnexpectedError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, userInfo, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			if tt.expectedError == nil {
				req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))
			} else {
				req = req.WithContext(context.WithValue(req.Context(), currentUserKey, nil))
			}

			rr := httptest.NewRecorder()
			http.HandlerFunc(c.UserInfo).ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				assert.Equal(t, tt.expectedEmail, response["email"])
			}
		})
	}
}
