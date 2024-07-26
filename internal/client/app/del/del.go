package del

import (
	"fmt"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
)

var app = config.LoadConfig().App.Name
var Del = &cobra.Command{
	Use:   "del",
	Short: "Del resources",
	Long:  `Del different types of resources like cards or logins.`,
	Example: fmt.Sprintf(`
# Get a card
%s del card -i <uuid>

# Get a login
%s del login -i <uuid>

# Get a note
%s del note -i <uuid>

# Get a binary
%s del binary -i <uuid>
	`, app, app, app, app),
}

func init() {
	Del.AddCommand(Card)
	Del.AddCommand(Login)
	Del.AddCommand(Note)
}
