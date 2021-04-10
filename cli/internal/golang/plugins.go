package golang

import (
	"github.com/gobuffalo/plugins"

	"github.com/alex-held/devctl/cli/internal/golang/download"
	"github.com/alex-held/devctl/cli/internal/golang/list"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Cmd{
			Plugins: []plugins.Plugin{
				&download.GoDownloadCmd{},
				&list.GoListerCmd{},
			},
		},
	}
}
