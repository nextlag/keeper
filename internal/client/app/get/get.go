package get

import (
	"fmt"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
)

var App = config.LoadConfig().App.Name
var Get = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
	Long:  "Get different types of resources like login, card, note or binary.",
	Example: fmt.Sprintf(`
# Get a login
%s get login -i login_id

# Get a card
%s get card -i card_id

# Get a note
%s get note -i note_id

# Get a binary
%s get note -i binary_id
	`, App, App, App, App),
}

func init() {
	Get.AddCommand(Card)
	Get.AddCommand(Login)
	Get.AddCommand(Note)
	Get.AddCommand(Binary)
}
