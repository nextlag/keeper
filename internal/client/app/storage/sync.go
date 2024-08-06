package storage

import (
	"fmt"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
	"github.com/nextlag/keeper/internal/client/usecase"
)

var SyncUserData = &cobra.Command{
	Use:   "sync",
	Short: "Sync user`s data",
	Long: fmt.Sprintf(`This command update users private data from server 
Usage: %s sync`, config.Load().App.Name),
	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			return
		}
		usecase.GetClientUseCase().Sync(userPassword)
	},
}
