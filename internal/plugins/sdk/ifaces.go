package sdk

import (
	"context"
	"fmt"
	"io"
	"plugin"
	"reflect"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type PluginFn func(context.Context, []string) error

var pluginFnType = reflect.TypeOf((PluginFn)(nil))

type SDKPlugin interface {
	SetStdout(io.Writer) error
	PluginName() string
	Install(ctx context.Context, args []string) error
	Download(ctx context.Context, args []string) error
	List(ctx context.Context, args []string) error
	Current(ctx context.Context, args []string) error
	Use(ctx context.Context, args []string) error
}

type sdkPluginW struct {
	plug       plugin.Plugin
	NameFn     func() string
	InstallFn  PluginFn
	DownloadFn PluginFn
	ListFn     PluginFn
	CurrentFn  PluginFn
	UseFn      PluginFn
	OutFn      func(writer io.Writer) error
}

type SDKRegistry interface {
	LoadSDKPlugins() (plugins []SDKPlugin, err error)
}

func LoadSDKPlugin(path string) (p SDKPlugin, err error) {
	log.Debugf("Loading plugin from path '%s'\n", path)

	plug, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	w := &sdkPluginW{
		plug: *plug,
	}

	var lookupFns = []func() error{
		func() error { return w.LookupPluginFn("Download") },
		func() error { return w.LookupPluginFn("Install") },
		func() error { return w.LookupPluginFn("Current") },
		func() error { return w.LookupPluginFn("Use") },
		func() error { return w.LookupPluginFn("List") },
		func() error {
			return w.Lookup("PluginName", reflect.TypeOf((func() string)(nil)), func(w *sdkPluginW, i interface{}) {
				w.NameFn = func() string {
					return i.(func() string)()
				}
			})
		}, func() error {
			var setStdoutSym, err = w.plug.Lookup("SetStdout")
			if err != nil {
				return fmt.Errorf("error while looking up symbol SDKPlugin.SetStdout\nErr=%+v\n", err)
			}
			if setStdoutFn, ok := setStdoutSym.(func(io.Writer) error); ok {
				w.OutFn = setStdoutFn
				return nil
			}
			return nil
		},
	}

	for i, lookupFn := range lookupFns {
		if err = lookupFn(); err != nil {
			return w, errors.Wrapf(err, "failed to lookup PluginFn; i=%d\n", i)
		}
	}
	return w, err
}

func (w *sdkPluginW) PluginName() string                               { return w.NameFn() }
func (w *sdkPluginW) Install(ctx context.Context, args []string) error { return w.InstallFn(ctx, args) }
func (w *sdkPluginW) Download(ctx context.Context, args []string) error {
	return w.DownloadFn(ctx, args)
}
func (w *sdkPluginW) List(ctx context.Context, args []string) error    { return w.ListFn(ctx, args) }
func (w *sdkPluginW) Current(ctx context.Context, args []string) error { return w.CurrentFn(ctx, args) }
func (w *sdkPluginW) Use(ctx context.Context, args []string) error     { return w.UseFn(ctx, args) }
func (w *sdkPluginW) SetStdout(writer io.Writer) error                 { return w.OutFn(writer) }

type symbolSetter func(*sdkPluginW, interface{})

func (w *sdkPluginW) Lookup(name string, asType reflect.Type, setter symbolSetter) (err error) {
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

func (w *sdkPluginW) LookupPluginFn(name string) (err error) {
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
