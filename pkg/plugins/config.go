package plugins

import (
	"bytes"
	"os"
	"path"
	template2 "text/template"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Version    string `yaml:"version"`
	PluginSpec `yaml:"plugin"`
}

type CommandSpec struct {
	Cmd         string        `yaml:"cmd"`
	Help        string        `yaml:"help,omitempty"`
	Subcommands []CommandSpec `yaml:"subcommands,omitempty"`
}

type ConfigSpec struct {
	Values  map[string]string
	Static  map[string]interface{}
	Dynamic map[string]string
}

func ResolveDynamic(c ConfigSpec, e *Engine) (cfg ConfigSpec, err error) {
	cfg = c

	templates := map[string]*template2.Template{}

	for k, v := range c.Dynamic {
		if tmpl, err := template2.New(k).Parse(v); err == nil {
			templates[k] = tmpl
			continue
		}
		return cfg, err
	}

	home, _ := os.UserHomeDir()
	values := map[string]interface{}{
		"DEVCTL_PATH_USERHOME": home,
		"DEVCTL_PATH_ROOT":     e.cfg.Pather.ConfigRoot(),
		"DEVCTL_PATH_CONFIG":   e.cfg.Pather.Config(),
		"DEVCTL_PATH_SDK":      e.cfg.Pather.SDK(),
		"DEVCTL_PATH_PLUGIN":   e.cfg.Pather.Plugin(),
		"values":               c.Values,
	}

	for key, tmpl := range templates {
		w := &bytes.Buffer{}
		if err := tmpl.Execute(w, values); err != nil {
			return cfg, err
		}
		cfg.Dynamic[key] = w.String()
	}

	return cfg, nil
}

type PluginSpec struct {
	*CommandSpec `yaml:"cmd,inline"`
	Name         string     `yaml:"name"`
	Pkg          string     `yaml:"pkg"`
	Config       ConfigSpec `yaml:"config,omitempty"`
}

func (e *Engine) LoadPlugin(manifestPath string) (p *Plugin, err error) {
	m, err := e.loadManifest(manifestPath)
	if err != nil {
		return nil, err
	}

	rootPath := path.Dir(manifestPath)
	b, err := afero.ReadFile(e.cfg.Fs, path.Join(rootPath, "main.go"))
	if err != nil {
		return nil, err
	}

	p = &Plugin{
		Manifest: m,
		Source:   string(b),
		RootPath: rootPath,
	}

	e.pluginCache[p.Cmd] = p
	return p, nil
}

func (e *Engine) loadManifest(manifestPath string) (m *Manifest, err error) {
	bytes, err := afero.ReadFile(e.cfg.Fs, manifestPath)
	if err != nil {
		return nil, err
	}
	m = &Manifest{}
	err = yaml.Unmarshal(bytes, m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

type Plugins []*Plugin

func (e *Engine) LoadPlugins() (plugins Plugins) {
	pluginsRoot := e.cfg.Pather.Plugin()
	fis, err := afero.ReadDir(e.cfg.Fs, pluginsRoot)
	if err != nil {
		return plugins
	}
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}

		manifestPath := path.Join(pluginsRoot, fi.Name(), "plugin.yaml")
		if _, err = e.cfg.Fs.Stat(manifestPath); err != nil {
			continue
		}
		p, err := e.LoadPlugin(manifestPath)
		if err != nil {
			continue
		}
		plugins = append(plugins, p)
	}

	return plugins
}
