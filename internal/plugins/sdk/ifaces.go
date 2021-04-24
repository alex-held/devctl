package sdk

import (
	"context"
	"fmt"
	"path"
	"plugin"
	"reflect"

	"github.com/alex-held/devctl/pkg/devctlpath"
	"github.com/spf13/afero"
)

type SDKPlugin interface {
	Name() string
	Install(ctx context.Context, args []string) error
	Download(ctx context.Context, args []string) error
	List(ctx context.Context, args []string) error
	Current(ctx context.Context, args []string) error
	Use(ctx context.Context, args []string) error
}
type PluginFn func(context.Context, []string) error

type wrapper struct {
	plug       plugin.Plugin
	NameFn     string
	InstallFn  PluginFn
	DownloadFn PluginFn
	ListFn     PluginFn
	CurrentFn  PluginFn
	UseFn      PluginFn
}

func (w *wrapper) Name() string                                      { return w.Name() }
func (w *wrapper) Install(ctx context.Context, args []string) error  { return w.InstallFn(ctx, args) }
func (w *wrapper) Download(ctx context.Context, args []string) error { return w.DownloadFn(ctx, args) }
func (w *wrapper) List(ctx context.Context, args []string) error     { return w.ListFn(ctx, args) }
func (w *wrapper) Current(ctx context.Context, args []string) error  { return w.CurrentFn(ctx, args) }
func (w *wrapper) Use(ctx context.Context, args []string) error      { return w.UseFn(ctx, args) }

type registry struct {
	pather devctlpath.Pather
	fs     afero.Fs
}

func (r *registry) LoadSDKPlugins() (plugins []SDKPlugin, err error) {
	dir := r.pather.ConfigRoot("plugins", "sdk")
	infos, err := afero.ReadDir(r.fs, dir)
	if err != nil {
		return plugins, err
	}

	for _, fi := range infos {
		filename := fi.Name()
		path := path.Join(dir, filename)
		p, err := plugin.Open(path)
		_ = p
		if err != nil {
			return plugins, err
		}

	}
	return plugins, err
}

type symbolSetter func(*wrapper, interface{})

func (w *wrapper) Lookup(name string, asType reflect.Type, setter symbolSetter) (err error) {
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

var pluginFnType = reflect.TypeOf((PluginFn)(nil))

func LookupSDKPlugin(p *plugin.Plugin) (_ SDKPlugin, errs []error) {
	fmt.Printf("Plugin: %+v\n", *p)
	w := &wrapper{
		plug: *p,
	}

	if err := w.Lookup("Install", pluginFnType, func(w *wrapper, v interface{}) {
		as, ok := v.(func(context.Context, []string) error)
		if !ok {
			return
		}
		w.InstallFn = as
	}); err != nil {
		errs = append(errs, err)
	}

	lookupFns := []func() error{
		func() error { return w.LookupPluginFn("Download") },
		func() error { return w.LookupPluginFn("List") },
	}

	for i, lookupFn := range lookupFns {
		err := lookupFn()
		if err != nil {
			fmt.Printf("failed to lookup PluginFn; i=%d;err=%v\n", i, err)
			errs = append(errs, err)
		}
	}

	return w, errs
}

func (w *wrapper) LookupPluginFn(name string) (err error) {
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
