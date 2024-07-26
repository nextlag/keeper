package del

import "github.com/spf13/cobra"

var Del = &cobra.Command{
	Use:   "del",
	Short: "Del resources",
	Long:  `Del different types of resources like cards or logins.`,
	Example: `
# Get a card
del card -i <uuid>

# Get a login
del login -i <uuid>
	`,
}

func init() {
	Del.AddCommand(Card)
	Del.AddCommand(Login)
}
