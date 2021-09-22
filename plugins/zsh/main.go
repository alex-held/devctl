package zsh

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alex-held/devctl-kit/pkg/plugins"
)

type Config struct {
	OMZ_PLUGINS []string

	Completions map[string]string

	*plugins.Context
}

func CreateConfig() *Config {
	return &Config{
		OMZ_PLUGINS: []string{},
		Completions: map[string]string{},
	}
}

var ErrWrongArgumentsProvided = errors.New("number of arguments is invalid")

func Exec(cfg *Config, args []string) (err error) {
	if len(args) == 0 {
		return ErrWrongArgumentsProvided
	}

	args = args[0:]

	switch args[0] {
	case "init":
		sb := strings.Builder{}
		fmt.Printf(sb.String())
	}

	return nil
}
