package progress

import "io"

type progressBar struct {
	o  io.Writer
	pr progressReporter
}

type progressReporter struct {
}

func (r progressReporter) report(size int, current int) error {
	return nil
}
