package vault

import (
	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var ShowVault = &cobra.Command{
	Use:   "showvault",
	Short: "Show user vault",
	Long: `
This command show user vault
Usage: showvault -o \"a|c|l|n\" 
Flags:
  -o, --option string     Option for listing (default "a")
	a - all
	c - cards
	l - logins
	n - notes
	b - bynaries
  `,
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
