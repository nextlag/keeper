package v1

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	config "github.com/nextlag/keeper/config/server"
	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/server/mw/gzip"
	"github.com/nextlag/keeper/internal/server/mw/request"
	"github.com/nextlag/keeper/pkg/logger/l"
)

// UseCase defines the interface for the business logic operations used by the controller.
type UseCase interface {
	HealthCheck() error
	AddLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error
	SignUpUser(ctx context.Context, email, password string) (entity.User, error)
	SignInUser(ctx context.Context, email, password string) (entity.JWT, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (entity.JWT, error)
	GetDomainName() string
	CheckAccessToken(ctx context.Context, accessToken string) (entity.User, error)
}

// Controller represents the HTTP handlers controller.
type Controller struct {
	uc  UseCase // The UseCase used to perform business logic operations.
	cfg *config.Config
	log *l.Logger
}

// NewController creates a new instance of the controller.
func NewController(uc UseCase, cfg *config.Config, log *l.Logger) *Controller {
	return &Controller{uc: uc, cfg: cfg, log: log}
}

// NewServer creates a new HTTP server with specified routes and middleware.
func (c *Controller) NewServer(handler *chi.Mux) *http.Server {
	handler.Use(middleware.RequestID)
	handler.Use(request.MwRequest(c.log))
	handler.Use(gzip.MwGzip())
	handler.Use(middleware.Recoverer)

	// Routes for health check
	handler.Get("/ping", c.HealthCheck) // Endpoint for health check

	// Routes for authentication
	handler.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.SignUpUser)       // Endpoint for user registration
		r.Post("/login", c.SignInUser)          // Endpoint for user authentication
		r.Get("/refresh", c.RefreshAccessToken) // Endpoint for refresh token
		r.Get("/logout", c.LogoutUser)          // Endpoint for logout user
	})

	// Routes for user operations.
	handler.Route("/user", func(r chi.Router) {
		r.Use(c.MwAuth())            // Middleware for user authentication
		r.Get("/me", c.UserInfo)     // Endpoint for retrieving current user information
		r.Post("/login", c.AddLogin) // Endpoint for adding login credentials for the current user
	})

	return &http.Server{
		Addr:    c.cfg.Network.Host,
		Handler: handler,
	}
}
