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
//
//go:generate mockgen -destination=mocks/mocks.go -package=mocks github.com/nextlag/keeper/internal/server/controller/http/v1 UseCase
type UseCase interface {
	HealthCheck() error
	SignUpUser(ctx context.Context, email, password string) (entity.User, error)
	SignInUser(ctx context.Context, email, password string) (entity.JWT, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (entity.JWT, error)
	GetDomainName() string
	CheckAccessToken(ctx context.Context, accessToken string) (entity.User, error)

	GetLogins(ctx context.Context, user entity.User) ([]entity.Login, error)
	AddLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error
	DelLogin(ctx context.Context, loginID, userID uuid.UUID) error
	UpdateLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error

	GetCards(ctx context.Context, user entity.User) ([]entity.Card, error)
	AddCard(ctx context.Context, card *entity.Card, userID uuid.UUID) error
	DelCard(ctx context.Context, cardUUID, userID uuid.UUID) error
	UpdateCard(ctx context.Context, card *entity.Card, userID uuid.UUID) error

	GetNotes(ctx context.Context, user entity.User) ([]entity.SecretNote, error)
	AddNote(ctx context.Context, note *entity.SecretNote, userID uuid.UUID) error
	DelNote(ctx context.Context, noteID, userID uuid.UUID) error
	UpdateNote(ctx context.Context, note *entity.SecretNote, userID uuid.UUID) error
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

	handler.Get("/ping", c.HealthCheck) // Endpoint for health check

	// Routes for authentication
	handler.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.SignUpUser)
		r.Post("/login", c.SignInUser)
		r.Get("/refresh", c.RefreshAccessToken)
		r.Get("/logout", c.LogoutUser)
	})

	// Routes for user operations.
	handler.Route("/user", func(r chi.Router) {
		r.Use(c.MwAuth())        // Middleware for user authentication
		r.Get("/me", c.UserInfo) // Endpoint for retrieving current user information

		r.Post("/logins", c.AddLogin)
		r.Get("/logins", c.GetLogins)
		r.Delete("/logins/{id}", c.DelLogin)
		r.Patch("/logins/{id}", c.UpdateLogin)

		r.Get("/cards", c.GetCards)
		r.Post("/cards", c.AddCard)
		r.Delete("/cards/{id}", c.DelCard)
		r.Patch("/cards/{id}", c.UpdateCard)

		r.Get("/notes", c.GetNotes)
		r.Post("/notes", c.AddNote)
		r.Delete("/notes/{id}", c.DelNote)
		r.Patch("/notes/{id}", c.UpdateNote)

	})

	return &http.Server{
		Addr:    c.cfg.Network.Host,
		Handler: handler,
	}
}
