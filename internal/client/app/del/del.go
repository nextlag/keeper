package del

import (
	"fmt"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
)

var App = config.LoadConfig().App.Name
var Del = &cobra.Command{
	Use:   "del",
	Short: "Del resources",
	Long:  `Del different types of resources like login, card, note or binary.`,
	Example: fmt.Sprintf(`
# Get a card
%s del card -i card_id

# Get a login
%s del login -i login_id

# Get a note
%s del note -i note_id

# Get a binary
%s del binary -i binary_id
	`, App, App, App, App),
}

func init() {
	Del.AddCommand(Card)
	Del.AddCommand(Login)
	Del.AddCommand(Note)
	Del.AddCommand(Binary)
}
