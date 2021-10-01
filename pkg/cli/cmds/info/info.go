package info

import (
	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/pkg/cli/options"
	"github.com/alex-held/devctl/pkg/cli/util"
)

type InfoOptions struct {
	options.IOStreams
}

func NewOptions(streams options.IOStreams) *InfoOptions {
	return &InfoOptions{
		IOStreams: streams,
	}
}

func (o *InfoOptions) Run(f util.Factory, cmd *cobra.Command) error {
	f.Logger().Infof("devctl info\nPath=%s\n", f.Pather().ConfigRoot())
	return nil
}

func NewCmd(f util.Factory, streams options.IOStreams) (cmd *cobra.Command) {
	o := NewOptions(streams)

	cmd = &cobra.Command{
		Use:                   "info",
		DisableFlagsInUseLine: true,
		Short:                 "prints devctl info",
		Long:                  "prints devctl info",
		Example: `
		To get devctl info:
			devctl info
		`,
		Run: func(cmd *cobra.Command, args []string) {
			// util.CheckErr(o.Complete(f, cmd))
			// util.CheckErr(o.ValidateArgs(cmd, args))
			util.CheckErr(o.Run(f, cmd))
		},
	}

	return cmd
}
