package del

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var Card = &cobra.Command{
	Use:   "card",
	Short: "Delete user card by id",
	Long: fmt.Sprintf(`
This command remove card
Usage: %s del card -i <card_id>`, App),

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
		log.Println(err)
		return
	}
}
