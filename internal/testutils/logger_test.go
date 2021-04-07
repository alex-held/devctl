package testutils

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	. "github.com/alex-held/devctl/pkg/logging"
)

func TestLogger_Captures_LogMessages(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Logger", func() {
		var logger Log

		g.Describe("GIVEN the go test -v (verbose) flag has not been set", func() {
			var b *bytes.Buffer

			g.JustBeforeEach(func() {
				ToggleTestVerbose(false)
				b = &bytes.Buffer{}
				logger = NewLogger(
					WithBuffer(b),
					WithLevel(LogLevelWarn),
				)
			})

			g.It("WHEN loglevel is info --> THEN no log messages get captured", func() {
				logger.Infof("To learn more about this cli visit our docs.")
				output := b.String()
				Expect(output).To(BeEmpty())
			})

			g.It("WHEN loglevel is debug --> THEN no log messages get captured", func() {
				logger.Debugf("To learn more about this cli visit our docs.")
				output := b.String()
				Expect(output).To(BeEmpty())
			})

			g.It("WHEN loglevel is warn --> THEN captures log message", func() {
				logger.Warnf("To learn more about this cli visit our docs.")
				output := b.String()
				Expect(output).To(ContainSubstring("To learn more about this cli visit our docs."))
			})
		})

		g.Describe("GIVEN the go test -v (verbose) flag has been set", func() {
			var b *bytes.Buffer

			g.JustBeforeEach(func() {
				ToggleTestVerbose(true)
				b = &bytes.Buffer{}
				logger = NewLogger(WithBuffer(b), WithLevel(LogLevelDebug))
			})

			g.It("WHEN loglevel is info --> THEN capture log messages", func() {
				logger.Infof("To learn more about this cli visit our docs. devctl-path=%s; temp-dir-path=%s;\n", os.ExpandEnv("$HOME/.devctl"), t.TempDir())
				output := b.String()
				Expect(output).To(ContainSubstring("To learn more about this cli visit our docs."))
				//	Expect(b.Len()).To(Equal(logger.(logging.Log).Output.Len()))
			})

			g.It("WHEN loglevel is debug --> THEN capture log messages", func() {
				logger.Debugf("To learn more about this cli visit our docs. devctl-path=%s; temp-dir-path=%s\n", os.ExpandEnv("$HOME/.devctl"), t.TempDir())
				output := b.String()
				Expect(output).To(ContainSubstring("To learn more about this cli visit our docs."))
			})

			g.It("WHEN loglevel is warn --> THEN captures log messages", func() {
				logger.Warnf("To learn more about this cli visit our docs. devctl-path=%s; temp-dir-path=%s\n", os.ExpandEnv("$HOME/.devctl"), t.TempDir())
				output := b.String()
				Expect(output).To(ContainSubstring("To learn more about this cli visit our docs."))
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

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Logger", func() {
		var logger Log
		var b *bytes.Buffer
		err := errors.New("fail this test")

		g.JustBeforeEach(func() {
			b = &bytes.Buffer{}
			logger = NewLogger(
				WithBuffer(b),
				WithLevel(LogLevelError),
				WithExitFn(func(int) { /* ignore logger.Fatal exit codes for tests*/ }),
			)
		})

		g.It("WHEN loglevel is error --> THEN captures error output", func() {
			logger.Errorf("%v", err)
			output := b.String()
			Expect(output).To(ContainSubstring("fail this test"))
			Expect(output).To(ContainSubstring("ERROR"))
		})

		g.It("WHEN loglevel is fatal --> THEN captures error output", func() {
			logger.Fatalf("%s", errors.Wrapf(err, "wrapping the failing error"))
			output := b.String()
			Expect(output).Should(ContainSubstring("FATAL"))
			Expect(output).Should(ContainSubstring("wrapping the failing error: fail this test"))
		})
	})
}
