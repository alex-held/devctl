package plugins

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alex-held/devctl/pkg/plugins/sdk"
)

func TestRegistryLoadsPlugins(t *testing.T) {
	reg := pluginRegistry{
		SDKPlugins: []sdk.SDKPlugin{},
	}

	gopath := os.Getenv("GOPATH")
	pluginDir := path.Join(gopath, "github.com/alex-held/devctl-sdkplugin-go")

	plugins := []string{}
	filepath.Walk(pluginDir, func(p string, info fs.FileInfo, err error) error {
		if path.Ext(p) == ".so" {
			plugins = append(plugins, p)
			println(p)
			return nil
		}
		if err != nil {
			return err
		}
		return err
	})

	for _, plugin_path := range plugins {
		plugin, err := reg.Load(plugin_path)
		if err != nil {
			t.Fatalf(err.Error())
		}
		assert.NotNil(t, plugin)
	}
}
