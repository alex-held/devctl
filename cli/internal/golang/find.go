package golang

import (
	"path"
	"strings"

	"github.com/gobuffalo/plugins"

	"github.com/alex-held/devctl/cli/cmds/sdk"
)

func FindSubcommandFromArgs(args []string, plugs []plugins.Plugin) plugins.Plugin {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		return FindSubcommand(a, plugs)
	}
	return nil
}

func FindSubcommand(name string, plugs []plugins.Plugin) plugins.Plugin {
	for _, p := range plugs {
		c, ok := p.(sdk.Sdker)
		if !ok {
			continue
		}

		if n, ok := c.(sdk.Namer); ok {
			if n.CmdName() == name {
				return n
			}
		}
		if name == path.Base(p.PluginName()) {
			return p
		}
	}
	return nil
}
