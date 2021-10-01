package options

import (
	"os"
	"sync"

	"github.com/spf13/pflag"
)

type ConfigFlags struct {
	DevctlRoot *string

	lock sync.Mutex
}

func NewConfigFlags() *ConfigFlags {
	return &ConfigFlags{
		DevctlRoot: stringptr(os.ExpandEnv("$HOME/.devctl")),
		lock:       sync.Mutex{},
	}
}

func stringptr(s string) *string {
	return &s
}

// AddFlags binds client configuration flags to a given flagset
func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	if f.DevctlRoot != nil {
		flags.StringVar(f.DevctlRoot, "devctl-path", *f.DevctlRoot, "Path to the devctl root")
	}
}
