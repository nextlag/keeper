package get

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"

	"github.com/nextlag/keeper/internal/client/app/add"
)

var Binary = &cobra.Command{
	Use:   "binary",
	Short: "Get user file by id",
	Long: fmt.Sprintf(`
This command get user binary info and encode it for path
Usage: %s get binary -i binary_id -f some_file`, App),

	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
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
		log.Fatal(err)
	}
	if err := add.Binary.MarkFlagRequired("file"); err != nil {
		log.Fatal(err)
	}
}
