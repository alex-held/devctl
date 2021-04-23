package cmds

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/plugins"

	"github.com/alex-held/devctl/cli/cmds/sdk"
	"github.com/alex-held/devctl/cli/cmds/version"
	"github.com/alex-held/devctl/cli/internal/golang"
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

	if IsDevctl(root) {
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


func IsDevctl(mod string) bool {
	if !strings.HasPrefix(mod, "go.mod") {
		mod = filepath.Join(mod, "go.mod")
	}

	b, err := ioutil.ReadFile(mod)
	if err != nil {
		return false
	}

	return bytes.Contains(b, []byte("github.com/alex-held/devctl"))
}
