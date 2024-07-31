package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
	"github.com/nextlag/keeper/internal/client/usecase"
	"github.com/nextlag/keeper/internal/entity"
)

var RequiredUserArgs = 2

var LoginUser = &cobra.Command{
	Use:   "login",
	Short: "Login user to the service",
	Long: fmt.Sprintf(`This is the user login command.
Usage: %s login <user_email> <user_password>`, config.Load().App.Name),
	Args: cobra.MinimumNArgs(RequiredUserArgs),
	Run: func(cmd *cobra.Command, args []string) {
		account := entity.User{
			Email:    args[0],
			Password: args[1],
		}
		usecase.GetClientUseCase().Logout()
		usecase.GetClientUseCase().Login(&account)
	},
}
