package logging

import (
	"bytes"
	"io"

	"github.com/sirupsen/logrus"
)

type Option func(*Logger) *Logger

func WithDebugFunc(debugFn DebugFunc) Option {
	return func(l *Logger) *Logger {
		l.DebugFunc = debugFn
		return l
	}
}

func WithFormatter(format logrus.Formatter) Option {
	return func(l *Logger) *Logger {
		l.logger.SetFormatter(format)
		return l
	}
}

func WithBuffer(b *bytes.Buffer) Option {
	return func(l *Logger) *Logger {
		l.Output = b
		return l
	}
}

func WithOutputs(w ...io.Writer) Option {
	return func(l *Logger) *Logger {
		l.OutWriter = io.MultiWriter(w...)
		return l
	}
}

func WithLevel(level LogLevel) Option {
	return func(l *Logger) *Logger {
		l.level = level
		lvl, _ := logrus.ParseLevel(level.String())
		l.logger.SetLevel(lvl)
		l.logger.Level = lvl
		return l
	}
}

func WithName(name string) Option {
	return func(l *Logger) *Logger {
		l.name = name
		return l
	}
}

func WithVerbose(verbose bool) Option {
	return func(l *Logger) *Logger {
		l.IsVerbose = verbose
		return l
	}
}
