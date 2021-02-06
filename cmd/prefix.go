package cmd

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/alex-held/devctl/pkg/cli"
)

// NewPrefixCommand creates the `devenv prefix` commands
func NewPrefixCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prefix",
		Short: "Get SDK prefix",
		Args:  cobra.RangeArgs(0, 1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return listSdks(), cobra.ShellCompDirectiveDefault
		},
		Run: func(c *cobra.Command, args []string) {
			config := viper.ConfigFileUsed()
			homeDir, err := homedir.Dir()
			if err != nil {
				cli.ExitWithError(1, err)
			}
			rootDir := filepath.Dir(config)
			sdkLinkDir := filepath.Join(homeDir, rootDir, "sdks", "current")

			if len(args) == 0 {
				println(sdkLinkDir)
			}

			sdkLink := filepath.Join(sdkLinkDir, args[0])
			println(sdkLink)
		},
	}

	return cmd
}
