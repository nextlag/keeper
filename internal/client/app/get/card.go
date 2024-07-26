package get

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var Card = &cobra.Command{
	Use:   "card",
	Short: "Show user card by id",
	Long: `
This command add card
Usage: card -i <card_id> 
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
	Card.Flags().StringVarP(&getCardID, "id", "i", "", "Card id")
	if err := Card.MarkFlagRequired("id"); err != nil {
		log.Fatal(err)
	}
}
