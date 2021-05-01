package cli

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"

	"github.com/alex-held/devctl/cli/cmds"
)

var _ plugcmd.SubCommander = &Devctl{}
var _ plugins.Plugin = &Devctl{}
var _ plugins.Scoper = &Devctl{}
var _ plugprint.Describer = &Devctl{}

// Devctl represents the `devctl` cli.
type Devctl struct {
	Plugins []plugins.Plugin
	root    string
}

func NewFromRoot(root string) (*Devctl, error) {
	b := &Devctl{
		root: root,
	}

	//	b.Plugins = append(b.Plugins, cmds.AvailablePlugins(root)...)
	alexCommands := cmds.AvailablePlugins(root)
	_ = alexCommands

	// pre scope the plugins to thin the initial set
	plugs := b.ScopedPlugins()
	plugins.Sort(plugs)
	b.Plugins = plugs

	pfn := b.ScopedPlugins

	for _, p := range b.Plugins {
		if t, ok := p.(plugins.Needer); ok {
			t.WithPlugins(pfn)
		}
	}

	return b, nil
}

func (b Devctl) ScopedPlugins() []plugins.Plugin {
	root := b.root
	plugs := make([]plugins.Plugin, 0, len(b.Plugins))
	for _, p := range b.Plugins {
		if t, ok := p.(AvailabilityChecker); ok {
			if !t.PluginAvailable(root) {
				continue
			}
		}
		plugs = append(plugs, p)
	}
	return plugs
}

func (b Devctl) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range b.ScopedPlugins() {
		if _, ok := p.(Commander); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (Devctl) PluginName() string { return "devctl" }

func (Devctl) String() string { return "devctl" }

func (Devctl) Description() string {
	return "Tools for working with devctl applications"
}
