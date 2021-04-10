package golang

import (
	"path"
	"strings"

	"github.com/gobuffalo/plugins"
)

func FindSubcommandFromArgs(args []string, plugs []plugins.Plugin) GoSDKCommander {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		return FindSubcommand(a, plugs)
	}
	return nil
}

func FindSubcommand(name string, plugs []plugins.Plugin) GoSDKCommander {
	for _, p := range plugs {
		if cmd, ok := p.(GoSDKCommander); ok {
			if cmd.CmdName() == name {
				return cmd
			}
		}
		if name == path.Base(p.PluginName()) {
			return p.(GoSDKCommander)
		}
	}
	return nil
}
