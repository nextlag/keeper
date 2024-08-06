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
	healthCheck = "/api/v1/ping"

	// Authentication
	userAuth     = "/api/v1/user/auth"
	authRegister = "/api/v1/auth/register"
	authLogin    = "/api/v1/auth/login"
	authRefresh  = "/api/v1/auth/refresh"
	authLogout   = "/api/v1/auth/logout"

	// User
	userInfo          = "/api/v1/user/me"
	userLogins        = "/api/v1/user/logins"
	userCards         = "/api/v1/user/cards"
	userNotes         = "/api/v1/user/notes"
	userBinaryAddMeta = "/user/binary/{id}/meta"
	userBinary        = "/api/v1/user/binary"
)

func loadTest(t *testing.T) (*Controller, *mocks.MockUseCase, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUseCase := mocks.NewMockUseCase(ctrl)
	cfg, err := config.Load()
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
