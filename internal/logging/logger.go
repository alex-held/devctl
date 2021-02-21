package logging

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var defaults = []Option{
	WithFormatter(&logrus.JSONFormatter{}),
	WithLevel(logrus.InfoLevel),
	WithOutput(os.Stdout),
}

type Option func(*logrus.Logger) *logrus.Logger

func NewLogger(opts ...Option) *logrus.Logger {
	logger := logrus.New()

	for _, opt := range defaults {
		logger = opt(logger)
	}

	for _, opt := range opts {
		logger = opt(logger)
	}

	return logger
}

func WithFormatter(format logrus.Formatter) Option {
	return func(l *logrus.Logger) *logrus.Logger {
		l.SetFormatter(format)
		return l
	}
}

func WithOutput(w io.Writer) Option {
	return func(l *logrus.Logger) *logrus.Logger {
		l.SetOutput(w)
		return l
	}
}

func WithLevel(level logrus.Level) Option {
	return func(l *logrus.Logger) *logrus.Logger {
		l.SetLevel(level)
		return l
	}
}
