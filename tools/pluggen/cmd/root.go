package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "pluggen",
	}
	cmd.AddCommand(NewGenCmd())
	return cmd
}
