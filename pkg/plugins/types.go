package plugins

import (
	"runtime"
)

type SDKPlugin interface {
	Name() string
	ListVersions() []string
	InstallE(version string) (err error)
	Download(version string) (err error)
	Install(version string)
	NewFunc() SDKPlugin
}

type RuntimeInfo struct {
	OS, Arch string
}

type RuntimeInfoGetter interface {
	Get() (info RuntimeInfo)
}

type OSRuntimeInfoGetter struct{}

func (OSRuntimeInfoGetter) Get() (info RuntimeInfo) {
	osID := runtime.GOOS
	archID := runtime.GOARCH
	if archID == "arm" {
		archID = "arm64"
	}
	return RuntimeInfo{
		OS:   osID,
		Arch: archID,
	}
}
