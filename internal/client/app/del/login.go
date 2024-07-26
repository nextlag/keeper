package del

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
	"github.com/nextlag/keeper/internal/client/usecase"
)

var Login = &cobra.Command{
	Use:   "login",
	Short: "Delete user login by id",
	Long: fmt.Sprintf(`
This command remove login
Usage: %s del login -i <login_id>
  `, config.LoadConfig().App.Name),

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
