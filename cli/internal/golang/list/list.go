package list

import (
	"context"
	"fmt"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/cli/cmds/sdk"
)

var _ plugcmd.Namer = &GoListerCmd{}
var _ plugins.Plugin = &GoListerCmd{}
var _ sdk.Sdker = &GoListerCmd{}

type GoListerCmd struct {
	pluginsFn plugins.Feeder
}

func (l *GoListerCmd) CmdName() string {
	return "go"
}

func (l *GoListerCmd) Sdk(ctx context.Context, root string, args []string) error {
	return l.List(ctx, root, args)
}

func (l *GoListerCmd) PluginName() string {
	return "sdk/go/list"
}

func (l *GoListerCmd) List(ctx context.Context, root string, args []string) error {
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
