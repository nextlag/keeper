package get

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var Card = &cobra.Command{
	Use:   "card",
	Short: "Show user card by id",
	Long: fmt.Sprintf(`
This command add card
Usage: %s get card -i card_id`, App),

	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			color.Red("Authentication required. Error: %v", err)
			return
		}
		usecase.GetClientUseCase().ShowCard(userPassword, getCardID)
	},
}

var getCardID string

func init() {
	Card.Flags().StringVarP(&getCardID, "id", "i", "", "Card id")
	if err := Card.MarkFlagRequired("id"); err != nil {
		color.Red("%v", err)
		return
	}
}
