package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

var registerUserCmd = &cobra.Command{ // cobra style guide
	Use:   "register",
	Short: "Register user to the service",
	Long: `
This command register new user.
Usage: client register user_login user_password`,
	Args: cobra.MinimumNArgs(RequiredUserArgs),
	Run: func(cmd *cobra.Command, args []string) {
		login := loginUser{
			login:    args[0],
			password: args[1],
		}
		fmt.Println(login)
		// TODO: add register logic
	},
}

func init() {
	rootCmd.AddCommand(registerUserCmd)
}
