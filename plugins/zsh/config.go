package zsh

import (
	"bytes"
	"io"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/alex-held/devctl-kit/pkg/generation/banner"

	"github.com/alex-held/devctl-kit/pkg/plugins"
)

type CompletionsSpec struct {
	CLI     map[string]string `yaml:"cli,omitempty"`
	Command map[string]string `yaml:"command,omitempty"`
}

type Config struct {
	*plugins.Context `yaml:"-"`
	Vars             map[string]string `yaml:"vars,omitempty"`
	Exports          map[string]string `yaml:"exports,omitempty"`
	Aliases          map[string]string `yaml:"aliases,omitempty"`
	Completions      CompletionsSpec   `yaml:"completions,omitempty"`
}

func ReadConfigFile(r io.Reader) (f *Config, err error) {
	b := &bytes.Buffer{}
	if _, err = io.Copy(b, r); err != nil {
		return nil, err
	}

	f = CreateConfig()

	if err = yaml.Unmarshal(b.Bytes(), f); err != nil {
		return nil, err
	}
	return f, nil
}

func (c *Config) DeepEqual(other *Config) bool {
	cC := Config{
		Context:     nil,
		Vars:        c.Vars,
		Exports:     c.Exports,
		Aliases:     c.Aliases,
		Completions: c.Completions,
	}

	otherC := Config{
		Context:     nil,
		Vars:        other.Vars,
		Exports:     other.Exports,
		Aliases:     other.Aliases,
		Completions: other.Completions,
	}

	return reflect.DeepEqual(cC, otherC)
}

func (c *Config) TemplateConfigs() (cfgs map[string]interface{}) {
	return map[string]interface{}{
		"completions": CompletionsTmplData{
			Header:          banner.GenerateBanner("Completions", banner.KIND_SHELL),
			COMPLETIONS_DIR: c.Pather.Config("zsh", "completions"),
			CONFIGFILE:      c.Pather.Config("zsh", "config.yaml"),
			Completions:     c.Completions,
		},
		"exports": c.Exports,
		"aliases": c.Aliases,
	}
}
