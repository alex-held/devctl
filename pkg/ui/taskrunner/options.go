package taskrunner

import (
	"time"

	"github.com/pterm/pterm"
)



type Option func(tr *taskRunner) *taskRunner


var defaultOptions = []Option{
	WithPTermOutput(&ptermTaskRunnerOutput{
		Initializer: func() *pterm.ProgressbarPrinter {
			return pterm.DefaultProgressbar.
				WithTitle("Default Task Runner")
		},
		Err: pterm.Error,
		Out: pterm.Success,
	}),
	WithTimeout(500 * time.Millisecond),
	WithTitle("Default Task Runner"),
}






type NoOpOutput struct{}

func (n NoOpOutput) ErrorF(_ string, _ ...interface{}) {}
func (n NoOpOutput) Printf(_ string, _ ...interface{}) {}
func (n NoOpOutput) PrintTaskProgress(_ string)        {}
func (n NoOpOutput) Next()                             {}

func WithDiscardOutput() Option {
	return func(tr *taskRunner) *taskRunner {
		tr.output = &NoOpOutput{}
		return tr
	}
}

func WithPTermOutput(output *ptermTaskRunnerOutput) Option {
	return func(tr *taskRunner) *taskRunner {
		tr.output = output
		return tr
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(tr *taskRunner) *taskRunner {
		tr.AfterTaskTimeout = timeout
		return tr
	}
}

func WithTasks(tasks ...Tasker) Option {
	return func(tr *taskRunner) *taskRunner {
		tr.Tasks = append(tr.Tasks, tasks...)
		return tr
	}
}

func WithTitle(title string) Option {
	return func(tr *taskRunner) *taskRunner {
		tr.Title = title
		return tr
	}
}
