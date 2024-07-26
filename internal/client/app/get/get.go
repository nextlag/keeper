package get

import "github.com/spf13/cobra"

var Get = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
	Long:  `Get different types of resources like cards or logins.`,
	Example: `
# Get a card
get card -i <uuid>

# Get a login
get login -i <uuid>
	`,
}

func init() {
	Get.AddCommand(Card)
	Get.AddCommand(Login)
}
