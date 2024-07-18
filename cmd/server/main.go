package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	config "github.com/nextlag/keeper/config/server"
	"github.com/nextlag/keeper/internal/server/app"
)

func main() {
	cfg, err := config.NewConfig()
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
