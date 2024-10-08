package add

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/nextlag/keeper/internal/client/usecase"
	"github.com/nextlag/keeper/internal/entity"
	utils "github.com/nextlag/keeper/internal/utils/client"
)

var Binary = &cobra.Command{
	Use:   "binary",
	Short: "Add binary",
	Long: fmt.Sprintf(`
This command add user file
Example:
  %s add binary -t "name" -f "file_location" --meta '[{"name":"meta","value":"value"}]'`, App),

	Run: func(cmd *cobra.Command, args []string) {
		userPassword, err := usecase.GetClientUseCase().GetTempPass()
		if err != nil {
			color.Red("Authentication required. Error: %v", err)
			return
		}
		usecase.GetClientUseCase().AddBinary(userPassword, &binaryForAdditing)
	},
}

var binaryForAdditing entity.Binary

func init() {
	Binary.Flags().StringVarP(&binaryForAdditing.Name, "title", "t", "", "Login title")
	Binary.Flags().StringVarP(&binaryForAdditing.FileName, "file", "f", "", "User file")
	Binary.Flags().Var(&utils.JSONFlag{Target: &binaryForAdditing.Meta}, "meta", `Add meta fields for entity`)

	if err := Binary.MarkFlagRequired("title"); err != nil {
		color.Red("%v", err)
		return
	}
	if err := Binary.MarkFlagRequired("file"); err != nil {
		color.Red("%v", err)
		return
	}
}
