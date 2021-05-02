package version

import (
	"github.com/gobuffalo/plugins"
)

var Version = "devctl/unknown"

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Cmd{},
	}
}
