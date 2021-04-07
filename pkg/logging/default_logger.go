package logging

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

type LogLevel int

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelError:
		return "ERROR"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	default:
		panic(fmt.Errorf("this LogLevel `%s` is not supported! ", strconv.Itoa(int(l))))
	}
}

const (
	// Debug messages, write to debug logs only by logutils.Debug.
	LogLevelDebug LogLevel = 0

	// Information messages, don't write too much messages,
	// only useful ones: they are shown when running with -v.
	LogLevelInfo LogLevel = 1

	// Hidden errors: non critical errors: work can be continued, no need to fail whole program;
	// tests will crash if any warning occurred.
	LogLevelWarn LogLevel = 2

	// Only not hidden from user errors: whole program failing, usually
	// error logging happens in 1-2 places: in the "main" function.
	LogLevelError LogLevel = 3
)

type Log interface {
	logrus.FieldLogger
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Logf(level LogLevel, format string, args ...interface{})

	Child(name string) Log
	SetLevel(level LogLevel)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf("%s%s", l.prefix(), fmt.Sprintf(format, args...))
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	v := fmt.Sprintf("%s%s", l.prefix(), fmt.Sprintf(format, args...))
	panic(v)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.level > LogLevelError {
		return
	}
	l.logger.Errorf("%s%s", l.prefix(), fmt.Sprintf(format, args...))
	// don't call exitIfTest() because the idea is to
	// crash on hidden errors (warnings); but Errorf MUST NOT be
	// called on hidden errors, see log levels comments.
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.level > LogLevelWarn {
		return
	}
	l.logger.Warnf("%s%s", l.prefix(), fmt.Sprintf(format, args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.level > LogLevelDebug {
		return
	}
	l.logger.Debugf("%s%s", l.prefix(), fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.level > LogLevelInfo {
		return
	}
	l.logger.Infof("%s%s", l.prefix(), fmt.Sprintf(format, args...))
}

func (l *Logger) Logf(lvl LogLevel, format string, args ...interface{}) {
	if l.level > lvl {
		return
	}
	parsedLevel, _ := logrus.ParseLevel(lvl.String())
	l.logger.Logf(parsedLevel, format, args...)
}

func (l *Logger) Child(name string) Log {
	prefix := ""
	if l.name != "" {
		prefix = l.name + "/"
	}
	child := l
	child.name = prefix + name
	return child
}

func (l *Logger) SetLevel(level LogLevel) { l.level = level }

func (l *Logger) prefix() string {
	prefix := ""
	if l.name != "" {
		prefix = fmt.Sprintf("[%s:>", l.name)
	}
	return prefix
}
