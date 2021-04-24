package plugins

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"plugin"
	"reflect"

	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/plugins/sdk"
	"github.com/alex-held/devctl/pkg/devctlpath"
)

type symbolSetter func(*pluginWrapper, interface{})
type PluginFn func(context.Context, []string) error

var PluginFnType = reflect.TypeOf((PluginFn)(nil))

func (r *registry) LoadSDKPlugins() (plugins []sdk.SDKPlugin, err error) {
	dir := r.Pather.ConfigRoot("plugins")
	files, err := afero.ReadDir(r.Fs, dir)
	if err != nil {
		return r.sdkPlugins, err
	}
	for _, fi := range files {
		if !fi.IsDir() && fi.Mode()&os.ModeType == 0 {
			filepath := path.Join(dir, fi.Name())
			sdkPlugin, err := sdk.LoadSDKPlugin(filepath)
			if err != nil {
				return r.sdkPlugins, err
			}
			r.sdkPlugins = append(r.sdkPlugins, sdkPlugin)
			return r.sdkPlugins, err
		}
	}
	return r.sdkPlugins, err
}

func NewRegistry(pather devctlpath.Pather, fs afero.Fs) *registry {
	r := &registry{
		Pather: pather,
		Fs:     fs,
		Store: &store{
			Pather: pather,
			Fs:     fs,
		},
	}
	r.sdkRegistry = (sdk.SDKRegistry)(r)
	return r
}

type registry struct {
	Pather devctlpath.Pather
	Fs     afero.Fs
	Store  Store

	sdkPlugins  []sdk.SDKPlugin
	sdkRegistry sdk.SDKRegistry
}

func (r *registry) ReloadPlugins() (plugins []Plugin, err error) {
	if r.sdkPlugins, err = r.sdkRegistry.LoadSDKPlugins(); err != nil {
		return plugins, err
	}
	for _, sdkPlugin := range r.sdkPlugins {
		plugins = append(plugins, sdkPlugin)
	}
	return plugins, err
}

type Plugin interface {
	PluginName() string
}

type Registry interface {
	ReloadPlugins() (plugins []Plugin, err error)
}

type pluginWrapper struct {
	plug      plugin.Plugin
	Name      string
	LookupFns []func() error
	stdout    io.Writer
}

func (w *pluginWrapper) LookupPluginFn(name string) (err error) {
	symbol, err := w.plug.Lookup(name)
	fn, ok := symbol.(func(context.Context, []string) error)
	if !ok {
		return fmt.Errorf("Could not load symbol '%s' as type 'func(context.Context, []string) error'\nSymbol: %+v\n\n", name, symbol)
	}
	pluginFn := PluginFn(fn)
	pluginFnV := reflect.ValueOf(pluginFn)
	wrapperV := reflect.ValueOf(w).Elem()
	fieldV := wrapperV.FieldByName(name + "Fn")
	fieldV.Set(pluginFnV)
	return nil
}

func (w *pluginWrapper) lookupSymbol(name string, asType reflect.Type, setter symbolSetter) (err error) {
	sym, err := w.plug.Lookup(name)
	if err != nil {
		return err
	}
	ok := reflect.TypeOf(sym).AssignableTo(asType)
	if !ok {
		return fmt.Errorf("symbol '%+v' is not assignable to asType '%s'\n", sym, asType.Name())
	}
	setter(w, sym)
	return nil
}
