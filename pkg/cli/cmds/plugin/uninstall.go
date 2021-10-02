package plugin

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/alex-held/devctl/pkg/env"
	"github.com/alex-held/devctl/pkg/index/installation"
	"github.com/alex-held/devctl/pkg/index/validate"
)

// uninstallCmd represents the uninstall command
func NewUninstallCmd(f env.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall plugins",
		Long: `Uninstall one or more plugins.
Example:
  kubectl krew uninstall NAME [NAME...]
Remarks:
  Failure to uninstall a plugin will result in an error and exit immediately.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, name := range args {
				if isCanonicalName(name) {
					return errors.New("uninstall command does not support INDEX/PLUGIN syntax; just specify PLUGIN")
				} else if !validate.IsSafePluginName(name) {
					return unsafePluginNameErr(name)
				}
				klog.V(4).Infof("Going to uninstall plugin %s\n", name)
				if err := installation.Uninstall(f, name); err != nil {
					return errors.Wrapf(err, "failed to uninstall plugin %s", name)
				}
				fmt.Fprintf(os.Stderr, "Uninstalled plugin: %s\n", name)
			}
			return nil
		},
		PreRunE: func(c *cobra.Command, args []string) error {
			return checkIndex(f, c, args)
		},
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"remove"},
	}
}

func unsafePluginNameErr(n string) error { return errors.Errorf("plugin name %q not allowed", n) }
