package sdk

import (
	"path"
	"strings"

	"github.com/gobuffalo/plugins"
)

func FindSdkerFromArgs(args []string, plugs []plugins.Plugin) Sdker {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		return FindSdker(a, plugs)
	}
	return nil
}

func FindSdker(name string, plugs []plugins.Plugin) Sdker {
	// Find wraps the other cmd finders into a mega finder for cmds
	for _, p := range plugs {
		c, ok := p.(Sdker)
		if !ok {
			continue
		}

		if n, ok := c.(Namer); ok {
			if n.CmdName() == name {
				return c
			}
		}

		if name == path.Base(c.PluginName()) {
			return c
		}
	}
	return nil
}
