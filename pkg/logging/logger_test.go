package logging

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFormatLogMessage(t *testing.T) {
	log := logrus.New()
	b := &bytes.Buffer{}
	log.Out = io.MultiWriter(os.Stdout, b)
	log.Formatter = &DevCtlFormatter{}
	log.Level = logrus.TraceLevel
	log.Info("Log messager here ✅")

	expected := "\x1b[;2m[ \x1b[0m\x1b[38;2;23;175;225;1;4mINFO\x1b[0m\x1b[;2m ]\x1b[0m\x1b[;2m\t\x1b[0mLog messager here ✅\n"

	assert.Equal(t, expected, b.String())
}

func TestColorOutput(t *testing.T) {
	type message struct {
		Level    LogLevel
		Message  string
		Expected string
	}

	var logMsgs = []message{
		{
			Level:    LogLevelDebug,
			Message:  "The file /Users/dev/go/src/github.com/alex-held/devctl/pkg.golang.org does not exist.",
			Expected: "\x1b[;2m[ \x1b[0m\x1b[38;2;23;175;225;1;4mDEBUG\x1b[0m\x1b[;2m ]\x1b[0m\x1b[;2m\t\x1b[0mThe file /Users/dev/go/src/github.com/alex-held/devctl/pkg.golang.org does not exist.\n",
		},
		{
			Level:    LogLevelInfo,
			Message:  "insert into the template",
			Expected: "\x1b[;2m[ \x1b[0m\x1b[38;2;23;175;225;1;4mINFO\x1b[0m\x1b[;2m ]\x1b[0m\x1b[;2m\t\x1b[0minsert into the template\n",
		},
		{
			Level:    LogLevelWarn,
			Message:  "Could not find 'sdk/current' PATH",
			Expected: "\x1b[;2m[ \x1b[0m\x1b[38;2;249;254;0;1;4mWARNING\x1b[0m\x1b[;2m ]\x1b[0m\x1b[;2m\t\x1b[0mCould not find 'sdk/current' PATH\n",
		},
		{
			Level:    LogLevelError,
			Message:  "Prepare some data to insert into the template.",
			Expected: "\x1b[;2m[ \x1b[0m\x1b[38;2;255;60;109;1;4mERROR\x1b[0m\x1b[;2m ]\x1b[0m\x1b[;2m\t\x1b[0mPrepare some data to insert into the template.\n",
		},
	}

	for _, msg := range logMsgs {
		b := &bytes.Buffer{}
		log := NewLogger(WithBuffer(b), WithLevel(LogLevelDebug))

		log.Logf(msg.Level, msg.Message)

		actual := b.String()
		assert.Equal(t, msg.Expected, actual)
	}
}

func TestNewLogger(t *testing.T) {
	b := &bytes.Buffer{}

	l := NewLogger(
		WithFormatter(&DevCtlFormatter{}),
		WithBuffer(b),
		WithName("TestNewLogger"),
		WithLevel(LogLevelDebug),
	)

	l.Infof("hello world!")

	actual := b.String()
	assert.Contains(t, actual, "hello world!")
}
