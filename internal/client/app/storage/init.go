package storage

import (
	"fmt"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
	"github.com/nextlag/keeper/internal/client/usecase"
)

var InitLocalStorage = &cobra.Command{
	Use:   "init",
	Short: "Initialize local storage",
	Long: fmt.Sprintf(`This command will register sqlite db to store personal data...
Usage: %s init`, config.Load().App.Name),
	Run: func(cmd *cobra.Command, args []string) {
		usecase.GetClientUseCase().InitDB()
	},
}
