package taskrunner

import (
	"fmt"

	"github.com/pterm/pterm"
)

/*
	ptermTaskRunnerOutput
*/

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

/*
	bufferTaskRunnerOutput
*/

type bufferTaskRunnerOutput struct {
	ResetFn                    func() error
	NextFn                     func()
	ProgressBarPrinterFn       func(string string)
	ErrPrinterFn, OutPrinterFn func(format string, args ...interface{})
}

func (b *bufferTaskRunnerOutput) Next() {
	b.NextFn()
}

func (b *bufferTaskRunnerOutput) Reset() error {
	return b.ResetFn()
}

func (b *bufferTaskRunnerOutput) ErrorF(format string, args ...interface{}) {
	if b.ErrPrinterFn != nil {
		b.ErrPrinterFn(format, args...)
		return
	}
}

func (b *bufferTaskRunnerOutput) Printf(format string, args ...interface{}) {
	if b.OutPrinterFn != nil {
		b.OutPrinterFn(format, args...)
		return
	}
}

func (b *bufferTaskRunnerOutput) PrintTaskProgress(title string) {
	if b.ProgressBarPrinterFn != nil {
		b.ProgressBarPrinterFn(title)
		return
	}
}
