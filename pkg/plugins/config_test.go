package plugins

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
)

func TestLoadManifests(t *testing.T) {
	testFs := afero.NewBasePathFs(afero.NewOsFs(), path.Join("testdata", "test-load-manifests-fs"))
	_, err := testFs.Stat("/home/user/.devctl/plugins/plugin1/plugin.yaml")
	if err != nil {
		t.Fatalf("/home/user/.devctl/plugins/plugin1/plugin.yaml does not exists!")
	}

	sut := NewEngine(func(c *Config) *Config {
		c.Pather = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			return "/home/user/.devctl"
		}))
		c.Fs = testFs
		return c
	})

	manifests := sut.LoadPlugins()
	assert.Len(t, manifests, 2)
	assert.Equal(t, manifests[0].Name, "plugin1")
	assert.Equal(t, manifests[1].Name, "plugin2")
}

func TestResolveDynamic(t *testing.T) {
	e := NewEngine(func(c *Config) *Config {
		c.Pather = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			return "/home/user/.devctl"
		}))
		return c
	})

	configSpec := ConfigSpec{
		Values: map[string]string{
			"version": "1.0.0",
		},
		Static: map[string]interface{}{
			"test": "test_value",
			"data": struct {
				Name string
			}{
				Name: "myname",
			},
		},
		Dynamic: map[string]string{
			"home":         "{{ .DEVCTL_PATH_USERHOME }}",
			"root_path":    "{{ .DEVCTL_PATH_ROOT }}",
			"install_path": "{{ .DEVCTL_PATH_SDK }}",
			"version":      "{{ .values.version }}",
		},
	}

	resolved, err := ResolveDynamic(configSpec, e)
	assert.NoError(t, err)

	h, _ := os.UserHomeDir()
	assert.Equal(t, h, resolved.Dynamic["home"])
	assert.Equal(t, e.cfg.Pather.ConfigRoot(), resolved.Dynamic["root_path"])
	assert.Equal(t, e.cfg.Pather.SDK(), resolved.Dynamic["install_path"])
	assert.Equal(t, "1.0.0", resolved.Dynamic["version"])
}

func TestMarshalManifest(t *testing.T) {
	expected := "version: v1.2.3\nplugin:\n    cmd: plug\n    help: |-\n        USAGE\n              plug <flags> [subcommand]\n    subcommands:\n        - cmd: list\n          help: |-\n            USAGE\n                      plug <flags> list\n        - cmd: view\n          help: |-\n            USAGE\n                      plug <flags> view [name]\n    name: plugin1\n    pkg: plugin_1\n"
	manifest := Manifest{
		Version: "v1.2.3",
		PluginSpec: PluginSpec{
			Name: "plugin1",
			Pkg:  "plugin_1",
			CommandSpec: &CommandSpec{
				Cmd:  "plug",
				Help: "USAGE\n      plug <flags> [subcommand]",
				Subcommands: []CommandSpec{
					{
						Cmd:         "list",
						Help:        "USAGE\n          plug <flags> list",
						Subcommands: nil,
					},
					{
						Cmd:         "view",
						Help:        "USAGE\n          plug <flags> view [name]",
						Subcommands: nil,
					}},
			},
		},
	}

	b, err := yaml.Marshal(&manifest)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(b))
}

func TestUnmarshalManifest(t *testing.T) {
	testFs := afero.NewBasePathFs(afero.NewOsFs(), path.Join("testdata", "test-load-manifests-fs"))
	_, err := testFs.Stat("/home/user/.devctl/plugins/plugin1/plugin.yaml")
	if err != nil {
		t.Fatalf("/home/user/.devctl/plugins/plugin1/plugin.yaml does not exists!")
	}

	sut := NewEngine(func(c *Config) *Config {
		c.Pather = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			return "/home/user/.devctl"
		}))
		c.Fs = testFs
		return c
	})

	manifests := sut.LoadPlugins()
	assert.Len(t, manifests, 2)
	assert.Equal(t, manifests[0].Name, "plugin1")
	assert.Equal(t, manifests[1].Name, "plugin2")
}
