package add

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
	"github.com/nextlag/keeper/internal/entity"
	utils "github.com/nextlag/keeper/internal/utils/client"
)

var Note = &cobra.Command{
	Use:   "note",
	Short: "add note",
	Long: fmt.Sprintf(`
This command add user note
Example: 
 %s add note -t name -n content --meta '[{"name":"meta","value":"value"}]'`, App),

	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			return
		}
		usecase.GetClientUseCase().AddNote(userPassword, &noteForAdditing)
	},
}

var (
	noteForAdditing entity.SecretNote
)

func init() {
	Note.Flags().StringVarP(&noteForAdditing.Name, "title", "t", "", "Login title")
	Note.Flags().StringVarP(&noteForAdditing.Note, "note", "n", "", "User note")
	Note.Flags().Var(&utils.JSONFlag{Target: &noteForAdditing.Meta}, "meta", `Add meta fields for entity`)

	if err := Note.MarkFlagRequired("title"); err != nil {
		log.Fatal(err)
	}
}
