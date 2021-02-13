package testutils

import (
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"golang.org/x/exp/errors"
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
	logger := zaptest.NewLogger(t).WithOptions(zap.ErrorOutput(os.Stderr))
	
	defer logger.Sync()
	logger.Error("error", zap.Error(errors.New("Fail this test")))
	// logger.AssertLenErrorMessages(1)
	// spy.AssertFailed()
}
