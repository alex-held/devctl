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
	osId := runtime.GOOS
	archId := runtime.GOARCH
	if archId == "arm" {
		archId = "arm64"
	}
	return RuntimeInfo{
		OS:   osId,
		Arch: archId,
	}
}
