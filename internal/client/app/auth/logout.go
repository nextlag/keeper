package auth

import (
	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var LogoutUser = &cobra.Command{
	Use:   "logout",
	Short: "Logout user",
	Long: `
This command drops users tokens`,
	Run: func(cmd *cobra.Command, args []string) {
		usecase.GetClientUseCase().Logout()
	},
}
