package devctlpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alex-held/devctl/internal/devctlpath/xdg"
)

const (
	ConfigHomeRootEnvVar = "DEVCTL_CONFIG_HOME"
	CacheHomeEnvVar      = "DEVCTL_CACHE_HOME"
)

type lazypath string
type lazypathFinder struct {
	lp     lazypath
	finder finder
}

func NewLazyFinder(appPrefix string, userHomeFn UserHomePathFinder, cacheFn CachePathFinder, configRootFn ConfigRootFinder) lazypathFinder {
	return lazypathFinder{
		lp: lazypath(appPrefix),
		finder: finder{
			GetUserHomeFn:   userHomeFn,
			GetCachePathFn:  cacheFn,
			GetConfigRootFn: configRootFn,
		},
	}
}

func (f lazypathFinder) resolveSubDir(sub string, elem ...string) string {
	subConfig := f.configRoot(sub)
	return filepath.Join(subConfig, filepath.Join(elem...))
}

func (f lazypathFinder) configRoot(elem ...string) string {

	// There is an order to checking for a path.
	// 1. GetConfigRootFn has been provided
	// 1. GetUserHomeFn + AppPrefix has been provided
	// 2. See if a devctl specific environment variable has been set.
	// 2. Check if an XDG environment variable is set
	// 3. Fall back to a default

	if f.finder.GetConfigRootFn != nil {
		p := f.finder.ConfigRoot()
		return filepath.Join(p, filepath.Join(elem...))
	}

	if f.finder.GetUserHomeFn != nil {
		p := f.finder.GetUserHomeFn()
		p = filepath.Join(p, f.lp.getAppPrefix())
		return filepath.Join(p, filepath.Join(elem...))
	}

	base := os.Getenv(ConfigHomeRootEnvVar)
	if base != "" {
		return filepath.Join(base, filepath.Join(elem...))
	}

	base = os.Getenv(xdg.ConfigHomeEnvVar)
	if base != "" {
		confRoot := filepath.Join(base, f.lp.getAppPrefix())
		return filepath.Join(confRoot, filepath.Join(elem...))
	}

	base = configHome()(f.lp)
	return filepath.Join(base, filepath.Join(elem...))
}

// cachePath resolves the path where devctl will cache data
// There is an order to checking for a path.
// 1. GetCachePathFn has been provided
// 2. See if a devctl specific environment variable has been set.
// 2. Check if an XDG environment variable is set
// 3. Fall back to a default
func (f lazypathFinder) cachePath(elem ...string) string {

	fqrdn := fmt.Sprintf("io.alexheld%s", f.lp.getAppPrefix())

	if f.finder.GetCachePathFn != nil {
		p := f.finder.CachePath()
		p = filepath.Join(p, fqrdn)
		return filepath.Join(p, filepath.Join(elem...))
	}

	p := os.Getenv(CacheHomeEnvVar)
	if p != "" {
		p = filepath.Join(p, fqrdn)
		return filepath.Join(p, filepath.Join(elem...))
	}

	p = os.Getenv(xdg.ConfigHomeEnvVar)
	if p != "" {
		p := filepath.Join(p, fqrdn)
		return filepath.Join(p, filepath.Join(elem...))
	}


	p = cacheHome()
	p = filepath.Join(p, fqrdn)
	return filepath.Join(p, filepath.Join(elem...))
}

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
	return l.path(ConfigHomeRootEnvVar, xdg.ConfigHomeEnvVar, func() string {
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
