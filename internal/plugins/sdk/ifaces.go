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

type SdkPluginW struct {
	Plug       plugin.Plugin
	NameFn     func() string
	InstallFn  PluginFn
	DownloadFn PluginFn
	ListFn     PluginFn
	CurrentFn  PluginFn
	UseFn      PluginFn
	OutFn      func(writer io.Writer) error
}

func (w *SdkPluginW) GetPlugin() *plugin.Plugin { return &w.Plug }
func (w *SdkPluginW) Lookup(symName string) (sym plugin.Symbol, err error) {
	fmt.Printf("looking up symbol '%s'\n", symName)
	sym, err = w.Plug.Lookup(symName)
	fmt.Printf("Symbol for '%s':\t%#v\nErr:%+v\n", symName, sym, err)
	return sym, err
}

func (w *SdkPluginW) PluginName() string                               { return w.NameFn() }
func (w *SdkPluginW) Install(ctx context.Context, args []string) error { return w.InstallFn(ctx, args) }
func (w *SdkPluginW) Download(ctx context.Context, args []string) error {
	return w.DownloadFn(ctx, args)
}
func (w *SdkPluginW) List(ctx context.Context, args []string) error    { return w.ListFn(ctx, args) }
func (w *SdkPluginW) Current(ctx context.Context, args []string) error { return w.CurrentFn(ctx, args) }
func (w *SdkPluginW) Use(ctx context.Context, args []string) error     { return w.UseFn(ctx, args) }
func (w *SdkPluginW) SetStdout(writer io.Writer) error                 { return w.OutFn(writer) }

type PluginGetter interface {
	GetPlugin() *plugin.Plugin
}

type SymbolLoader interface {
	Lookup(symName string) (plugin.Symbol, error)
}

type SDKPluginLoaderFn func(elfModule string) (plug PluginGetter, err error)
type SDKPluginLoader interface {
	LoadSDKPlugin(elfModule string) (plug PluginGetter, err error)
}

func LoadSDKPlugin(elfModule string) (plug PluginGetter, err error) {
	var fn *SDKPluginLoaderFn
	return fn.LoadSDKPlugin(elfModule)
}

type pluginGetter struct {
	Plug *plugin.Plugin
}

func (p *pluginGetter) GetPlugin() (plug *plugin.Plugin) {
	return p.Plug
}

func (loaderFnP *SDKPluginLoaderFn) LoadSDKPlugin(elfModule string) (plug PluginGetter, err error) {
	p, err := plugin.Open(elfModule)
	pg := &pluginGetter{Plug: p}
	return pg, err
}

type SDKPluginBinder interface {
	Bind(p *plugin.Plugin) (plugin SDKPlugin, err error)
}

// SDKPluginBinderFn already implements interface SDKPluginBinder
type SDKPluginBinderFn func(p *plugin.Plugin) (plugin SDKPlugin, err error)

func Bind(plug *plugin.Plugin) (sdkPlugin SDKPlugin, err error) {
	var fn *SDKPluginBinderFn
	sdkPlugin, err = fn.Bind(plug)
	return sdkPlugin, err
}

func (binderFnP *SDKPluginBinderFn) Bind(plug *plugin.Plugin) (sdkPlugin SDKPlugin, err error) {
	if binderFnP != nil {
		return (*binderFnP)(plug)
	}

	sdkPluginW := &SdkPluginW{Plug: *plug}
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
