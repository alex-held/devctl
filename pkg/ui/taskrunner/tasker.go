package taskrunner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

func NewTaskRunner(opts ...Option) (runner *taskRunner) {
	tr := &taskRunner{
		DoneC:  make(chan struct{}),
		TaskMC: make(chan TaskRunnerMsg),
	}
	for _, opt := range defaultOptions {
		tr = opt(tr)
	}
	for _, opt := range opts {
		tr = opt(tr)
	}
	runner = tr
	return runner
}

type taskRunner struct {
	Title            string
	Tasks            Tasks
	AfterTaskTimeout time.Duration
	output TaskRunnerOutput
	DoneC  chan struct{}
	TaskMC chan TaskRunnerMsg
}

func (r *taskRunner) Run(ctx context.Context) error {
	if pTermOutput, ok := r.output.(*ptermTaskRunnerOutput); ok {
		pTermOutput.p, _ = pTermOutput.Initializer().
			WithRemoveWhenDone(false).
			WithTotal(len(r.Tasks)).
			Start()
	}

	pterm.DefaultSection.
		WithStyle(pterm.NewStyle(pterm.Bold, pterm.FgWhite)).
		WithTopPadding(2).
		WithBottomPadding(1).
		WithIndentCharacter("-->").
		Println(r.Title)

	go func() {
		for _, t := range r.Tasks {
			// Communicate that a Task will be started
			msg := t.Describe()
			r.TaskMC <- &taskRunnerStartMsg{msg}

			// Start the Task
			err := t.Task(ctx)
			time.Sleep(r.AfterTaskTimeout)

			// Communicate that a Task has been completed
			r.TaskMC <- &taskRunnerEndMsg{
				message: t.Describe(),
				error:   err,
			}

			if err != nil {
				break
			}
		}
		r.DoneC <- struct{}{}
	}()

	for {
		select {
		case <-r.DoneC:
			return nil
		case m := <-r.TaskMC:
			m.Print(r.output)
			if err := m.Error(); err != nil {
				return err
			}
			//goland:noinspection GoLinterLocal
		default:
			// no-op
		}
	}
}

func (r *taskRunner) Wrap(executeWhenTrue ConditionalExecutorFn) Tasker {
	tasker := NewConditionalTask(
		r.Title,
		func(ctx context.Context) error {
			return r.Run(ctx)
		},
		executeWhenTrue,
	)
	return tasker
}

func (r *taskRunner) Describe() string {
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "Title:\t%s\n", r.Title)
	fmt.Fprintf(sb, "Tasks:\t%v\n", r.Tasks)
	_, _ = fmt.Fprintln(sb)
	return sb.String()
}

type taskRunnerStartMsg struct {
	message string
}

func (t *taskRunnerStartMsg) Error() error             { return nil }
func (t *taskRunnerStartMsg) Print(o TaskRunnerOutput) { o.PrintTaskProgress(t.message) }

type taskRunnerEndMsg struct {
	message string
	error   error
}

func (t *taskRunnerEndMsg) Error() error { return t.error }
func (t *taskRunnerEndMsg) Print(o TaskRunnerOutput) {
	if t.error != nil {
		o.ErrorF("%v", t.error)
		return
	}
	o.Printf(t.message)
}
