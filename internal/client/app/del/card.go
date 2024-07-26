package del

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var Card = &cobra.Command{
	Use:   "card",
	Short: "Delete user card by id",
	Long: `
This command remove card
Usage: card -i <card_id>
Flags:
  -i, --id string card id
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
	Card.Flags().StringVarP(&delCardID, "id", "i", "", "Card id")
	if err := Card.MarkFlagRequired("id"); err != nil {
		log.Fatal(err)
	}
}
