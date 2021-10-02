package plugin

import (
	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/pkg/env"
)

func NewCmd(f env.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "manages devctl plugins",
	}

	cmd.AddCommand(newSearchCmd(f))
	cmd.AddCommand(newUpdateCmd(f))
	cmd.AddCommand(NewIndexCommand(f))
	cmd.AddCommand(NewInstallCmd(f))
	cmd.AddCommand(NewUninstallCmd(f))
	
	return cmd
}
