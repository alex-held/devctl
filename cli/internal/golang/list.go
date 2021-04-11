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

func (c *GoListerCmd) Init() {
	if c.path == nil {
		c.path = devctlpath.DefaultPather()
	}
	if c.fs == nil {
		c.fs = afero.NewOsFs()
	}
}

func (l *GoListerCmd) CmdName() string {
	return "list"
}

func (l *GoListerCmd) PluginName() string {
	return "sdk/go/list"
}

func (l *GoListerCmd) ExecuteCommand(_ context.Context, _ string, _ []string) error {
	l.Init()
	fis, err := afero.ReadDir(l.fs, l.path.SDK("go"))
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
