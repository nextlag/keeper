package usecase

import (
	config "github.com/nextlag/keeper/config/server"
	"github.com/nextlag/keeper/internal/server/usecase/repository"
	c "github.com/nextlag/keeper/pkg/cache"
	"github.com/nextlag/keeper/pkg/logger/l"
)

// UseCase represents the use case layer for business logic operations.
type UseCase struct {
	repo  repository.Repository // Repository interface for data access
	cfg   *config.Config
	cache c.Cacher // Cache interface for caching operations
	log   *l.Logger
}

// New creates a new instance of UseCase with provided dependencies.
func New(r repository.Repository, cfg *config.Config, cache c.Cacher, log *l.Logger) *UseCase {
	return &UseCase{
		repo:  r,
		cfg:   cfg,
		cache: cache,
		log:   log,
	}
}

// HealthCheck performs a health check on the repository/database.
func (uc *UseCase) HealthCheck() error {
	return uc.repo.DBHealthCheck()
}

// GetDomainName - retrieves the domain name configured for the application.
// This domain name is used for setting cookies and other domain-specific configurations.
func (uc *UseCase) GetDomainName() string {
	return uc.cfg.Security.Domain
}
