package cards

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var GetCard = &cobra.Command{
	Use:   "getcard",
	Short: "Show user card by id",
	Long: `
This command add card
Usage: getcard -i \"card_id\" 
Flags:
  -i, --id string Card id
  `,
	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			return
		}
		usecase.GetClientUseCase().ShowCard(userPassword, getCardID)
	},
}

var getCardID string

func init() {
	GetCard.Flags().StringVarP(&getCardID, "id", "i", "", "Card id")
	if err := GetCard.MarkFlagRequired("id"); err != nil {
		log.Fatal(err)
	}
}
