package testutils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewTestLogSpy(t testing.TB) *LogSpy {
	return &LogSpy{TB: t, failed: false, Messages: []string{}}
}

// LogSpy is a testing.TB that captures logged messages.
type LogSpy struct {
	testing.TB
	failed        bool
	Messages      []string
	ErrorMessages []errorEntry
}

type errorEntry struct {
	Error   error
	Message string
}

func (e errorEntry) String() string {
	if e.Error != nil {
		return e.Error.Error()
	}
	return e.Message
}

func (t *LogSpy) Fail() {
	t.failed = true
	t.TB.Fail()
}

func (t *LogSpy) Failed() bool {
	return t.failed || t.TB.Failed()
}

func (t *LogSpy) FailNow() {
	t.failed = true
	t.TB.FailNow()
}

// Log messages are in the format,
//
//   2017-10-27T13:03:01.000-0700	DEBUG	your message here	{data here}
//
// We strip the first part of these messages because we can't really test
// for the timestamp from these tests.
func (t *LogSpy) Logf(format string, args ...interface{}) {
	m := fmt.Sprintf(format, args...)
	if strings.Contains(m, "[test.error]") {
		errorEntry := errorEntry{Error: fmt.Errorf(format, args...)}
		t.TB.Errorf(errorEntry.String())
		t.ErrorMessages = append(t.ErrorMessages, errorEntry)
	}
	m = m[strings.IndexByte(m, '\t')+1:]
	t.Messages = append(t.Messages, m)
	t.TB.Log(m)
}

func (t *LogSpy) AssertMessages(messages ...string) {
	assert.Equal(t.TB, messages, t.Messages, "logged messages did not match")
}

func (t *LogSpy) AssertPassed() {
	t.assertFailed(false, "expected test to pass")
}

func (t *LogSpy) AssertFailed() {
	t.assertFailed(true, "expected test to fail")
}

func (t *LogSpy) assertFailed(v bool, msg string) {
	assert.Equal(t.TB, v, t.failed, msg)
}

func (t *LogSpy) AssertLenErrorMessages(i int) {
	assert.Equal(t.TB, i, len(t.ErrorMessages))
}
