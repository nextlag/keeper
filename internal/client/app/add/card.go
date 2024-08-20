package add

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
	"github.com/nextlag/keeper/internal/entity"
	utils "github.com/nextlag/keeper/internal/utils/client"
)

var cardForAdditing entity.Card

var Card = &cobra.Command{
	Use:   "card",
	Short: "Add card",
	Long: fmt.Sprintf(`
This command adds a card.
Example: 
  %s add card -t "Card Title" -n "1234 5678 9012 3456" -o "Card Owner" -b "VISA" -c "123" -m "12" -y "2025" \
  --meta '[{"name":"meta","value":"value"}]'`, App),

	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			color.Red("Authentication required. Error: %v", err)
			return
		}
		usecase.GetClientUseCase().AddCard(userPassword, &cardForAdditing)
	},
}

func init() {
	Card.Flags().StringVarP(&cardForAdditing.Name, "title", "t", "", "Card title")
	Card.Flags().StringVarP(&cardForAdditing.Number, "number", "n", "", "Card number")
	Card.Flags().StringVarP(&cardForAdditing.CardHolderName, "owner", "o", "", "Cardholder name")
	Card.Flags().StringVarP(&cardForAdditing.Brand, "brand", "b", "", "Card brand")
	Card.Flags().StringVarP(&cardForAdditing.SecurityCode, "code", "c", "", "CVV/CVC")
	Card.Flags().StringVarP(&cardForAdditing.ExpirationMonth, "month", "m", "", "Card expiration month")
	Card.Flags().StringVarP(&cardForAdditing.ExpirationYear, "year", "y", "", "Card expiration year")
	Card.Flags().Var(&utils.JSONFlag{Target: &cardForAdditing.Meta}, "meta", `Add meta fields for entity`)

	if err := Card.MarkFlagRequired("title"); err != nil {
		color.Red("%v", err)
		return
	}
	if err := Card.MarkFlagRequired("number"); err != nil {
		color.Red("%v", err)
		return
	}
}
