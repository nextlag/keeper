package app

import (
	"errors"

	"github.com/go-chi/chi/v5"

	config "github.com/nextlag/keeper/config/server"
	v1 "github.com/nextlag/keeper/internal/server/controller/http/v1"
	"github.com/nextlag/keeper/internal/server/usecase"
	"github.com/nextlag/keeper/internal/server/usecase/repository"
	"github.com/nextlag/keeper/pkg/cache"
	"github.com/nextlag/keeper/pkg/logger/l"
)

// App represents the main application with its components.
type App struct {
	cfg    *config.Config   // Configuration for the application.
	log    *l.Logger        // Logger for application logging.
	router *chi.Mux         // HTTP router for routing requests.
	uc     *usecase.UseCase // Use case layer handling business logic.
	repo   *repository.Repo // Repository layer for data access.
	ctrl   *v1.Controller   // HTTP controller for handling requests.
}

// NewApp initializes a new App instance with the given configuration.
// It sets up the logger, repository, use case, and controller.
// Returns the initialized App or an error if the configuration is invalid.
func NewApp(cfg *config.Config) (*App, error) {
	if cfg == nil {
		return nil, errors.New("failed to initialize: config is invalid")
	}

	log := l.NewLogger(cfg)
	log.Debug("Configuration initialized", "config", cfg.Network)

	repo := repository.New(cfg.PG.DSN, log)
	repo.Migrate()

	uc := usecase.New(repo, cfg, cache.New(cfg.Cache.DefaultExpiration, cfg.Cache.CleanupInterval), log)

	r := chi.NewRouter()
	ctrl := v1.NewController(uc, cfg, log)

	return &App{
		repo:   repo,
		uc:     uc,
		router: r,
		ctrl:   ctrl,
		cfg:    cfg,
		log:    log,
	}, nil
}
