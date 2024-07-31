package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
	"github.com/nextlag/keeper/internal/client/usecase"
	"github.com/nextlag/keeper/internal/entity"
)

var RegisterUser = &cobra.Command{
	Use:   "register",
	Short: "User registration in the service",
	Long: fmt.Sprintf(`This command registers a new user.
Usage: %s register <user_email> <user_password>`, config.Load().App.Name),
	Args: cobra.MinimumNArgs(RequiredUserArgs),
	Run: func(cmd *cobra.Command, args []string) {
		account := entity.User{
			Email:    args[0],
			Password: args[1],
		}
		usecase.GetClientUseCase().Register(&account)
	},
}
