package get

import (
	"fmt"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
)

var app = config.LoadConfig().App.Name
var Get = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
	Long:  "Get different types of resources like cards or logins.",
	Example: fmt.Sprintf(`
# Get a card
%s get card -i <uuid>

# Get a login
%s get login -i <uuid>

# Get a note
%s get note -i <uuid>
	`, app, app, app),
}

func init() {
	Get.AddCommand(Card)
	Get.AddCommand(Login)
	Get.AddCommand(Note)
}
