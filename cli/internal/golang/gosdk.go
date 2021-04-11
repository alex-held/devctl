package golang

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"

	"github.com/alex-held/devctl/cli/cmds/sdk"
)

var _ plugcmd.SubCommander = &GoSDKCmd{}
var _ Namer = &GoSDKCmd{}
var _ plugins.Plugin = &GoSDKCmd{}
var _ plugins.Scoper = &GoSDKCmd{}
var _ plugprint.Describer = &GoSDKCmd{}

type GoSDKCmd struct {
	Plugins   []plugins.Plugin
	pluginsFn plugins.Feeder
	flags     *pflag.FlagSet
	help      bool
}

func (c *GoSDKCmd) ExecuteCommand(ctx context.Context, root string, args []string) error {
	return c.Main(ctx, root, args)
}

func (c *GoSDKCmd) CmdName() string {
	return "go"
}

func (c *GoSDKCmd) PluginName() string {
	return "sdk/go"
}

func (c *GoSDKCmd) Description() string {
	return "manages the installations of the go sdk"
}

func (c *GoSDKCmd) SubCommands() []plugins.Plugin {
	return []plugins.Plugin{
		&GoListerCmd{},
		&GoDownloadCmd{},
		&GoUseCmd{},
		&GoInstallCmd{},
	}
}

func (c *GoSDKCmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if c.pluginsFn == nil {
		return plugs
	}

	plugs = append(plugs, c.Plugins...)
	for _, p := range c.pluginsFn() {
		switch p.(type) {
		case sdk.Sdker:
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (c *GoSDKCmd) WithPlugins(f plugins.Feeder) {
	c.pluginsFn = f
}
