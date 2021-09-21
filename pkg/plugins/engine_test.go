package plugins

import (
	"bytes"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func newTestEngine() *Engine {
	return NewEngine(func(c *Config) *Config {
		c.Pather = nil // we don't want really to access the file system in this test
		return c
	})
}

func TestLoadManifest(t *testing.T) {
	sut := newTestEngine()
	manifest, err := sut.LoadPlugin("testdata/plugin1.yaml")
	assert.NoError(t, err)

	assert.Equal(t, manifest.Version, "v1.2.3")
	assert.Equal(t, manifest.PluginSpec.Name, "plugin1")
}

func TestExecutePlugin(t *testing.T) {
	pluginRootPath := "testdata/example-plugin"
	out := &bytes.Buffer{}
	sut := NewEngine(func(c *Config) *Config {
		c.Out = out
		c.Pather = nil
		c.Fs = afero.NewOsFs()
		return c
	})

	manifest, err := sut.LoadPlugin(path.Join(pluginRootPath, "plugin.yaml"))
	assert.NoError(t, err)

	assert.Equal(t, manifest.Version, "v1.2.3")
	assert.Equal(t, manifest.PluginSpec.Name, "example-plugin")
	assert.Equal(t, manifest.PluginSpec.Pkg, "example_plugin")

	err = sut.Execute("example-plugin", []string{"hello world"})
	assert.Equal(t, "some error", err.Error())

	expected := "fmt\n[INFO]\t  New called with args=[hello world]\n[DEBUG]\t  Example Plugin name=New()\n"
	actual := out.String()
	assert.Equal(t, expected, actual)
}

func TestExecute(t *testing.T) {
	buf := &bytes.Buffer{}

	sut := NewEngine(func(c *Config) *Config {
		c.Pather = nil
		c.Out = buf
		return c
	})

	sut.pluginCache["go"] = &Plugin{
		Source: `package plugin_go
		import "fmt"
		func New(args []string)(error){
			fmt.Printf("arglen=%d", len(args))
		    return nil
		}`,
		RootPath: "plugin-go",
		Manifest: &Manifest{
			Version: "v1.0.0",
			PluginSpec: PluginSpec{
				Name: "go",
				CommandSpec: &CommandSpec{
					Cmd:  "go",
					Help: "USAGE\n go [args]",
				},
				Pkg: "plugin_go",
			},
		},
	}

	err := sut.Execute("go", []string{"hello", "world"})
	assert.NoError(t, err)

	actualOut := buf.String()
	assert.Equal(t, "arglen=2", actualOut)

	println("")
	println("--- actual out ---")
	println(actualOut)
	println("------------------")
}
