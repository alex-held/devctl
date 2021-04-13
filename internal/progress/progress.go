package progress

import "io"

type progressBar struct {
	o io.Writer
	pr progressReporter
}

type progressReporter struct {

}
