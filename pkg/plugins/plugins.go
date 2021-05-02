package plugins

import (
	"context"
	"io"

	"github.com/gobuffalo/plugins"
)

type Executor interface {
	plugins.Plugin
	ExecuteCommand(ctx context.Context, root string, args []string) error
}

type StdoutNeeder interface {
	SetStdout(io.Writer) error
}

type Plugin interface {
	PluginName() string
}

type SDKPlugin interface {
	StdoutNeeder
	Plugin
	Install(ctx context.Context, args []string) error
	Download(ctx context.Context, args []string) error
	List(ctx context.Context, args []string) error
	Current(ctx context.Context, args []string) error
	Use(ctx context.Context, args []string) error
}
