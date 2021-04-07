package sdk

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/spf13/pflag"
)

var _ plugcmd.Commander = &Cmd{}
var _ plugcmd.SubCommander = &Cmd{}
var _ plugins.Needer = &Cmd{}
var _ plugins.Plugin = &Cmd{}
var _ plugins.Scoper = &Cmd{}

type Cmd struct {
	pluginsFn plugins.Feeder
	flags     *pflag.FlagSet
	help      bool
}

func (cmd *Cmd) PluginName() string {
	return "sdk"
}

func (cmd *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if cmd.pluginsFn == nil {
		return plugs
	}

	for _, p := range cmd.pluginsFn() {
		switch p.(type) {
		case Needer:
			plugs = append(plugs, p)
		case Sdker:
			plugs = append(plugs, p)
		case Installer:
			plugs = append(plugs, p)
		case Lister:
			plugs = append(plugs, p)
		case VersionUser:
			plugs = append(plugs, p)
		case VersionLister:
			plugs = append(plugs, p)
		case Downloader:
			plugs = append(plugs, p)
		}
	}

	return plugs
}

func (cmd *Cmd) WithPlugins(f plugins.Feeder) {
	cmd.pluginsFn = f
}

func (cmd *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin

	for _, p := range cmd.ScopedPlugins() {
		switch p.(type) {
		case Sdker:
			plugs = append(plugs, p)
		}
	}

	return plugs
}
