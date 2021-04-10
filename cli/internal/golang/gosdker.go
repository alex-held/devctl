package golang

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"

	"github.com/alex-held/devctl/cli/internal/golang/download"
	"github.com/alex-held/devctl/cli/internal/golang/list"
)

var _ plugcmd.SubCommander = &Cmd{}
var _ Namer = &Cmd{}
var _ plugins.Plugin = &Cmd{}
var _ plugins.Scoper = &Cmd{}
var _ plugprint.Describer = &Cmd{}

type Cmd struct {
	Plugins   []plugins.Plugin
	pluginsFn plugins.Feeder
	flags     *pflag.FlagSet
	help      bool
}

func (c *Cmd) Description() string {
	return "manages the installations of the go sdk"
}

func (c *Cmd) SubCommands() []plugins.Plugin {
	return []plugins.Plugin{
		&list.GoListerCmd{},
		&download.GoDownloadCmd{},
	}
}

func (c *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if c.pluginsFn == nil {
		return plugs
	}

	plugs = append(plugs, c.Plugins...)
	for _, p := range c.pluginsFn() {
		switch p.(type) {
		case GoSDKCommander:
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (c *Cmd) WithPlugins(f plugins.Feeder) {
	c.pluginsFn = f
}

func (c *Cmd) CmdName() string {
	return "sdk/go"
}

func (c *Cmd) PluginName() string {
	return "go"
}

func (c *Cmd) Sdk(ctx context.Context, root string, args []string) error {
	plugs := c.ScopedPlugins()
	subcommand := FindSubcommandFromArgs(args, plugs)
	err := subcommand.ExecuteCommand(ctx, root, args)
	return err
}
