package cards

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var DelCard = &cobra.Command{
	Use:   "delcard",
	Short: "Delete user card by id",
	Long: `
This command remove card
Usage: delcard -i \"card_id\" 
Flags:
  -i, --id string Card id
  `,
	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			return
		}
		usecase.GetClientUseCase().DelCard(userPassword, delCardID)
	},
}

var delCardID string

func init() {
	DelCard.Flags().StringVarP(&delCardID, "id", "i", "", "Card id")
	if err := DelCard.MarkFlagRequired("id"); err != nil {
		log.Fatal(err)
	}
}
