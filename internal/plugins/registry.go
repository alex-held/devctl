package plugins

import (
	"context"
	"os"
	"path"
	"reflect"

	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/plugins/sdk"
	"github.com/alex-held/devctl/pkg/devctlpath"
)

// GlobalRegistry is a v variable
var GlobalRegistry Registry

func init() {
	var globalRegistry = registry{}
	GlobalRegistry = globalRegistry
}

type PluginFn func(context.Context, []string) error

var PluginFnType = reflect.TypeOf((PluginFn)(nil))

type PluginProvider func() map[string]Plugin

func (m *pluginManager) LoadSDKPlugins() (err error) {

	dir := m.pather.ConfigRoot("plugins")
	files, err := afero.ReadDir(m.fs, dir)
	if err != nil {
		return err
	}
	for _, fi := range files {
		pluginPath := path.Join(dir, fi.Name())
		if m.registry.IsRegistered(pluginPath) {
			continue
		}
		if !fi.IsDir() && fi.Mode()&os.ModeType == 0 {
			unboundPlugin, err := m.sdkLoaderFn.LoadSDKPlugin(pluginPath)
			if err != nil {
				return err
			}
			sdkPlugin, err := m.sdkBinderFn.Bind(unboundPlugin)
			if err != nil {
				return err
			}
			m.registry.Register(pluginPath, sdkPlugin)
		}
	}
	return nil
}

type registry map[string]Plugin

func (r *registry) register(pluginPath string, plug Plugin) {
	(*r)[pluginPath] = plug
}

func (r *registry) isRegistered(pluginPath string) bool {
	_, ok := (*r)[pluginPath]
	return !ok
}

type pluginManager struct {
	provider PluginProvider
	fs       afero.Fs
	pather   devctlpath.Pather
	registry registry

	sdkPlugins  map[string]Plugin
	sdkLoaderFn sdk.SDKPluginLoaderFn
	sdkBinderFn sdk.SDKPluginBinderFn
}

func (m *pluginManager) GetProvider() PluginProvider {
	copiedSnapShot := make(map[string]Plugin, len(m.registry))
	for k, v := range m.registry {
		copiedSnapShot[k] = v
	}
	return func() map[string]Plugin {
		return copiedSnapShot
	}
}

func (r registry) Register(pluginPath string, Plugin Plugin) {
	r[pluginPath] = Plugin
}

func (r registry) IsRegistered(pluginPath string) bool {
	_, ok := r[pluginPath]
	return !ok
}

type Plugin interface {
	PluginName() string
}

type Registry interface {
	Register(pluginPath string, plug Plugin)
	IsRegistered(pluginPath string) bool
}
