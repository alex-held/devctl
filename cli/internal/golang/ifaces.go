package golang

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"

	"github.com/alex-held/devctl/cli/cmds/sdk"
)

type Stdouter = plugio.Outer
type Needer = plugins.Needer
type Namer = sdk.Namer

type SubCommander = plugcmd.SubCommander

type Initer interface {
	Init()
}
