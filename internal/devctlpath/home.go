// Package devctlpath calculates filesystem paths to devctl's configuration, cache and data.
package devctlpath

const devctlConfigFileName = "config.yaml"
var lf = NewLazyFinder(".devctl", nil,nil,nil)

type PathFinder interface {
	UserHome() string
	CachePath() string
	ConfigRoot() string
}

type finder struct {
	GetUserHomeFn   UserHomePathFinder
	GetCachePathFn  CachePathFinder
	GetConfigRootFn ConfigRootFinder
}

func (f *finder) UserHomePathFinder() string {return f.GetUserHomeFn()}
func (f *finder) CachePath() string          {return f.GetCachePathFn()}
func (f *finder) ConfigRoot() string         {return f.GetConfigRootFn()}

type UserHomePathFinder func() string
type CachePathFinder func() string
type ConfigRootFinder func() string



// Path DevCtlConfigRoot the path where Helm stores configuration.
func DevCtlConfigRoot(elem ...string) string { return lf.DevCtlConfigRoot(elem...) }
func (f *lazypathFinder) DevCtlConfigRoot(elem ...string) string { return f.configRoot(elem...) }


// Path DevCtlConfigFile  path where Helm stores configuration.
func DevCtlConfigFilePath() string { return lf.DevCtlConfigFilePath() }
func (f *lazypathFinder) DevCtlConfigFilePath(elem ...string) string { return f.configRoot(devctlConfigFileName) }

// ConfigPath returns the path where Helm stores configuration.
func ConfigPath(elem ...string) string                { return lf.Config(elem...) }
func (f *lazypathFinder) Config(elem ...string) string { return f.resolveSubDir("config", elem...) }

// BinPath returns the path where Helm stores configuration.
func BinPath(elem ...string) string { return lf.Bin(elem...) }
func (f *lazypathFinder) Bin(elem ...string) string { return f.resolveSubDir("bin", elem...) }

// ConfigPath returns the path where Helm stores configuration.
func DownloadPath(elem ...string) string { return lf.Download( elem...) }
func (f *lazypathFinder) Download(elem ...string) string { return f.resolveSubDir("downloads", elem...) }

// Path returns the path where Helm stores configuration.
func SDKsPath(elem ...string) string { return lf.SDK(elem...) }
func (f *lazypathFinder) SDK(elem ...string) string { return f.resolveSubDir("sdks", elem...) }

// CachePath returns the path where Helm stores cached objects.
func CachePath(elem ...string) string { return lf.Cache(elem...) }
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
