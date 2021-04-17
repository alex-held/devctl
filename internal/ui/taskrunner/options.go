package taskrunner

import (
	"time"
)

type Option func(tr *taskRunner) *taskRunner

func WithPTermOutput(output *ptermTaskRunnerOutput) Option {
	return func(tr *taskRunner) *taskRunner {
		tr.output = output
		return tr
	}
}

func WithTaskRunnerOutput(progressbarFn func(title string), outFn func(format string, args ...interface{}), errFn func(format string, args ...interface{})) Option {
	return func(tr *taskRunner) *taskRunner {
		tr.output = &bufferTaskRunnerOutput{
			ProgressBarPrinterFn: progressbarFn,
			ErrPrinterFn:         errFn,
			OutPrinterFn:         outFn,
		}
		return tr
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(tr *taskRunner) *taskRunner {
		tr.AfterTaskTimeout = timeout
		return tr
	}
}

func WithTasks(tasks ...task) Option {
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
