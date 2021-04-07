package cmds

import (
	"github.com/gobuffalo/plugins"

	"github.com/alex-held/devctl/cli/cmds/sdk"
	"github.com/alex-held/devctl/cli/cmds/version"
	"github.com/alex-held/devctl/cli/internal/golang"
	"github.com/alex-held/devctl/meta"
)

func Plugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, insidePlugins()...)
	plugs = append(plugs, outsidePlugins()...)
	plugs = append(plugs, version.Plugins()...)
	return plugs
}

func AvailablePlugins(root string) []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, version.Plugins()...)
	plugs = append(plugs, sdk.Plugins()...)
	plugs = append(plugs, golang.Plugins()...)

	if meta.IsDevctl(root) {
		plugs = append(plugs, insidePlugins()...)
		return plugs
	}
	plugs = append(plugs, outsidePlugins()...)
	return plugs
}

func outsidePlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, sdk.Plugins()...)
	plugs = append(plugs, golang.Plugins()...)
	return plugs
}

func insidePlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, golang.Plugins()...)
	return plugs
}
