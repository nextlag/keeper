package get

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var Login = &cobra.Command{
	Use:   "login",
	Short: "Show user login by id",
	Long: fmt.Sprintf(`
This command login
Usage: %s get login -i <login_id>`, App),

	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			color.Red("Authentication required. Error: %v", err)
			return
		}
		usecase.GetClientUseCase().ShowLogin(userPassword, getLoginID)
	},
}

var getLoginID string

func init() {
	Login.Flags().StringVarP(&getLoginID, "id", "i", "", "Card id")

	if err := Login.MarkFlagRequired("id"); err != nil {
		color.Red("%v", err)
		return
	}
}
