package cli

import (
	"context"
	"io"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"
)

type Aliaser = plugcmd.Aliaser
type Commander = plugcmd.Commander
type Needer = plugins.Needer
type StderrNeeder = plugio.ErrNeeder
type StdinNeeder = plugio.InNeeder
type StdoutNeeder = plugio.OutNeeder

// AvailabilityChecker
type AvailabilityChecker interface {
	PluginAvailable(root string) bool
}

type ErrNeeder interface {
	SetErr(w io.Writer)
}
type Errer interface {
	Err() io.Writer
}
type OutNeeder interface {
	SetOut(w io.Writer)
}
type Outer interface {
	Out() io.Writer
}
type InNeeder interface {
	SetIn(r io.Reader)
}
type Inner interface {
	In() io.Reader
}

type Initer interface {
	Init(ctx context.Context, root string, args []string) error
}
