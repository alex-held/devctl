package plugins

import (
	"bytes"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
)

func newTestEngine() *Engine {
	return NewEngine(func(c *Config) *Config {
		c.Pather = nil // we don't want really to access the file system in this test
		return c
	})
}

func TestLoadManifest(t *testing.T) {
	sut := newTestEngine()
	manifest, err := sut.LoadPlugin("testdata/plugin1/plugin.yaml")
	assert.NoError(t, err)

	assert.Equal(t, manifest.Version, "v1.2.3")
	assert.Equal(t, manifest.PluginSpec.Name, "plugin1")
}

func TestExecPlugin(t *testing.T) {
	pluginRootPath := "testdata/example-exec-plugin"
	out := &bytes.Buffer{}
	sut := NewEngine(func(c *Config) *Config {
		c.Out = out
		c.Pather = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			return "/home/user/.devctl"
		}))
		c.Fs = afero.NewOsFs()
		return c
	})

	p, err := sut.LoadPlugin(path.Join(pluginRootPath, "plugin.yaml"))
	assert.NoError(t, err)

	cfg := map[string]interface{}{
		"InstallPath": "/home/user/.devctl/sdks/go",
	}

	args := []string{"use", "1.16.8"}
	execP, err := sut.NewExecutablePlugin(p, cfg)
	assert.NoError(t, err)

	result := execP.Exec(args)

	assert.Equal(t, "cfg.InstallPath=/home/user/.devctl/sdks/go\nargs[0]=use\nargs[1]=1.16.8\n", out.String())
	assert.Equal(t, "exec done", result.Error())

	t.Log(out.String())
}
