package app

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
	"github.com/nextlag/keeper/internal/client/app/auth"
	"github.com/nextlag/keeper/internal/client/app/build"
	"github.com/nextlag/keeper/internal/client/app/vault"
	"github.com/nextlag/keeper/internal/client/usecase"
	"github.com/nextlag/keeper/internal/client/usecase/api"
	"github.com/nextlag/keeper/internal/client/usecase/repo"
)

var (
	cfg *config.Config

	rootCmd = &cobra.Command{
		Use:   config.LoadConfig().App.Name,
		Short: "App for storing private data",
		Long:  `User can save cards, note and logins`,
		Run: func(cmd *cobra.Command, args []string) {
			build.PrintBuildInfo()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initApp)
	commands := []*cobra.Command{
		vault.RegisterInitLocalStorage,

		auth.LoginUserCmd,
		auth.RegisterUserCmd,
		auth.LogoutUser,
	}

	rootCmd.AddCommand(commands...)
}

func initApp() {
	cfg = config.LoadConfig()

	log.Println(cfg.Network.Host)

	uc := usecase.GetClientUseCase()
	clientOpts := []usecase.UseCaseOpts{
		usecase.SetAPI(api.New(cfg.Network.Host)),
		usecase.SetConfig(cfg),
		usecase.SetRepo(repo.New(cfg.SQLite.DSN)),
	}

	for _, opt := range clientOpts {
		opt(uc)
	}

	if _, err := os.Stat(cfg.FilesStorage.Location); os.IsNotExist(err) {
		err = os.MkdirAll(cfg.FilesStorage.Location, os.ModePerm)
		if err != nil {
			log.Fatalf("App.Init - os.MkdirAll - %v", err)
		}
	}
}
