package testutils

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"

	"github.com/alex-held/devctl/internal/logging"
)

func TestLogger_Captures_LogMessages(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Logger", func() {
		RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
		var logger *logging.Logger

		g.Describe("GIVEN the go test -v (verbose) flag has not been set", func() {
			g.JustBeforeEach(func() {
				ToggleTestVerbose(false)
				logger = logging.NewLogger()
			})

			g.It("WHEN loglevel is info --> THEN no log messages get captured", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Infof("To learn more about this cli visit our docs.")

				Expect(logger.Output.String()).To(BeEmpty())
			})

			g.It("WHEN loglevel is debug --> THEN no log messages get captured", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Debugf("To learn more about this cli visit our docs.")
				Expect(logger.Output.String()).To(BeEmpty())
			})

			g.It("WHEN loglevel is warn --> THEN captures log message", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Warnf("To learn more about this cli visit our docs.")
				Expect(logger.Output.String()).To(ContainSubstring("To learn more about this cli visit our docs."))
			})
		})

		g.Describe("GIVEN the go test -v (verbose) flag has been set", func() {
			g.JustBeforeEach(func() {
				ToggleTestVerbose(true)
				logger = NewLogger()
			})

			g.It("WHEN loglevel is info --> THEN capture log messages", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Debugln("To learn more about this cli visit our docs.")
				Expect(logger.Output.String()).To(ContainSubstring("To learn more about this cli visit our docs."))
				Expect(logger.Output.Len()).To(Equal(logger.Output.Len()))
			})

			g.It("WHEN loglevel is debug --> THEN capture log messages", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Debugf("To learn more about this cli visit our docs.")
				Expect(logger.Output.String()).To(ContainSubstring("To learn more about this cli visit our docs."))
				Expect(logger.Output.Len()).To(Equal(logger.Output.Len()))
			})

			g.It("WHEN loglevel is warn --> THEN captures log messages", func() {
				logger.
					WithField("devctl-path", os.ExpandEnv("$HOME/.devctl")).
					WithField("temp-dir-path", t.TempDir()).
					Warnf("To learn more about this cli visit our docs.")
				Expect(logger.Output.String()).To(ContainSubstring("To learn more about this cli visit our docs."))
				Expect(logger.Output.Len()).To(Equal(logger.Output.Len()))
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
		var logger *logging.Logger
		err := errors.New("fail this test")

		g.JustBeforeEach(func() {
			logger = NewLogger()
			logger.ExitFunc = func(int) { /* ignore logger.Fatal exit codes for tests*/ }
		})

		g.It("WHEN loglevel is error --> THEN captures error output", func() {
			logger.WithError(err).Error()
			Expect(logger.Error.String()).To(MatchRegexp("time=.*\\slevel=error error=.*fail this test\""))
			Expect(logger.Error.Len()).To(Equal(logger.Error.Len()))
		})

		g.It("WHEN loglevel is fatal --> THEN captures error output", func() {
			logger.WithError(err).Fatal()
			Expect(logger.Error.String()).To(MatchRegexp("time=.*\\slevel=fatal error=.*fail this test\""))
		})
	})
}
