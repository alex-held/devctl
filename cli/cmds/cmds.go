package cmds

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/alex-held/devctl/pkg/plugins"
)

func Plugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, insidePlugins()...)
	plugs = append(plugs, outsidePlugins()...)
	plugs = append(plugs, Plugins()...)
	return plugs
}

func AvailablePlugins(root string) []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, Plugins()...)
	plugs = append(plugs, Plugins()...)
	plugs = append(plugs, Plugins()...)

	if IsDevctl(root) {
		plugs = append(plugs, insidePlugins()...)
		return plugs
	}
	plugs = append(plugs, outsidePlugins()...)
	return plugs
}

func outsidePlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, Plugins()...)
	plugs = append(plugs, Plugins()...)
	return plugs
}

func insidePlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, Plugins()...)
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
