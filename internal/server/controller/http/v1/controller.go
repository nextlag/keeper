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

//go:generate mockgen -destination=mocks/mocks.go -package=mocks github.com/nextlag/keeper/internal/server/controller/http/v1 UseCase
type UseCase interface {
	HealthCheck() error
	CheckAccessToken(ctx context.Context, accessToken string) (entity.User, error)
	AddLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error
	SignUpUser(ctx context.Context, email, password string) (entity.User, error)
	SignInUser(ctx context.Context, email, password string) (entity.JWT, error)
}

type Controller struct {
	uc  UseCase
	cfg *config.Config
	log *l.Logger
}

func NewController(uc UseCase, cfg *config.Config, log *l.Logger) *Controller {
	return &Controller{uc: uc, cfg: cfg, log: log}
}

func (c *Controller) NewServer(handler *chi.Mux) *http.Server {
	handler.Use(middleware.RequestID)
	handler.Use(request.MwRequest(c.log))
	handler.Use(gzip.MwGzip())
	handler.Use(middleware.Recoverer)

	handler.Get("/ping", c.HealthCheck) // healthCheck

	handler.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.SignUpUser)
		r.Post("/login", c.SignInUser)
	})

	handler.Route("/user", func(r chi.Router) {
		r.Use(c.MwAuth())        // middleware for checking authorization
		r.Get("/me", c.UserInfo) // getting information about the current user

		r.Post("/login", c.AddLogin)
	})

	return &http.Server{
		Addr:    c.cfg.Network.Host,
		Handler: handler,
	}
}
