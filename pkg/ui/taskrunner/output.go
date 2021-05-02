package taskrunner

import (
	"fmt"

	"github.com/pterm/pterm"
)

type ptermTaskRunnerOutput struct {
	Initializer func() *pterm.ProgressbarPrinter
	p           *pterm.ProgressbarPrinter
	Err         pterm.PrefixPrinter
	Out         pterm.PrefixPrinter
}

func (p *ptermTaskRunnerOutput) Next() {
	p.p.Increment()
}

func (p *ptermTaskRunnerOutput) ErrorF(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	p.Err.Println(msg)
	p.p.Increment()
}

func (p *ptermTaskRunnerOutput) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	p.Out.Println(msg)
	p.p.Increment()
}

func (p *ptermTaskRunnerOutput) PrintTaskProgress(title string) {
	p.p.Title = title
}
