package env

import (
	"os"
	"path/filepath"

	"github.com/alex-held/devctl-kit/pkg/log"
	"github.com/pkg/errors"

	"github.com/alex-held/devctl/pkg/constants"
)

type Paths struct {
	base string
	tmp  string
}

// MustGetPaths returns the inferred paths for devctl. By default, it assumes
// $HOME/.devctl as the base path, but can be overridden via DEVCTL_ROOT environment
// variable.
func MustGetPaths() Paths {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(errors.Wrap(err, "cannot get user home dir"))
	}
	base := filepath.Join(homeDir, constants.DefaultDevctlDir)
	if fromEnv := os.Getenv(constants.DEVCTL_ROOT_KEY); fromEnv != "" {
		base = fromEnv
		log.Infof("using environment override %s=%s", constants.DEVCTL_ROOT_KEY, fromEnv)
	}

	base, err = filepath.Abs(base)
	if err != nil {
		panic(errors.Wrap(err, "cannot get absolute path"))
	}
	return NewPaths(base)
}

func NewPaths(base string) Paths {
	return Paths{
		base: base,
		tmp:  os.TempDir(),
	}
}

// Base returns the devctl base directory
func (p Paths) Base() string { return p.base }

// IndexBase returns the devctl index directory
func (p Paths) IndexBase() string { return filepath.Join(p.base, "index") }

// // Index returns the devctl index directory
// func (p Paths) Index(paths ...string) string {
// 	return filepath.Join(p.base, filepath.Join(paths...), "index")
// }

// IndexPath returns the directory where a plugin index repository is cloned.
//
// e.g. {BasePath}/index/default or {BasePath}/index
func (p Paths) IndexPath(name string) string {
	return filepath.Join(p.base, "index", name)
}

// IndexPluginsPath returns the plugins directory of an index repository.
//
// e.g. {BasePath}/index/default/plugins/ or {BasePath}/index/plugins/
func (p Paths) IndexPluginsPath(name string) string {
	return filepath.Join(p.IndexPath(name), "plugins")
}

// IndexPluginManifestPath returns the plugins directory of an index repository.
//
// e.g. {BasePath}/index/default/plugins/ or {BasePath}/index/plugins/
func (p Paths) IndexPluginManifestPath(indexName, pluginName string) string {
	return filepath.Join(p.IndexPluginsPath(indexName), pluginName+constants.ManifestExtension)
}

// InstallReceiptsPath returns the base directory where plugin receipts are stored.
//
// e.g. {BasePath}/receipts
func (p Paths) InstallReceiptsPath() string { return filepath.Join(p.base, "receipts") }

// BinPath returns the path where plugin executable symbolic links are found.
// This path should be added to $PATH in client machine.
//
// e.g. {BasePath}/bin
func (p Paths) BinPath() string { return filepath.Join(p.base, "bin") }

// InstallPath returns the base directory for plugin installations.
//
// e.g. {BasePath}/store
func (p Paths) InstallPath() string { return filepath.Join(p.base, "store") }

// PluginInstallPath returns the path to install the plugin.
//
// e.g. {InstallPath}/{version}/{..files..}
func (p Paths) PluginInstallPath(plugin string) string {
	return filepath.Join(p.InstallPath(), plugin)
}

// PluginInstallReceiptPath returns the path to the install receipt for plugin.
//
// e.g. {InstallReceiptsPath}/{plugin}.yaml
func (p Paths) PluginInstallReceiptPath(plugin string) string {
	return filepath.Join(p.InstallReceiptsPath(), plugin+constants.ManifestExtension)
}

// PluginVersionInstallPath returns the path to the specified version of specified
// plugin.
//
// e.g. {PluginInstallPath}/{plugin}/{version}
func (p Paths) PluginVersionInstallPath(plugin, version string) string {
	return filepath.Join(p.InstallPath(), plugin, version)
}

// Realpath evaluates symbolic links. If the path is not a symbolic link, it
// returns the cleaned path. Symbolic links with relative paths return error.
func Realpath(path string) (string, error) {
	if !IsOsFs(fs) {
		panic("realpath is only supported for afero.OsFs")
	}

	s, err := os.Lstat(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to stat the currently executed path (%q)", path)
	}

	if s.Mode()&os.ModeSymlink != 0 {
		if path, err = os.Readlink(path); err != nil {
			return "", errors.Wrap(err, "failed to resolve the symlink of the currently executed version")
		}
		if !filepath.IsAbs(path) {
			return "", errors.Errorf("symbolic link is relative (%s)", path)
		}
	}
	return filepath.Clean(path), nil
}
