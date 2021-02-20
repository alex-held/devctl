package testutils

import (
	"flag"
	"io"

	"github.com/alex-held/devctl/internal/logging"

	"github.com/sirupsen/logrus"
)

func NewLogger(out io.Writer) *logrus.Logger {
	verbose := flag.CommandLine.Lookup("test.v")

	logger := logging.NewLogger(
		func(l *logrus.Logger) *logrus.Logger {
			l.SetFormatter(&logrus.TextFormatter{})
			return l
		},
		func(l *logrus.Logger) *logrus.Logger {
			if out != nil {
				l.SetOutput(out)
			}
			return l
		},
		func(l *logrus.Logger) *logrus.Logger {
			switch verbose.Value.String() {
			case "true":
				l.SetLevel(logrus.DebugLevel)
			default:
				l.SetLevel(logrus.WarnLevel)
			}
			return l
		})

	return logger
}
