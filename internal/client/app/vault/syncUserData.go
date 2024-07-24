package vault

import (
	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var SyncUserData = &cobra.Command{
	Use:   "sync",
	Short: "Sync user`s data",
	Long: `
This command update users private data from server
Usage: keeper sync"`,
	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			return
		}
		usecase.GetClientUseCase().Sync(userPassword)
	},
}
