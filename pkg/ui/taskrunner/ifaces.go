package taskrunner

import (
	"context"

	"github.com/gobuffalo/plugins"
)

type TaskRunnerMsg interface {
	Print(output TaskRunnerOutput)
	Error() error
}

type TaskRunnerOutput interface {
	ErrorF(format string, args ...interface{})
	Printf(format string, args ...interface{})
	PrintTaskProgress(string string)
	Next()
}

type Executer interface {
	plugins.Plugin
	ExecuteCommand(ctx context.Context, root string, args []string) error
}

type Describer interface {
	Describe() string
}

type Runner interface {
	Describer
	Run(ctx context.Context) error
}
