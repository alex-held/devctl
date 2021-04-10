package plugins

import (
	"fmt"
	"runtime"
	"strings"
)

type RuntimeInfo struct {
	OS, Arch string
}

func (info RuntimeInfo) Format(pattern string, args ...interface{}) string {
	formatted := fmt.Sprintf(pattern, args...)
	formatted = strings.ReplaceAll(formatted, "[os]", strings.ToLower(info.OS))
	formatted = strings.ReplaceAll(formatted, "[OS]", strings.ToUpper(info.OS))
	formatted = strings.ReplaceAll(formatted, "[arch]", strings.ToLower(info.Arch))
	formatted = strings.ReplaceAll(formatted, "[ARCH]", strings.ToUpper(info.Arch))
	return formatted
}

type RuntimeInfoGetter interface {
	Get() (info RuntimeInfo)
}

type OSRuntimeInfoGetter struct{}

// Format formats pattern string using the provided args and the runtime info
// [os] = OS (darwin, linux)
// [OS] = OS (DARWIN, LINUX)
// [arch] = amd64, arm64
// [ARCH] = AMD64, ARM64
func (g OSRuntimeInfoGetter) Format(pattern string, args ...interface{}) string {
	info := g.Get()
	return info.Format(pattern, args...)
}

func (OSRuntimeInfoGetter) Get() (info RuntimeInfo) {
	var osID = runtime.GOOS
	var archID = runtime.GOARCH
	if archID == "arm" {
		archID = "arm64"
	}
	return RuntimeInfo{
		OS:   osID,
		Arch: archID,
	}
}
