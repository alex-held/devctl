package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

type LoggerOptions func(*logrus.Logger) *logrus.Logger

func NewLogger(opts ...LoggerOptions) *logrus.Logger {
	logger := logrus.New()
	for _, opt := range GetDefaultLoggerOptions() {
		logger = opt(logger)
	}

	for _, opt := range opts {
		logger = opt(logger)
	}
	return logger
}

func GetDefaultLoggerOptions() (opts []LoggerOptions) {
	opts = []LoggerOptions{
		func(l *logrus.Logger) *logrus.Logger {
			l.SetFormatter(&logrus.JSONFormatter{})
			return l
		},
		func(l *logrus.Logger) *logrus.Logger {
			l.SetLevel(logrus.InfoLevel)
			return l
		},
		func(l *logrus.Logger) *logrus.Logger {
			l.SetOutput(os.Stdout)
			return l
		},
	}
	return opts
}
