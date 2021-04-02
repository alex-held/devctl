package logging

import (
	"bytes"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger    *logrus.Logger
	Output    *bytes.Buffer
	IsVerbose bool
	OutWriter io.Writer
	DebugFunc DebugFunc
	level     LogLevel
	name      string
}

var defaults = []Option{
	WithFormatter(&DevCtlFormatter{}),
	WithOutputs(os.Stdout),
	WithName(""),
}

func NewLogger(opts ...Option) Log {
	logger := &Logger{
		logger: logrus.New(),
		Output: &bytes.Buffer{},
	}

	for _, opt := range defaults {
		logger = opt(logger)
	}

	for _, opt := range opts {
		logger = opt(logger)
	}

	logger.logger.Out = io.MultiWriter(logger.Output, logger.OutWriter)
	return logger
}
