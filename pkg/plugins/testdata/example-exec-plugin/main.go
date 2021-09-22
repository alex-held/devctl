package example_exec_plugin

import (
	"errors"
	"fmt"

	"github.com/alex-held/devctl-kit/pkg/plugins"
)

type Config struct {
	*plugins.Context  `yaml:"-"`
	InstallPath string `yaml:"install_path"`
}

func CreateConfig() *Config {
	return &Config{
		Context: &plugins.Context{},
	}
}

func Exec(cfg *Config, args []string) (err error) {
	fmt.Fprintf(cfg.Out, "cfg.InstallPath=%s\n", cfg.InstallPath)
	for i, arg := range args {
		fmt.Fprintf(cfg.Out, "args[%d]=%s\n", i, arg)
	}
	return errors.New("exec done")
}
