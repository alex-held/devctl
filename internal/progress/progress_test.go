package progress

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProgress(t *testing.T) {
	pr := progressReporter{}
	err := pr.report(100, 50)
	assert.Equal(t, nil, err)
}
