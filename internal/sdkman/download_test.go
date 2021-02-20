package sdkman

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadServiceResolve(t *testing.T) {
	client := NewSdkManClient()
	dlPath := client.Download.Resolve()("scala", "2.13.4")
	assert.Equal(t, "/Users/dev/.devctl/downloads/scala/2.13.4", dlPath)
}
