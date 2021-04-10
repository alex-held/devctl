package golang

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"

	"github.com/alex-held/devctl/cli/cmds/sdk"
)

type GoSDKCommander interface {
	sdk.Command
}

type Stdouter = plugio.Outer
type Needer = plugins.Needer
type Namer = sdk.Namer

type Downloader sdk.Downloader
type Lister sdk.Lister
type SubCommander = plugcmd.SubCommander
