package plugins

import (
	"fmt"
	"io"

	"github.com/schollz/progressbar/v3"
)

func NewProgress(out io.Writer, size int, description string, opts ...progressbar.Option) (bar *progressbar.ProgressBar) {
	options := []progressbar.Option{
		progressbar.OptionSetWriter(out),
		progressbar.OptionShowBytes(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(40), //nolint:gomnd
		progressbar.OptionSpinnerType(1),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan][1/1][reset] %s", description)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			_, _ = io.WriteString(out, "\n")
		}),
	}

	bar = progressbar.NewOptions64(int64(size),
		append(options, opts...)...,
	)

	return bar
}
