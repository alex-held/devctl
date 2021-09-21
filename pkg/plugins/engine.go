package plugins

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/afero"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"

	"github.com/alex-held/devctl/pkg/devctlpath"
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

	// get symbols
	val, err := i.Eval(p.Pkg + ".New")
	newFunc := val.Interface().(func([]string) error)

	// execute plugin
	return newFunc(args)
}
