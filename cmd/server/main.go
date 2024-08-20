package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	config "github.com/nextlag/keeper/config/server"
	_ "github.com/nextlag/keeper/docs"
	"github.com/nextlag/keeper/internal/server/app"
)

// @title Keeper Server
// @version 1.0.0
// @description keeper project
// @contact.name Nexbug
// @contact.url https://github.com/nextlag
// @contact.email nextbug@ya.ru
// @host localhost:8080
// @BasePath /api/v1
// Main func.
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	srv, err := app.NewApp(cfg)
	if srv.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
