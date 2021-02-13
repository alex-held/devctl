package testutils

import (
	"bytes"
	"errors"
	"regexp"
	"testing"
	
	"github.com/franela/goblin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTestLogger_Can_Spy_On_Infof_Messages(t *testing.T) {
	spy, teardown := SetupTestLogger(t)
	z := zap.S()
	defer teardown()
	
	z.Infof("This is a test, containging a map: %+v", map[string]interface{}{"value1": 1, "value2": 15})
	
	assert.Len(t, spy.Messages, 1)
}

func TestTestSpy_Captures_FailedState(t *testing.T) {
	spy, teardown := SetupTestLogger(t)
	z := zap.S()
	defer teardown()
	z.Errorw("error", errors.New("Fail this test"))
	spy.AssertLenErrorMessages(1)
	spy.AssertFailed()
}

func TestDefaultLogger_Captures_FailedState(t *testing.T) {
	g := goblin.Goblin(t)
	
	g.Describe("Logger", func() {
		g.It("captures error output", func() {
			err := errors.New("fail this test")
			buf := new(bytes.Buffer)
			logger := logrus.New()
			logger.SetOutput(buf)
			
			logger.WithError(err).Error()
			
			g.Assert(buf.String()).IsNotZero("there should be something in the stderr")
			assert.Regexp(g, regexp.MustCompile("time=.*\\slevel=error error=.*fail this test\""), buf.String())
		})
	})
}
