package list

import (
	"context"
	"fmt"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/spf13/afero"
)

var _ plugcmd.Namer = &GoListerCmd{}
var _ plugins.Plugin = &GoListerCmd{}

type GoListerCmd struct {
	pluginsFn plugins.Feeder
}

func (l *GoListerCmd) CmdName() string {
	return "list"
}

func (l *GoListerCmd) PluginName() string {
	return "sdk/go/list"
}

func (l *GoListerCmd) ExecuteCommand(ctx context.Context, root string, args []string) error {
	fs := afero.NewOsFs()
	fis, err := afero.ReadDir(fs, "/Users/dev/.devctl/sdks/go")
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if fi.Name() != "current" {
			fmt.Println(fi.Name())
		}
	}

	return nil
}
