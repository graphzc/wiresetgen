package commands

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "wiresetgen",
		Short: "A generator for wire cli to auto generate wireset",
		Long:  `A generator for wire cli to auto generate wireset`,
	}
}
