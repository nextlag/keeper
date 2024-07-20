package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/nextlag/keeper/pkg/logger/l"
)

const shutdown = time.Second * 15

// Run starts the HTTP server and handles graceful shutdown
func (a *App) Run(ctx context.Context) {
	defer a.repo.ShutDown()

	srv := &http.Server{
		Addr:    a.cfg.Network.Host,
		Handler: a.ctrl.NewServer(a.router).Handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("HTTP server ListenAndServe:", l.ErrAttr(err))
		}
	}()

	<-ctx.Done()

	a.log.Debug("Starting graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdown)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		a.log.Error("HTTP server Shutdown:", l.ErrAttr(err))
	}

	a.log.Info("Server shutdown gracefully")
}
