package testutils

import (
	"flag"
	"os"

	"github.com/alex-held/devctl/pkg/logging"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logging.Logger {
	verboseFlag := flag.CommandLine.Lookup("test.v")

	var level logrus.Level
	var verbose bool
	switch verboseFlag.Value.String() {
	case "true":
		verbose = true
		level = logrus.TraceLevel
	default:
		verbose = false
		level = logrus.WarnLevel
	}

	logger := logging.NewLogger(
		logging.WithVerbose(verbose),
		logging.WithLevel(level),
		logging.WithOutputs(),
		logging.WithErrorOutputs(os.Stderr),
		logging.WithFormatter(&logrus.TextFormatter{}),
	)

	return logger
}
