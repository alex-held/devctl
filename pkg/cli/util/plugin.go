package util

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/spf13/afero"

	devctlerrors "github.com/alex-held/devctl/pkg/errors"
)

type PluginHandler interface {
	Lookup(name string) (executablePath string, ok bool)
	Execute(executablePath string, cmdArgs, environment []string) (err error)
}

// DefaultPluginHandler implements PluginHandler
type DefaultPluginHandler struct {
	ValidPrefixes []string
	Fs            afero.Fs
}

func NewDefaultPluginHandler() PluginHandler {
	return &DefaultPluginHandler{
		ValidPrefixes: []string{"devctl"},
		Fs:            afero.NewOsFs(),
	}
}

func (d *DefaultPluginHandler) Lookup(name string) (executablePath string, ok bool) {
	for _, prefix := range d.ValidPrefixes {
		path, err := exec.LookPath(fmt.Sprintf("%s-%s", prefix, name))
		if err != nil || len(path) == 0 {
			continue
		}
		return path, true
	}
	return "", false
}

func (d *DefaultPluginHandler) Execute(executablePath string, cmdArgs, environment []string) (err error) {
	if runtime.GOOS == "windows" {
		return devctlerrors.ErrWindowsNotSupported
	}

	// invoke cmd binary relaying the environment and args given
	// append executablePath to cmdArgs, as execve will make first argument the "binary name".
	return syscall.Exec(executablePath, append([]string{executablePath}, cmdArgs...), environment)
}
