// +build !windows,!darwin

package devctlpath

import (
	"os"
	"path/filepath"
)

// userHomeDir defines the user's home directory
func userHome() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = os.Getenv("HOME")
	}
	return home
}

// configHome defines the base directory relative to which user specific configuration files should
// be stored.
//
// If $XDG_CONFIG_HOME is either not set or empty, a default equal to $HOME/.config is used.
func configHome() func(lazypath) string {
	return func(lp lazypath) string {
		return filepath.Join(userHome(), ".config", lp.getAppPrefix())
	}
}

// cacheHome defines the base directory relative to which user specific non-essential data files
// should be stored.
//
// If $XDG_CACHE_HOME is either not set or empty, a default equal to $HOME/.cache is used.
func cacheHome() string {
	if cachehome := os.Getenv("XDG_CACHE_HOME"); cachehome != "" {
		return cachehome
	}
	return filepath.Join(userHome(), ".cache")
}
