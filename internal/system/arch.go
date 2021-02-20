package system

import (
	"runtime"

	"github.com/pkg/errors"
)

const (
	OsWindows = "windows"
	OsDarwin  = "darwin"
	OsLinux   = "linux"
)

// Arch The Processor Architecture the CLI is running at
type Arch string

func (a *Arch) String() string {
	return string(*a)
}

func GetCurrent() Arch {
	switch runtime.GOOS {
	case OsDarwin:
		return Darwin
	case OsWindows:
		panic(errors.Errorf("'%s' not yet supported", OsWindows))
	default:
		return Linux
	}
}

const (

	// Darwin
	Darwin Arch = OsDarwin

	// Linux
	Linux Arch = OsLinux

	// MacOsx64 amd64
	MacOsx64 Arch = "darwinx64"

	// Linux64
	Linux64 Arch = "linuxx64"

	// LinuxArm32
	LinuxArm32 Arch = "linuxarm32"
)
