package testutils

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func NewLogger(out *bytes.Buffer) *logrus.Logger {
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

func SetupTestLogger(t testing.TB) (spy *LogSpy, teardown Teardown) {
	spy = NewTestLogSpy(t)
	zl := NewDefaultZapLogger(spy)
	restoreGlobals := zap.ReplaceGlobals(zl)
	zLogger := zl.Sugar()

	teardown = func() {
		restoreGlobals()
		zl.Sync()
		_ = zLogger.Sync()
	}

	return spy, teardown
}

func NewDefaultZapLogger(spy *LogSpy) *zap.Logger {
	return zaptest.NewLogger(spy, zaptest.Level(zap.WarnLevel), zaptest.WrapOptions(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.OnFatal(zapcore.WriteThenGoexit),
		zap.Hooks(func(entry zapcore.Entry) (err error) {
			switch entry.Level {
			case zap.FatalLevel, zap.PanicLevel:
				spy.Fatalf("[test.fatal]\t%s\n", entry.Message)
			case zap.ErrorLevel:
				err = fmt.Errorf("[test.error]\t%s\n", entry.Message)
				//	spy.ErrorMessages = append(spy.ErrorMessages,errorEntry{Error: err} )
				//	spy.Error(err)
				return err
			case zap.WarnLevel:
				spy.Logf("[test.warn]\t%s\n", entry.Message)
			case zap.DebugLevel, zap.InfoLevel:
				if testing.Verbose() {
					spy.Logf("[test.%s]\t%s\n", strings.ToLower(entry.Level.String()), entry.Message)
				}
			default:
				err = fmt.Errorf("LogLevel %s not configured", entry.Level.String())
			}
			return err
		}),
	))
}
