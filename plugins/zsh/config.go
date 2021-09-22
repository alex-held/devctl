package zsh

import (
	"bytes"
	"io"

	"github.com/alex-held/devctl-kit/pkg/plugins"
	"gopkg.in/yaml.v3"
)

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

type Config struct {
	*plugins.Context `yaml:"-"`
	Vars             map[string]string `yaml:"vars,omitempty"`
	Exports          map[string]string `yaml:"exports,omitempty"`
	Aliases          map[string]string `yaml:"aliases,omitempty"`
	Completions      CompletionsSpec   `yaml:"completions,omitempty"`
}

type CompletionsSpec struct {
	CLI     map[string]string `yaml:"cli,omitempty"`
	Command map[string]string `yaml:"command,omitempty"`
}
