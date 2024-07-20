package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	config "github.com/nextlag/keeper/config/client"
	"github.com/nextlag/keeper/internal/client/app/build"
)

var rootCmd = &cobra.Command{
	Use:   config.LoadConfig().App.Name,
	Short: "App for storing private data",
	Long:  `User can save cards, note and logins`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		build.PrintBuildInfo()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
