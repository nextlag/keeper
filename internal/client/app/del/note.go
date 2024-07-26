package del

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
	"github.com/nextlag/keeper/internal/client/usecase"
)

var Note = &cobra.Command{
	Use:   "note",
	Short: "Delete user note by id",
	Long: fmt.Sprintf(`
This command remove note
Usage: %s del note -i <note_id> 
  `, config.LoadConfig().App.Name),

	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			return
		}
		usecase.GetClientUseCase().DelNote(userPassword, delNoteID)
	},
}

var delNoteID string

func init() {
	Note.Flags().StringVarP(&delNoteID, "id", "i", "", "Note id")

	if err := Note.MarkFlagRequired("id"); err != nil {
		log.Fatal(err)
	}
}
