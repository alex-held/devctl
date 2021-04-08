package golang

import (
	"context"
	"fmt"

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

// ScopedPlugins
func (c *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin

	fmt.Println("inside scoped plugins")
	if c.pluginsFn == nil {
		fmt.Printf("%#v", *c)
		return plugs
	}
	plugs = append(plugs, c.Plugins...)

	for _, p := range c.pluginsFn() {
		switch p.(type) {
		case Lister:
			plugs = append(plugs, p)
		case Downloader:
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
	fmt.Println("go sdk")

	plugs := c.ScopedPlugins()

	fmt.Printf("found %d scoped plugins = %v\n", len(plugs), plugs)
	for i, plugin := range plugs {
		fmt.Printf("scoped plugins %d= %v\n", i, plugin.PluginName())
		switch p := plugin.(type) {
		case Lister:
			fmt.Printf("switching gosdk plugin: Lister\n", plugin)
			if args[0] == "list" {
				fmt.Printf("execute gosdk plugin: Lister\n", plugin)
				return p.List(ctx, root, args)
			}
		case Downloader:
			fmt.Printf("switching gosdk plugin: Downloader\n", plugin)
			if args[0] == "download" {
				fmt.Printf("execute gosdk plugin: Downloader\n", plugin)
				return p.Download(ctx, root, args)
			}
		}
	}
	return nil
}
