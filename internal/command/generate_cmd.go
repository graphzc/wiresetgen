package command

import (
	"github.com/GraphZC/go-wireset-gen/internal/handlers"
	"github.com/spf13/cobra"
)

func NewGenerateCommand(generateHandler handlers.GenerateHandler) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate wire set",
		Long:  `Generate wire set`,
		Run: func(cmd *cobra.Command, args []string) {
			generateHandler.GenerateWireSet()
		},
	}
}
