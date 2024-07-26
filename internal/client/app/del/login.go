package del

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var Login = &cobra.Command{
	Use:   "login",
	Short: "Delete user login by id",
	Long: `
This command remove login
Usage: login -i <login_id>
Flags:
  -i, --id string login id
  `,
	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			return
		}
		usecase.GetClientUseCase().DelLogin(userPassword, delLoginID)
	},
}

var delLoginID string

func init() {
	Login.Flags().StringVarP(&delLoginID, "id", "i", "", "Card id")
	if err := Login.MarkFlagRequired("id"); err != nil {
		log.Fatal(err)
	}
}
