package app

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
	"github.com/nextlag/keeper/internal/client/app/add"
	"github.com/nextlag/keeper/internal/client/app/auth"
	"github.com/nextlag/keeper/internal/client/app/build"
	"github.com/nextlag/keeper/internal/client/app/del"
	"github.com/nextlag/keeper/internal/client/app/get"
	"github.com/nextlag/keeper/internal/client/app/storage"
	"github.com/nextlag/keeper/internal/client/app/vault"
	"github.com/nextlag/keeper/internal/client/usecase"
	"github.com/nextlag/keeper/internal/client/usecase/api"
	"github.com/nextlag/keeper/internal/client/usecase/repo"
)

var (
	cfg *config.Config

	rootCmd = &cobra.Command{
		Use:   config.Load().App.Name,
		Short: "App for storing private data",
		Long:  `User can save cards, notes, and logins.`,
		Run: func(cmd *cobra.Command, args []string) {
			build.PrintBuildInfo()
		},
	}
)

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// init initializes the application and adds subcommands to the root command.
func init() {
	cobra.OnInitialize(initApp)

	commands := []*cobra.Command{
		storage.InitLocalStorage, // Command to initialize local storage.
		storage.SyncUserData,     // Command to sync user data with the server.

		auth.LoginUser,    // Command to log in a user.
		auth.RegisterUser, // Command to register a new user.
		auth.LogoutUser,   // Command to log out a user.

		add.Add,    // Command to add new entities.
		add.Login,  // Command to add a new login.
		add.Card,   // Command to add a new card.
		add.Note,   // Command to add a new note.
		add.Binary, // Command to add a new binary file.

		get.Get,    // Command to retrieve entities.
		get.Login,  // Command to retrieve logins.
		get.Card,   // Command to retrieve cards.
		get.Note,   // Command to retrieve notes.
		get.Binary, // Command to retrieve binary files.

		del.Del,    // Command to delete entities.
		del.Login,  // Command to delete a login.
		del.Card,   // Command to delete a card.
		del.Note,   // Command to delete a note.
		del.Binary, // Command to delete a binary file.

		vault.ShowVault, // Command to display the vault.
	}

	rootCmd.AddCommand(commands...)
}

// initApp initializes the application configuration and use case.
// Sets up the necessary directories and configurations based on the provided settings.
func initApp() {
	cfg = config.Load()
	uc := usecase.GetClientUseCase()
	clientOpts := []usecase.OptsUseCase{
		usecase.SetAPI(api.New(cfg.Server.ServerURL)),
		usecase.SetConfig(cfg),
		usecase.SetRepo(repo.New(cfg.SQLite.DSN)),
	}

	for _, opt := range clientOpts {
		opt(uc)
	}

	if _, err := os.Stat(cfg.FilesStorage.ServerLocation); os.IsNotExist(err) {
		err = os.MkdirAll(cfg.FilesStorage.ServerLocation, os.ModePerm)
		if err != nil {
			log.Fatalf("App.Init - os.MkdirAll - %v", err)
		}
	}
}
