package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	config "github.com/nextlag/keeper/config/server"
	"github.com/nextlag/keeper/internal/server/controller/http/v1/mocks"
	"github.com/nextlag/keeper/pkg/logger/l"
)

func TestController_HealthCheck(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name             string
		mockHealthCheck  func(*mocks.MockUseCase)
		expectedStatus   int
		expectedResponse any
	}{
		{
			name: "successful health check",
			mockHealthCheck: func(mockUc *mocks.MockUseCase) {
				mockUc.EXPECT().HealthCheck().Return(nil).Times(1)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"status":"connected"}`,
		},
		{
			name: "failed health check",
			mockHealthCheck: func(mockUc *mocks.MockUseCase) {
				mockUc.EXPECT().HealthCheck().Return(errors.New("service unavailable")).Times(1)
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "Application not available\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUc := mocks.NewMockUseCase(ctrl)
			tt.mockHealthCheck(mockUc)
			controller := &Controller{
				uc:  mockUc,
				cfg: cfg,
				log: l.NewLogger(cfg),
			}

			req, err := http.NewRequest("GET", healthCheck, nil)
			assert.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.HealthCheck)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}
