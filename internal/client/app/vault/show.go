package vault

import (
	"fmt"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
	"github.com/nextlag/keeper/internal/client/usecase"
)

var ShowVault = &cobra.Command{
	Use:   "show",
	Short: "Show user vault",
	Long: fmt.Sprintf(`
This command show user vault
Usage: %s show -o a|c|l|n|b
Flags:
  -o, --option string     Option for listing (default "a")
	a - all
	c - cards
	l - logins
	n - notes
	b - binaries
  `, config.Load().App.Name),

	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			return
		}
		usecase.GetClientUseCase().ShowVault(userPassword, showVaultOption)
	},
}

var showVaultOption string

func init() {
	ShowVault.Flags().StringVarP(&showVaultOption, "option", "o", "a", "Option for listing")
}
