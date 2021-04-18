package plugins

import (
	"context"

	"github.com/gobuffalo/plugins"
)



type Executor interface {
	plugins.Plugin
	ExecuteCommand(ctx context.Context, root string, args []string) error
}
