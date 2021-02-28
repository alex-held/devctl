package logging

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
	Output    *bytes.Buffer
	Error     *bytes.Buffer
	IsVerbose bool
	OutWriter io.Writer
	ErrWriter io.Writer
}

var defaults = []Option{
	WithFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
		PadLevelText:     true,
		QuoteEmptyFields: true,
	}),
	WithLevel(logrus.TraceLevel),
	WithOutputs(os.Stdout),
	WithErrorOutputs(os.Stderr),
	WithVerbose(false),
}

type Option func(*Logger) *Logger

type errorHook struct {
	Writer io.Writer
}

func (e errorHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}
}

func (e errorHook) Fire(entry *logrus.Entry) error {
	s, err := entry.String()

	if err != nil {
		return errors.Wrapf(err, "failed to fire errhook")
	}
	_, err = e.Writer.Write([]byte(s))
	if err != nil {
		return errors.Wrapf(err, "failed to write errhook logmessage to output writer")
	}
	return nil
}

type outputHook struct {
	Verbose bool
	Writer  io.Writer
}

func (e outputHook) Levels() (levels []logrus.Level) {
	if !e.Verbose {
		return []logrus.Level{logrus.InfoLevel, logrus.WarnLevel}
	}
	return []logrus.Level{logrus.TraceLevel, logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel}
}

func (e outputHook) Fire(entry *logrus.Entry) error {
	s, err := entry.String()
	if err != nil {
		return errors.Wrapf(err, "failed to fire errhook")
	}
	_, err = e.Writer.Write([]byte(s))
	if err != nil {
		return errors.Wrapf(err, "failed to write errhook logmessage to output writer")
	}
	return nil
}

func NewLogger(opts ...Option) *Logger {
	logger := &Logger{
		Logger: logrus.New(),
		Output: &bytes.Buffer{},
		Error:  &bytes.Buffer{},
	}

	for _, opt := range defaults {
		logger = opt(logger)
	}

	for _, opt := range opts {
		logger = opt(logger)
	}

	oHook := outputHook{Writer: logger.OutWriter, Verbose: logger.IsVerbose}
	errHook := errorHook{Writer: logger.ErrWriter}
	logger.Logger.AddHook(oHook)
	logger.Logger.AddHook(errHook)

	logger.Logger.SetOutput(&strings.Builder{})

	return logger
}

func WithVerbose(verbose bool) Option {
	return func(l *Logger) *Logger {
		if verbose {
			l.Logger.SetLevel(logrus.TraceLevel)
			l.IsVerbose = true
			return l
		}
		l.Logger.SetLevel(logrus.WarnLevel)
		l.IsVerbose = false
		return l
	}
}

func WithFormatter(format logrus.Formatter) Option {
	return func(l *Logger) *Logger {
		l.Logger.SetFormatter(format)
		return l
	}
}

func WithOutputs(w ...io.Writer) Option {
	return func(l *Logger) *Logger {
		l.OutWriter = io.MultiWriter(append(w, l.Output)...)
		return l
	}
}

func WithErrorOutputs(w ...io.Writer) Option {
	return func(l *Logger) *Logger {
		l.ErrWriter = io.MultiWriter(append(w, l.Error)...)
		return l
	}
}

func WithLevel(level logrus.Level) Option {
	return func(l *Logger) *Logger {
		l.Logger.SetLevel(level)
		return l
	}
}
