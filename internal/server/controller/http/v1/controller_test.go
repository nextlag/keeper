package v1

import (
	"testing"

	"github.com/golang/mock/gomock"

	config "github.com/nextlag/keeper/config/server"
	"github.com/nextlag/keeper/internal/server/controller/http/v1/mocks"
	"github.com/nextlag/keeper/pkg/logger/l"
)

const (
	// Health check
	healthCheckPath = "/api/v1/ping"

	// Authentication
	authRegisterPath = "/api/v1/auth/register"
	authLoginPath    = "/api/v1/auth/login"
	authRefreshPath  = "/api/v1/auth/refresh"
	authLogoutPath   = "/api/v1/auth/logout"

	// User
	userInfoPath           = "/api/v1/user/me"
	userLoginsPath         = "/api/v1/user/logins"
	userLoginsIDPath       = "/api/v1/user/logins/{id}"
	userCardsPath          = "/api/v1/user/cards"
	userCardsIDPath        = "/api/v1/user/cards/{id}"
	userNotesPath          = "/api/v1/user/notes"
	userNotesIDPath        = "/api/v1/user/notes/{id}"
	userBinaryAddPath      = "/api/v1/user/binary"
	userBinaryAddMetaPath  = "/api/v1/user/binary/{id}/meta"
	userBinaryGetPath      = "/api/v1/user/binary"
	userBinaryDownloadPath = "/api/v1/user/binary/{id}"
	userBinaryDeletePath   = "/api/v1/user/binary/{id}"
)

func loadTest(t *testing.T) (*Controller, *mocks.MockUseCase, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUseCase := mocks.NewMockUseCase(ctrl)
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	log := l.NewLogger(cfg)

	c := &Controller{
		uc:  mockUseCase,
		cfg: cfg,
		log: log,
	}

	return c, mockUseCase, ctrl
}
