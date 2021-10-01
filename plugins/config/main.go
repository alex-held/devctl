package config

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
	"github.com/alex-held/devctl-kit/pkg/plugins"
)

type Config struct {
	*plugins.Context `yaml:"-"`
	Vars             map[string]string `yaml:"vars,omitempty"`
	Plugins          []PluginSpec      `yaml:"plugins,omitempty"`
}

type PluginSpec struct {
	Name     string `yaml:"name"`
	Enabled  bool   `yaml:"enabled"`
	Config   string `yaml:"config,omitempty`
	Manifest string `yaml:"manifest`
}

func CreateConfig() *Config {
	pather := devctlpath.DefaultPather()

	c := &Config{
		Context: &plugins.Context{
			Pather:  pather,
			Out:     os.Stdout,
			Context: context.Background(),
		},
		Vars: map[string]string{
			"DEVCTL_PATH_ROOT":   pather.ConfigRoot(),
			"DEVCTL_PATH_CONFIG": pather.ConfigRoot(),
			"DEVCTL_PATH_PLUGIN": pather.Plugin(),
			"DEVCTL_PATH_SDK":    pather.SDK(),
		},
		Plugins: []PluginSpec{
			{
				Name:     "golang",
				Enabled:  true,
				Config:   "{{ .DEVCTL_PATH_CONFIG }}/golang/config.yaml",
				Manifest: "{{ .DEVCTL_PATH_PLUGIN }}/golang/plugin.yaml",
			},
			{
				Name:     "zsh",
				Enabled:  true,
				Config:   "{{ .DEVCTL_PATH_CONFIG }}/zsh/config.yaml",
				Manifest: "{{ .DEVCTL_PATH_PLUGIN }}/zsh/plugin.yaml",
			},
		},
	}

	return c
}

//go:embed "testdata/config.yaml"
var testConfigYaml string

type cfgYaml struct {
	Vars map[string]string      `yaml:"vars"`
	Spec map[string]interface{} `yaml:"spec"`
}

func LoadCfg() (*cfgYaml, error) {
	cfg := &cfgYaml{}
	if e := yaml.Unmarshal([]byte(testConfigYaml), cfg); e != nil {
		fmt.Printf("err=%v\n", e)
		return nil, e
	}
	data := cfg.Vars

	t, e := template.New("").Parse(testConfigYaml)
	if e != nil {
		fmt.Printf("err=%v\n", e)
		return cfg, e
	}

	out := &bytes.Buffer{}
	if e = t.Execute(out, data); e != nil {
		fmt.Printf("err=%v\n", e)
		return cfg, e
	}

	newYaml := out.String()
	if e = yaml.Unmarshal([]byte(newYaml), cfg); e != nil {
		fmt.Printf("err=%v\n", e)
		return cfg, e
	}

	return cfg, nil
}

//
// func (c *Config) Resolve() (cfg *Config, err error) {
// 	data := c.Vars
//
// 	tmpl := template.New("")
//
// 	for _, plug := range c.Plugins {
// 		tmpl, _ = tmpl.New("Plugins." + plug.Name + "Enabled").Parse(fmt.Sprint(plug.Enabled))
// 		tmpl, _ = tmpl.New("Plugins." + plug.Name + "Manifest").Parse(fmt.Sprint(plug.Manifest))
// 		tmpl, _ = tmpl.New("Plugins." + plug.Name + "Config").Parse(fmt.Sprint(plug.Config))
// 		tmpl, _ = tmpl.New("Plugins." + plug.Name + "Name").Parse(fmt.Sprint(plug.Name))
// 	}
//
// 	return
// }

func Exec(cfg *Config, args []string) {
	// fmt.Printf("config.Exec cfg=%v; args=%v\n", *cfg, args)

	switch args[1] {
	case "view":
		handleView(cfg, args[1:])
	}
}

func handleView(cfg *Config, args []string) (err error) {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	yaml := string(b)
	fmt.Println(yaml)
	return nil
}
