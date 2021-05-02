package taskrunner

import (
	"context"
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

type Describer interface {
	Describe() string
}

type Runner interface {
	Describer
	Run(ctx context.Context) error
	Wrap(executeWhenTrue ConditionalExecutorFn) Tasker
}

type Tasker interface {
	Describer
	Task(ctx context.Context) (err error)
}

type ConditionalExecutorFn func() bool
type TaskActionFn func(ctx context.Context) error
