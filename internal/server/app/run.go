package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutdown = time.Second * 15

func (a *App) Run(ctx context.Context) {
	defer a.repo.ShutDown()
	srv := a.ctrl.NewServer(a.router)
	a.router.Mount("/", srv.Handler)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigint
		a.log.Info("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdown)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal("HTTP server Shutdown:", err)
		}
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("HTTP server ListenAndServe:", err)
	}

	<-ctx.Done()

	a.log.Info("Server Shutdown gracefully")
}
