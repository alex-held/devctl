package progress

import (
	"testing"
	"github.com/stretchr/testify/assert"
	)

func TestProgress(t *testing.T) {
	pr := progressReporter{}
	err := pr.report( 100, 50)
	assert.Equal(t, nil, err)
}
