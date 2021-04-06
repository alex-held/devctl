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
	exitFn    func(c int)
}

func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry { return l.logger.WithFields(fields) }
func (l *Logger) WithError(err error) *logrus.Entry             { return l.logger.WithError(err) }
func (l *Logger) Printf(format string, args ...interface{})     { l.logger.Printf(format, args...) }
func (l *Logger) Warningf(format string, args ...interface{})   { l.logger.Warningf(format, args...) }
func (l *Logger) Debug(args ...interface{})                     { l.logger.Debug(args...) }
func (l *Logger) Info(args ...interface{})                      { l.logger.Info(args...) }
func (l *Logger) Print(args ...interface{})                     { l.logger.Print(args...) }
func (l *Logger) Warn(args ...interface{})                      { l.logger.Warn(args...) }
func (l *Logger) Warning(args ...interface{})                   { l.logger.Warning(args...) }
func (l *Logger) Error(args ...interface{})                     { l.logger.Error(args...) }
func (l *Logger) Fatal(args ...interface{})                     { l.logger.Fatal(args...) }
func (l *Logger) Panic(args ...interface{})                     { l.logger.Panic(args...) }
func (l *Logger) Debugln(args ...interface{})                   { l.logger.Debugln(args...) }
func (l *Logger) Infoln(args ...interface{})                    { l.logger.Infoln(args...) }
func (l *Logger) Println(args ...interface{})                   { l.logger.Println(args...) }
func (l *Logger) Warnln(args ...interface{})                    { l.logger.Warnln(args...) }
func (l *Logger) Warningln(args ...interface{})                 { l.logger.Warningln(args...) }
func (l *Logger) Errorln(args ...interface{})                   { l.logger.Errorln(args...) }
func (l *Logger) Fatalln(args ...interface{})                   { l.logger.Fatalln(args...) }
func (l *Logger) Panicln(args ...interface{})                   { l.logger.Panicln(args...) }

var defaults = []Option{
	WithFormatter(&DevCtlFormatter{}),
	WithOutputs(os.Stdout),
	WithName(""),
	WithExitFn(os.Exit),
	WithLevel(LogLevelInfo),
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
