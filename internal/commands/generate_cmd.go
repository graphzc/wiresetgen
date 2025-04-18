package commands

import (
	"github.com/graphzc/wiresetgen/internal/handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewGenerateCommand(generateHandler handlers.GenerateHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate wire set",
		Long:  "Generate wire set",
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")

			if err := generateHandler.GenerateWireSet(verbose); err != nil {
				logrus.Error("Error generating wire set:", err)
			} else {
				logrus.Info("Wire set generated successfully")
			}
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	return cmd
}
