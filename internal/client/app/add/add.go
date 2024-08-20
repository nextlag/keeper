package add

import (
	"fmt"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
)

var App = config.Load().App.Name
var Add = &cobra.Command{
	Use:   "add",
	Short: "Add resources",
	Long:  `Add different types of resources like login, card, note or binary.`,
	Example: fmt.Sprintf(`
# Add a login
%s add login -t "Login Title" -l "user@example.com" -s "password" -u "https://example.com" --meta '[{"name":"meta","value":"value"}]'

# Add a card
%s add card -t "Card Title" -n "1234 5678 9012 3456" -o "Card Owner" -b "VISA" -c "123" -m "12" -y "2025" --meta '[{"name":"meta","value":"value"}]'

# Add a note
 %s add note -t "Name" -n "Content" --meta '[{"name":"meta","value":"value"}]'

# Add a binary
 %s add binary -t "name" -f "file_location" --meta '[{"name":"meta","value":"value"}]'
	`, App, App, App, App),
}

func init() {
	Add.AddCommand(Card)
	Add.AddCommand(Login)
	Add.AddCommand(Note)
	Add.AddCommand(Binary)
}
