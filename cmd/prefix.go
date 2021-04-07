package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/internal/app"
)

// NewPrefixCommand creates the `devenv prefix` commands
func NewPrefixCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prefix",
		Short: "devctl prefix",
		Args:  cobra.RangeArgs(0, 1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return listSdks(), cobra.ShellCompDirectiveDefault
		},
		Run: func(c *cobra.Command, args []string) {
			cli := app.GetCLI()
			pather := cli.GetPather()
			rootPath := pather.ConfigRoot()
			fmt.Println(rootPath)
		},
	}

	return cmd
}
