package taskrunner

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gobuffalo/plugins"
	"github.com/pterm/pterm"
)

var defaultOptions = []Option{
	WithPTermOutput(&ptermTaskRunnerOutput{
		Initializer: func() *pterm.ProgressbarPrinter {
			return pterm.DefaultProgressbar.
				WithTitle("Default Task Runner")
		},
		Err: pterm.Error,
		Out: pterm.Success,
	}),
	WithTimeout(defaultTaskTimeout),
	WithTitle("Default Task Runner"),
}

type Tasks []task

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

var defaultTaskTimeout = 500 * time.Millisecond

type TaskRunnerOutput interface {
	ErrorF(format string, args ...interface{})
	Printf(format string, args ...interface{})
	PrintTaskProgress(string string)
	Next()
}

type taskRunner struct {
	Title            string
	Tasks            Tasks
	AfterTaskTimeout time.Duration

	output TaskRunnerOutput

	DoneC  chan struct{}
	TaskMC chan TaskRunnerMsg
}

type Executer interface {
	plugins.Plugin
	ExecuteCommand(ctx context.Context, root string, args []string) error
}

type task struct {
	Plugin      Executer
	Description string
	Root        string
	Args        []string
}

func (r *taskRunner) Run(ctx context.Context) error {

//	const padding = 3
//	w := indent.NewWriterPipe(os.Stdout, 4, nil)

	//w := tabwriter.NewWriter(os.Stdout, 5, 0, padding, ' ', tabwriter.Debug)
//	pterm.SetDefaultOutput(w)

	if pTermOutput, ok := r.output.(*ptermTaskRunnerOutput); ok {
		pTermOutput.p, _ = pTermOutput.Initializer().
			WithRemoveWhenDone(false).
			WithTotal(len(r.Tasks)).
			Start()
	}

	go func() {
		for _, t := range r.Tasks {
			// Communicate that a task will be started
			r.TaskMC <- &taskRunnerStartMsg{t.Description}

			// Start the task
			err := t.Plugin.ExecuteCommand(ctx, t.Root, t.Args)

			time.Sleep(r.AfterTaskTimeout)

			// Communicate that a task has been completed
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

	err := r.display()

	return err
}

type indentWriter struct {
	indentChar string
	level      int
	w          io.Writer
}

func (w *indentWriter) GetIndent() string {
	return strings.Repeat(w.indentChar, w.level)
}

func (w *indentWriter) Write(p []byte) (n int, err error) {
	reader := bufio.NewReader(bytes.NewReader(p))
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var indentedBytes []byte

	for scanner.Scan() {
		ln := scanner.Text()
		indent := w.GetIndent()
		fmt.Fprintf(w.w, "%s%s\n", indent, ln)
		/*indentedLn := fmt.Sprintf("%s%s\n", indent, ln)
		indentedLnBytes := []byte(indentedLn)
		indentedBytes = append(indentedBytes, indentedLnBytes...)*/
	}

	return w.w.Write(indentedBytes)
}

func (r *taskRunner) display() error {
	pterm.DefaultSection.
		WithStyle(pterm.NewStyle(pterm.Bold, pterm.FgWhite)).
		WithTopPadding(2).
		WithBottomPadding(1).
		WithIndentCharacter("-->").
		Println(r.Title)

	for {
		select {
		case <-r.DoneC:
			return nil
		case m := <-r.TaskMC:
			m.Print(r.output)
			//goland:noinspection GoLinterLocal
		default:
			// no-op
		}
	}
}

type taskRunnerStartMsg struct {
	message string
}

func (t *taskRunnerStartMsg) Print(o TaskRunnerOutput) {
	o.PrintTaskProgress(t.message)
}

type taskRunnerEndMsg struct {
	message string
	error   error
}

func (t *taskRunnerEndMsg) Print(o TaskRunnerOutput) {
	if t.error != nil {
		o.ErrorF("%v", t.error)
		return
	}
	o.Printf(t.message)
}

type TaskRunnerMsg interface {
	Print(output TaskRunnerOutput)
}
