package golang

import (
	"fmt"
	"path"
	"strings"

	"github.com/gobuffalo/plugins"
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
	// Find wraps the other cmd finders into a mega finder for cmds
	for _, p := range plugs {
		c, ok := p.(GoSdker)
		if !ok {
			continue
		}

		if n, ok := c.(Namer); ok {
			fmt.Println("searching subcommand; name=" + n.CmdName())
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
