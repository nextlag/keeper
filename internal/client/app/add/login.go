package add

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
	"github.com/nextlag/keeper/internal/entity"
	utils "github.com/nextlag/keeper/internal/utils/client"
)

var loginForAdditing entity.Login

var Login = &cobra.Command{
	Use:   "login",
	Short: "Add login",
	Long: fmt.Sprintf(`This command adds a login for a site.
Example:
  %s add login -t "Login Title" -l "user@example.com" -s "password" -u "https://example.com" --meta '[{"name":"meta","value":"value"}]'
	`, App),
	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			log.Fatal(err)
		}
		usecase.GetClientUseCase().AddLogin(userPassword, &loginForAdditing)
	},
}

func init() {
	Login.Flags().StringVarP(&loginForAdditing.Name, "title", "t", "", "Login title")
	Login.Flags().StringVarP(&loginForAdditing.Login, "login", "l", "", "Site login")
	Login.Flags().StringVarP(&loginForAdditing.Password, "secret", "s", "", "Site password|secret")
	Login.Flags().StringVarP(&loginForAdditing.URI, "uri", "u", "", "Site endpoint")
	Login.Flags().Var(&utils.JSONFlag{Target: &loginForAdditing.Meta}, "meta", `Add meta fields for entity`)

	if err := Login.MarkFlagRequired("title"); err != nil {
		log.Fatal(err)
	}
}
