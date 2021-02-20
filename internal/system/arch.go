package system

import (
	"runtime"
	"strings"
)

const (
	OsWindows = "windows"
	OsDarwin  = "darwin"
	OsLinux   = "linux"
)

// Arch The Processor Architecture the CLI is running at
type Arch string

// String formats Arch to string
func (a *Arch) String() string {
	return strings.ToLower(string(*a))
}

func (a *Arch) IsOfFamily(family Arch) bool {
	switch {
	case family.IsWindows():
		return a.IsWindows()
	case family.IsLinux():
		return a.IsLinux()
	case family.IsDarwin():
		return a.IsDarwin()
	default:
		return strings.Contains(a.String(), family.String())
	}
}

func (a *Arch) IsWindows() bool { return strings.Contains(a.String(), OsWindows) }
func (a *Arch) IsDarwin() bool  { return strings.Contains(a.String(), OsDarwin) }
func (a *Arch) IsLinux() bool   { return strings.Contains(a.String(), OsLinux) }

var getDefaultGoosRuntimeGoos = func() string {
	return runtime.GOOS
}

var GetGoosFuncOverwrite func() string

func GetCurrent() Arch {
	var goosFunc func() string

	if GetGoosFuncOverwrite != nil {
		goosFunc = GetGoosFuncOverwrite
	} else {
		goosFunc = getDefaultGoosRuntimeGoos
	}

	switch goosFunc() {
	case OsDarwin:
		return Darwin
	case OsWindows:
		return Windows
	default:
		return Linux
	}
}

const (

	// Darwin
	Darwin Arch = OsDarwin

	// Linux
	Linux Arch = OsLinux

	// Windows
	Windows Arch = OsWindows

	// DarwinX64 amd64
	DarwinX64 Arch = "darwinx64"
	// LinuxX64
	LinuxX64 Arch = "linuxx64"
	// LinuxArm32
	LinuxArm32 Arch = "linuxarm32"
)
