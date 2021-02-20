// +build darwin

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

func cacheHome() string {
	return filepath.Join(userHome(), "Library", "Caches")
}

func configHome() func(lazypath) string {
	return func(lp lazypath) string {
		return filepath.Join(userHome(), lp.getAppPrefix())
	}
}
