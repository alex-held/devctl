package sdk

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
)

func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	plugs := cmd.ScopedPlugins()

	if p := FindSdkerFromArgs(args, plugs); p != nil {
		return p.ExecuteCommand(ctx, root, args[1:])
	}

	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return plugins.Wrap(cmd, err)
	}

	args = flags.Args()

	if cmd.help {
		return plugprint.Print(plugio.Stdout(plugs...), cmd)
	}

	return cmd.run(ctx, root, args)
}

func (cmd *Cmd) run(ctx context.Context, root string, args []string) error {
	plugs := cmd.ScopedPlugins()

	for _, p := range plugs {
		if s, ok := p.(Sdker); ok {
			if err := s.ExecuteCommand(ctx, root, args); err != nil {
				return plugins.Wrap(s, err)
			}
		}
	}

	return nil
}
