package plugins

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/afero"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
	"github.com/alex-held/devctl-kit/pkg/plugins"
)

type Engine struct {
	cfg         *Config
	pluginCache map[string]*Plugin
}

type Config struct {
	Out    io.Writer
	Fs     afero.Fs
	Pather devctlpath.Pather
}

type Option func(c *Config) *Config

func NewEngine(opts ...Option) *Engine {
	e := &Engine{
		cfg: &Config{
			Out:    os.Stdout,
			Fs:     afero.NewOsFs(),
			Pather: devctlpath.DefaultPather(),
		},
		pluginCache: map[string]*Plugin{},
	}
	for _, opt := range opts {
		opt(e.cfg)
	}
	return e
}

type Plugin struct {
	*Manifest
	Source   string
	RootPath string
}

var ErrNoPluginWithNameFound = fmt.Errorf("no plugin with that name could be found in the pluginCache")

// Execute tries to execute a named plugin from the cache
func (e *Engine) Execute(pluginName string, args []string) (err error) {
	if plugin, ok := e.pluginCache[pluginName]; ok {
		return e.execute(plugin, args)
	}
	return ErrNoPluginWithNameFound
}

func (e *Engine) execute(p *Plugin, args []string) (err error) {
	i := interp.New(interp.Options{
		GoPath: path.Join(p.RootPath, "_gopath"),
		Stdout: e.cfg.Out,
		Stderr: e.cfg.Out,
	})

	// imports
	i.Use(interp.Symbols)
	i.Use(syscall.Symbols)
	i.Use(unsafe.Symbols)
	i.Use(unrestricted.Symbols)
	i.Use(stdlib.Symbols)

	// load plugin code
	_, err = i.Eval(p.Source)
	if err != nil {
		return err
	}

	vConfig, err := i.Eval(p.Pkg + `.CreateConfig()`)
	if err != nil {
		return fmt.Errorf("failed to eval CreateConfig: %w", err)
	}

	vNewFn, err := i.Eval(p.Pkg + `.New`)
	if err != nil {
		return fmt.Errorf("failed to eval New: %w", err)
	}

	createConfigFn := vConfig.Interface().(func() interface{})
	newFn := vNewFn.Interface().(func(interface{}, []string) error)

	// execute plugin
	cfg := createConfigFn()
	fmt.Printf("CFG: %v", cfg)
	fmt.Printf("EXEC with args; args=%#v", args)

	return newFn(cfg, args)
}

func (execP *ExecutablePlugin) Exec(args []string) (err error) {
	execArgs := []reflect.Value{execP.Config, reflect.ValueOf(args)}
	result := execP.ExecFn.Call(execArgs)

	return result[0].Interface().(error)
}

//goland:noinspection GoUnhandledErrorResult
func (e *Engine) NewExecutablePlugin(p *Plugin, config map[string]interface{}) (execP *ExecutablePlugin, err error) {
	i := interp.New(interp.Options{
		GoPath: path.Join(p.RootPath, "_gopath"),
		Stdout: e.cfg.Out,
		Stderr: e.cfg.Out,
	})

	// imports
	i.Use(interp.Symbols)
	i.Use(syscall.Symbols)
	i.Use(unsafe.Symbols)
	i.Use(unrestricted.Symbols)
	i.Use(stdlib.Symbols)

	// load plugin code
	_, err = i.Eval(p.Source)
	if err != nil {
		return nil, err
	}

	vConfig, err := i.Eval(p.Pkg + `.CreateConfig()`)
	if err != nil {
		return nil, fmt.Errorf("failed to eval CreateConfig: %w", err)
	}

	fnExec, err := i.Eval(p.Pkg + `.Exec`)
	if err != nil {
		return nil, fmt.Errorf("failed to eval Exec: %w", err)
	}

	if err = e.decodeConfig(vConfig, config); err != nil {
		return nil, err
	}

	execP = &ExecutablePlugin{
		Plugin: p,
		Config: vConfig,
		ExecFn: fnExec,
	}

	return execP, nil
}

func (e *Engine) decodeConfig(vConfig reflect.Value, cfg map[string]interface{}) (err error) {
	cfg["Context"] = &plugins.Context{
		Out:     e.cfg.Out,
		Pather:  e.cfg.Pather,
		Context: context.Background(),
	}
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToSliceHookFunc(","),
		WeaklyTypedInput: true,
		Result:           vConfig.Interface(),
	})
	if err != nil {
		return err
	}

	return d.Decode(cfg)
}

type ExecutablePlugin struct {
	*Plugin
	Config reflect.Value
	ExecFn reflect.Value
}
