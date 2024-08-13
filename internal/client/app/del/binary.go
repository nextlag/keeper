package del

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
)

var Binary = &cobra.Command{
	Use:   "binary",
	Short: "Delete user file by id",
	Long: fmt.Sprintf(`
This command remove file
Usage: %s del binary -i binary_id`, App),

	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			return
		}
		usecase.GetClientUseCase().DelBinary(userPassword, delBinaryID)
	},
}

var delBinaryID string

func init() {
	Binary.Flags().StringVarP(&delBinaryID, "id", "i", "", "Binary id")
	if err := Binary.MarkFlagRequired("id"); err != nil {
		log.Println(err)
		return
	}
}
