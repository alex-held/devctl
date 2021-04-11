package golang

import (
	"context"
	"fmt"

	"github.com/alex-held/devctl/cli/cmds/sdk"
)

func (c *GoSDKCmd) Main(ctx context.Context, root string, args []string) error {
	plugs := c.ScopedPlugins()
	subcommand := FindSubcommandFromArgs(args, plugs)

	switch cmd := subcommand.(type) {
	case sdk.Sdker:
		if i, ok := cmd.(Initer); ok {
			i.Init()
		}
		return cmd.ExecuteCommand(ctx, root, args)
	default:
		return fmt.Errorf("plugin %s has a unsupported api", cmd.PluginName())
	}
}
