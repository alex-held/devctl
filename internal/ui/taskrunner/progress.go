package taskrunner

import (
	"context"
	"time"

	"github.com/pterm/pterm"
)

func NewTaskRunner(opts ...Option) (tr *taskRunner) {
	tr = &taskRunner{
		DoneC:  make(chan struct{}),
		TaskMC: make(chan TaskRunnerMsg),
	}

	for _, opt := range defaultOptions {
		tr = opt(tr)
	}

	for _, opt := range opts {
		tr = opt(tr)
	}

	return tr
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
			r.TaskMC <- &taskRunnerStartMsg{t.Description}

			// Start the Task
			err := t.Plugin.ExecuteCommand(ctx, t.Root, t.Args)

			time.Sleep(r.AfterTaskTimeout)

			// Communicate that a Task has been completed
			r.TaskMC <- &taskRunnerEndMsg{
				message: t.Description,
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

type taskRunnerStartMsg struct {
	message string
}
func (t *taskRunnerStartMsg) Error() error { return nil }
func (t *taskRunnerStartMsg) Print(o TaskRunnerOutput) {o.PrintTaskProgress(t.message)}




type taskRunnerEndMsg struct {
	message string
	error   error
}
func (t *taskRunnerEndMsg) Error() error {return t.error}
func (t *taskRunnerEndMsg) Print(o TaskRunnerOutput) {
	if t.error != nil {
		o.ErrorF("%v", t.error)
		return
	}
	o.Printf(t.message)
}
