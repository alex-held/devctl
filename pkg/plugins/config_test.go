package plugins

import (
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/alex-held/devctl/pkg/devctlpath"
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

func TestLoadMani(t *testing.T) {
	sut := NewEngine(func(c *Config) *Config {
		c.Pather = devctlpath.DefaultPather()
		c.Fs = afero.NewOsFs()
		return c
	})

	manifest, err := sut.LoadPlugin("testdata/plugin1.yaml")
	assert.NoError(t, err)

	// plugin
	assert.Equal(t, manifest.Version, "v1.2.3")
	assert.Equal(t, manifest.Name, "plugin1")
	assert.Equal(t, manifest.Pkg, "plugin_1")

	// cmd
	assert.Equal(t, manifest.Cmd, "plug")
	assert.Equal(t, manifest.Help, "USAGE\n      plug <flags> [subcommand]")

	// subcommands
	assert.Len(t, manifest.Subcommands, 2)
	assert.Equal(t, manifest.Subcommands[0], CommandSpec{
		Cmd:         "list",
		Help:        "USAGE\n          plug <flags> list",
		Subcommands: nil,
	})
	assert.Equal(t, manifest.Subcommands[1], CommandSpec{
		Cmd:         "view",
		Help:        "USAGE\n          plug <flags> view [name]",
		Subcommands: nil,
	})
}

func TestMarshalManifest(t *testing.T) {
	expected := "version: v1.2.3\nplugin:\n    cmd: plug\n    help: |-\n        USAGE\n              plug <flags> [subcommand]\n    subcommands:\n        - cmd: list\n          help: |-\n            USAGE\n                      plug <flags> list\n        - cmd: view\n          help: |-\n            USAGE\n                      plug <flags> view [name]\n    name: plugin1\n    pkg: plugin_1"
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
