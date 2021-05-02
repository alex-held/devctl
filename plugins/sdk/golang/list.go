package golang

import (
	"context"
	"fmt"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

var _ plugcmd.Namer = &GoListerCmd{}
var _ plugins.Plugin = &GoListerCmd{}

type GoListerCmd struct {
	fs   afero.Fs
	path devctlpath.Pather
}

func (cmd *GoListerCmd) Init() {
	if cmd.path == nil {
		cmd.path = devctlpath.DefaultPather()
	}
	if cmd.fs == nil {
		cmd.fs = afero.NewOsFs()
	}
}

func (cmd *GoListerCmd) CmdName() string {
	return "list"
}

func (cmd *GoListerCmd) PluginName() string {
	return "sdk/go/list"
}

func (cmd *GoListerCmd) ExecuteCommand(_ context.Context, _ string, _ []string) error {
	cmd.Init()
	fis, err := afero.ReadDir(cmd.fs, cmd.path.SDK("go"))
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
