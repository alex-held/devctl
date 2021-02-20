package testutils

import (
	"flag"
	"io"

	"github.com/sirupsen/logrus"
)

func NewLogger(out io.Writer) *logrus.Logger {
	logger := logrus.New()
	if out != nil {
		logger.SetOutput(out)
	}

	verbose := flag.CommandLine.Lookup("test.v")
	switch verbose.Value.String() {
	case "true":
		logger.SetLevel(logrus.DebugLevel)
	default:
		logger.SetLevel(logrus.WarnLevel)
	}

	return logger
}
