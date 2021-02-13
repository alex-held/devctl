package testutils

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestLogger_Captures_LogMessages(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Logger", func() {

		RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
		var out *bytes.Buffer
		var logger *logrus.Logger

		g.Describe("GIVEN the go test -v (verbose) flag has not been set", func() {

			g.JustBeforeEach(func() {
				ToggleTestVerbose(false)
				out = bytes.NewBuffer(nil)
				logger = NewLogger(out)
			})

			g.It("WHEN loglevel is info --> THEN no log messages get captured", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Infof("To learn more about this cli visit our docs.")
				Expect(out.String()).To(BeEmpty())
				Expect(out.Len()).To(Equal(0))
			})

			g.It("WHEN loglevel is debug --> THEN no log messages get captured", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Debugf("To learn more about this cli visit our docs.")
				Expect(out.String()).To(BeEmpty())
				Expect(out.Len()).To(Equal(0))
			})

			g.It("WHEN loglevel is warn --> THEN captures log message", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Warnf("To learn more about this cli visit our docs.")
				Expect(out.String()).To(ContainSubstring("To learn more about this cli visit our docs."))
				Expect(out.Len()).To(Equal(out.Len()))
			})
		})

		g.Describe("GIVEN the go test -v (verbose) flag has been set", func() {

			g.JustBeforeEach(func() {
				ToggleTestVerbose(true)
				out = bytes.NewBuffer(nil)
				logger = NewLogger(out)
			})

			g.It("WHEN loglevel is info --> THEN capture log messages", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Infof("To learn more about this cli visit our docs.")
				str := out.String()
				println(str)
				Expect(out.String()).To(ContainSubstring("To learn more about this cli visit our docs."))
				Expect(out.Len()).To(Equal(out.Len()))
			})

			g.It("WHEN loglevel is debug --> THEN capture log messages", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Debugf("To learn more about this cli visit our docs.")
				Expect(out.String()).To(ContainSubstring("To learn more about this cli visit our docs."))
				Expect(out.Len()).To(Equal(out.Len()))
			})

			g.It("WHEN loglevel is warn --> THEN captures log messages", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Warnf("To learn more about this cli visit our docs.")
				Expect(out.String()).To(ContainSubstring("To learn more about this cli visit our docs."))
				Expect(out.Len()).To(Equal(out.Len()))
			})
		})

	})
}

func ToggleTestVerbose(on bool) {
	v := flag.CommandLine.Lookup("test.v")
	err := v.Value.Set(fmt.Sprintf("%t", on))
	if err != nil {
		panic(err)
	}
}

func TestDefaultLogger_Captures_FailedState(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Logger", func() {
		RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
		var buf *bytes.Buffer
		var logger *logrus.Logger
		err := errors.New("fail this test")

		g.JustBeforeEach(func() {
			buf = bytes.NewBuffer(nil)
			logger = NewLogger(buf)
			logger.ExitFunc = func(int) { /* ignore logger.Fatal exit codes for tests*/ }
		})

		g.It("WHEN loglevel is error --> THEN captures error output", func() {
			logger.WithError(err).Error()
			Expect(buf.String()).To(MatchRegexp("time=.*\\slevel=error error=.*fail this test\""))
			Expect(buf.Len()).To(Equal(buf.Len()))
		})

		g.It("WHEN loglevel is fatal --> THEN captures error output", func() {
			logger.WithError(err).Fatal()
			Expect(buf.String()).To(MatchRegexp("time=.*\\slevel=fatal error=.*fail this test\""))
			Expect(buf.Len()).To(Equal(buf.Len()))
		})

	})
}
