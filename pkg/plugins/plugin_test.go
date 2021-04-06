package plugins

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryLoadsPlugins(t *testing.T) {
	reg := pluginRegistry{
		SDKPlugins: []SDKPlugin{},
	}

	gopath := os.Getenv("GOPATH")
	pluginDir := path.Join(gopath, "github.com/alex-held/devctl-sdkplugin-go")

	plugins := []string{}
	err := filepath.Walk(pluginDir, func(p string, info fs.FileInfo, err error) error {
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
	if err != nil {
		t.Fatalf(err.Error())
	}
	for _, pluginPath := range plugins {
		plugin, err := reg.Load(pluginPath)
		if err != nil {
			t.Fatalf(err.Error())
		}
		assert.NotNil(t, plugin)
	}
}
