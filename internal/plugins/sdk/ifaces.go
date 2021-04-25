package sdk

import (
	"context"
	"fmt"
	"io"
	"plugin"
	"reflect"
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
	Plug       plugin.Plugin
	NameFn     func() string
	InstallFn  PluginFn
	DownloadFn PluginFn
	ListFn     PluginFn
	CurrentFn  PluginFn
	UseFn      PluginFn
	OutFn      func(writer io.Writer) error
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

type SDKPluginLoaderFn func(elfModule string) (plug *plugin.Plugin, err error)
type SDKPluginLoader interface {
	LoadSDKPlugin(elfModule string) (plug *plugin.Plugin, err error)
}

func (loaderFnP SDKPluginLoaderFn) LoadSDKPlugin(elfModule string) (plug *plugin.Plugin, err error) {
	fmt.Printf("Executing %T nil value\n", (SDKPluginLoader)(nil))
	plug, err = plugin.Open(elfModule)
	return plug, err
}

type SDKPluginBinder interface {
	Bind(p *plugin.Plugin) (plugin SDKPlugin, err error)
}

//SDKPluginBinderFn already implements interface SDKPluginBinder
type SDKPluginBinderFn func(p *plugin.Plugin) (plugin SDKPlugin, err error)

func (binderFnP *SDKPluginBinderFn) Bind(plug *plugin.Plugin) (sdkPlugin SDKPlugin, err error) {
	if binderFn := *binderFnP; binderFn != nil {
		return binderFn(plug)
	}

	sdkPluginW := &sdkPluginW{Plug: *plug}
	bindPluginFn := func(name string) (err error) {
		symbol, err := sdkPluginW.Plug.Lookup(name)
		switch fn := symbol.(type) {
		case func(context.Context, []string) error:
			pluginFn := PluginFn(fn)
			pluginFnV := reflect.ValueOf(pluginFn)
			wrapperV := reflect.ValueOf(sdkPluginW).Elem()
			fieldV := wrapperV.FieldByName(name + "Fn")
			fieldV.Set(pluginFnV)
			return nil
		default:
			return fmt.Errorf("Could not load symbol '%s' as type 'func(context.Context, []string) error'\nSymbol: %+v\n\n", name, symbol)
		}
	}

	if err = bindPluginFn("Download"); err != nil {
		return sdkPluginW, err
	}
	if err = bindPluginFn("Install"); err != nil {
		return sdkPluginW, err
	}
	if err = bindPluginFn("Use"); err != nil {
		return sdkPluginW, err
	}
	if err = bindPluginFn("List"); err != nil {
		return sdkPluginW, err
	}
	if err = bindPluginFn("Install"); err != nil {
		return sdkPluginW, err
	}

	symbol, err := sdkPluginW.Plug.Lookup("SetStdout")
	switch fn := symbol.(type) {
	case func(writer io.Writer) error:
		sdkPluginW.OutFn = fn
	default:
		return sdkPluginW, fmt.Errorf("Could not load symbol '%s' as type 'func(io.Writer) error'\nSymbol: %+v\n\n", "SetStdout", symbol)
	}

	symbol, err = sdkPluginW.Plug.Lookup("PluginName")
	switch fn := symbol.(type) {
	case func() string:
		sdkPluginW.NameFn = fn
	default:
		return sdkPluginW, fmt.Errorf("Could not load symbol '%s' as type 'func() string'\nSymbol: %+v\n\n", "PluginName", symbol)
	}
	return sdkPluginW, nil
}
