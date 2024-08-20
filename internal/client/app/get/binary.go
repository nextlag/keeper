package get

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var Binary = &cobra.Command{
	Use:   "binary",
	Short: "Get user file by id",
	Long: fmt.Sprintf(`
This command get user binary info and encode it for path
Usage: %s get binary -i binary_id -f some_file`, App),

	Run: func(cmd *cobra.Command, args []string) {
		if getBinaryID == "" || filePath == "" {
			fmt.Println("Error: Both 'id' and 'file' flags are required")
			if err := cmd.Help(); err != nil {
				return
			}
			os.Exit(1)
		}

		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			color.Red("Authentication required. Error: %v", err)
			return
		}

		usecase.GetClientUseCase().GetBinary(userPassword, getBinaryID, filePath)
	},
}

var (
	getBinaryID string
	filePath    string
)

func init() {
	Binary.Flags().StringVarP(&getBinaryID, "id", "i", "", "Binary id")
	Binary.Flags().StringVarP(&filePath, "file", "f", "", "User file")

	if err := Binary.MarkFlagRequired("id"); err != nil {
		color.Red("%v", err)
		return
	}
	if err := Binary.MarkFlagRequired("file"); err != nil {
		color.Red("%v", err)
		return
	}
}
