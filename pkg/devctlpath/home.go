// Package devctlpath calculates filesystem paths to devctl's configuration, cache and data.
package devctlpath

const devctlConfigFileName = "config.yaml"

var lf = NewPather()

type finder struct {
	GetUserHomeFn   UserHomePathFinder
	GetCachePathFn  CachePathFinder
	GetConfigRootFn ConfigRootFinder
}

func (f *finder) UserHomePathFinder() string { return f.GetUserHomeFn() }
func (f *finder) CachePath() string          { return f.GetCachePathFn() }
func (f *finder) ConfigRoot() string         { return f.GetConfigRootFn() }

type UserHomePathFinder func() string
type CachePathFinder func() string
type ConfigRootFinder func() string

// DevCtlConfigRoot the path where Helm stores configuration.
func DevCtlConfigRoot(elem ...string) string               { return lf.ConfigRoot(elem...) }
func (f *lazypathFinder) ConfigRoot(elem ...string) string { return f.configRoot(elem...) }

// DevCtlConfigFilePath  path where Helm stores configuration.
func DevCtlConfigFilePath() string { return lf.ConfigFilePath() }
func (f *lazypathFinder) ConfigFilePath() string {
	return f.configRoot(f.cfgName)
}

// ConfigPath returns the path where Helm stores configuration.
func DefaultPather() Pather                     { return lf }
func (f *lazypathFinder) DefaultPather() Pather { return f }

// ConfigPath returns the path where Helm stores configuration.
func ConfigPath(elem ...string) string                 { return lf.Config(elem...) }
func (f *lazypathFinder) Config(elem ...string) string { return f.resolveSubDir("config", elem...) }

// BinPath returns the path where Helm stores configuration.
func BinPath(elem ...string) string                 { return lf.Bin(elem...) }
func (f *lazypathFinder) Bin(elem ...string) string { return f.resolveSubDir("bin", elem...) }

// DownloadPath returns the path where Helm stores configuration.
func DownloadPath(elem ...string) string { return lf.Download(elem...) }
func (f *lazypathFinder) Download(elem ...string) string {
	return f.resolveSubDir("downloads", elem...)
}

// SDKsPath returns the path where Helm stores configuration.
func SDKsPath(elem ...string) string                { return lf.SDK(elem...) }
func (f *lazypathFinder) SDK(elem ...string) string { return f.resolveSubDir("sdks", elem...) }

// PluginPath returns the path where devctl plugins are located.
func PluginPath(elem ...string) string                 { return lf.SDK(elem...) }
func (f *lazypathFinder) Plugin(elem ...string) string { return f.resolveSubDir("plugins", elem...) }

// CachePath returns the path where Helm stores cached objects.
func CachePath(elem ...string) string                 { return lf.Cache(elem...) }
func (f *lazypathFinder) Cache(elem ...string) string { return f.cachePath(elem...) }

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
