// Package devctlpath calculates filesystem paths to devctl's configuration, cache and data.
package devctlpath

var lp = lazypath(".devctl")

const devctlConfigFileName = "config.yaml"

// Path DevCtlConfigRoot the path where Helm stores configuration.
func DevCtlConfigRoot(elem ...string) string { return lp.configRoot(elem...) }

// Path DevCtlConfigFile  path where Helm stores configuration.
func DevCtlConfigFilePath() string { return lp.configRoot(devctlConfigFileName) }

// ConfigPath returns the path where Helm stores configuration.
func ConfigPath(elem ...string) string { return lp.configSubPath("config", elem...) }

// BinPath returns the path where Helm stores configuration.
func BinPath(elem ...string) string { return lp.configSubPath("bin", elem...) }

// ConfigPath returns the path where Helm stores configuration.
func DownloadPath(elem ...string) string { return lp.configSubPath("downloads", elem...) }

// Path returns the path where Helm stores configuration.
func SDKsPath(elem ...string) string { return lp.configSubPath("sdks", elem...) }

// CachePath returns the path where Helm stores cached objects.
func CachePath(elem ...string) string { return lp.cachePath(elem...) }

// CacheIndexFile returns the path to an index for the given named repository.
func CacheIndexFile(name string) string {
	if name != "" {
		name += "-"
	}
	return name + "index.yaml"
}

// CacheChartsFile returns the path to a text file listing all the charts
// within the given named repository.
func CacheChartsFile(name string) string {
	if name != "" {
		name += "-"
	}
	return name + "charts.txt"
}
