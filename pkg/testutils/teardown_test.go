package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeardownCombine_Combination_Contains_All_Teardowns(t *testing.T) {
	firstCalled := false
	secondCalled := false

	var first, second Teardown = func() {
		firstCalled = true
	}, func() {
		secondCalled = true
	}

	combination := first.CombineInto(second)
	combination()

	assert.True(t, firstCalled)
	assert.True(t, secondCalled)
}
