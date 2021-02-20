package devctlpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alex-held/devctl/internal/devctlpath/xdg"
)

type lazypath string

const (
	ConfigHomeEnvVar = "DEVCTL_CONFIG_HOME"
	CacheHomeEnvVar  = "DEVCTL_CACHE_HOME"
)

func (l lazypath) path(devctlEnvVar, xdgEnvVar string, defaultFn func() string, elem ...string) string {
	// There is an order to checking for a path.
	// 1. See if a Helm specific environment variable has been set.
	// 2. Check if an XDG environment variable is set
	// 3. Fall back to a default
	base := os.Getenv(devctlEnvVar)
	if base != "" {
		return filepath.Join(base, filepath.Join(elem...))
	}

	base = os.Getenv(xdgEnvVar)
	if base != "" {
		base = filepath.Join(base, l.getAppPrefix())
		return filepath.Join(base, filepath.Join(elem...))
	}
	if base == "" {
		base = defaultFn()
	}
	return filepath.Join(base, filepath.Join(elem...))
}

// cachePath defines the base directory relative to which user specific non-essential data files
// should be stored.
func (l lazypath) cachePath(elem ...string) string {
	return l.path(CacheHomeEnvVar, xdg.CacheHomeEnvVar, cacheHome, filepath.Join(elem...))
}

// configRoot defines the base directory relative to which user specific configuration files should
// be stored.
func (l lazypath) configRoot(elem ...string) string {
	return l.path(ConfigHomeEnvVar, xdg.ConfigHomeEnvVar, func() string {
		configHome := configHome()(l)
		return configHome
	}, elem...)
}

// configRoot defines the base directory relative to which user specific configuration files should
// be stored.
func (l lazypath) configSubPath(sub string, elem ...string) string {
	return filepath.Join(l.configRoot(sub), filepath.Join(elem...))
}

// userHomePath defines the base directory relative to which user home directory
func (l lazypath) userHomePath(elem ...string) string {
	return l.path("HOME", "HOME", userHome, filepath.Join(elem...))
}

func (l lazypath) getAppPrefix() (prefix string) {
	prefix = strings.ToLower(fmt.Sprintf(".%s", strings.TrimPrefix(string(l), ".")))
	return prefix
}
