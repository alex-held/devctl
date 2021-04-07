package golang

import (
	"github.com/gobuffalo/plugins"

	"github.com/alex-held/devctl/cli/internal/golang/download"
	"github.com/alex-held/devctl/cli/internal/golang/list"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&list.GoListerCmd{},
		&download.GoDownloadCmd{},
	}
}
