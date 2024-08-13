package get

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var Note = &cobra.Command{
	Use:   "note",
	Short: "Show user note by id",
	Long: fmt.Sprintf(`
This command show user note
Usage: %s get note -i <note_id>`, App),
	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			color.Red("Authentication required. Error: %v", err)
			return
		}
		usecase.GetClientUseCase().ShowNote(userPassword, getNoteID)
	},
}

var getNoteID string

func init() {
	Note.Flags().StringVarP(&getNoteID, "id", "i", "", "Note id")
	if err := Note.MarkFlagRequired("id"); err != nil {
		color.Red("%v", err)
		return
	}
}
