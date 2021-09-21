package golang

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/alex-held/devctl/pkg/devctlpath"
	"github.com/alex-held/devctl/pkg/plugins"
)

func TestExecuteGoPlugin(t *testing.T) {
	out := &bytes.Buffer{}
	sut := plugins.NewEngine(func(c *plugins.Config) *plugins.Config {
		c.Out = out
		c.Fs = afero.NewOsFs()
		c.Pather = devctlpath.DefaultPather()
		return c
	})

	p, err := sut.LoadPlugin("plugin.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, p)

	err = sut.Execute("go", []string{"current"})
	assert.NoError(t, err)
	assert.Equal(t, out.String(), "1.16.8")
}
