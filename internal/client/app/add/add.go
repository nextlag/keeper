package add

import (
	"github.com/spf13/cobra"
)

var Add = &cobra.Command{
	Use:   "add",
	Short: "Add resources",
	Long:  `Add different types of resources like cards or logins.`,
	Example: `
# Add a card
add card -t "Card Title" -n "1234 5678 9012 3456" -o "Card Owner" -b "VISA" -c "123" -m "12" -y "2025" --meta '[{"name":"meta","value":"value"}]'

# Add a login
add login -t "Login Title" -l "user@example.com" -s "password" -u "https://example.com" --meta '[{"name":"meta","value":"value"}]'
	`,
}

func init() {
	Add.AddCommand(Card)
	Add.AddCommand(Login)
}
