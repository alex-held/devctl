// +build darwin linux windows

package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/alex-held/devctl/internal/system"

	"github.com/pkg/errors"
)

// ConfigReader reads Config
type ConfigReader interface {
	io.Reader
}

// EnvSettings describes all of the environment settings.
type EnvSettings struct {
	namespace string
	// config    *genericclioptions.ConfigFlags

	// KubeConfig is the path to the kubeconfig file
	KubeConfig string
	// KubeContext is the name of the kubeconfig context.
	KubeContext string
	// Bearer KubeToken used for authentication
	KubeToken string
	// Username to impersonate for the operation
	KubeAsUser string
	// Groups to impersonate for the operation, multiple groups parsed from a comma delimited list
	KubeAsGroups []string
	// Kubernetes API Server Endpoint for authentication
	KubeAPIServer string
	// Custom certificate authority file.
	KubeCaFile string
	// Debug indicates whether or not Helm is running in Debug mode.
	Debug bool
	// RegistryConfig is the path to the registry config file.
	RegistryConfig string
	// RepositoryConfig is the path to the repositories file.
	RepositoryConfig string
	// RepositoryCache is the path to the repository cache directory.
	RepositoryCache string
	// PluginsDirectory is the path to the plugins directory.
	PluginsDirectory string
	// MaxHistory is the max release history maintained.
	MaxHistory int
}

type Env struct {
	sync.RWMutex

	// cmd        CLI
	//	config     ConfigReader
	HomeFinder HomeFinder
	// writer        ConfigWriter
	// Test          *TestParameters
	// updaterConfig UpdaterConfigReader
	//	file *JSONFile
}

type ConfigGetter func() string
type EnvGetter func(env string) string

type HomeFinder interface {
	BinDir() string
	ConfigDir() string
	DownloadsDir() string
	Home() string
	SDKDir(sdk ...string) string
	SDKRoot() string
	LogDir() string
}

// Base
// ===============
//

// Base shared cli.HomeFinder for multiple runtime.GOOS
type Base struct {
	appPrefix       string
	getUserHomeFunc ConfigGetter
	getenvFunc      EnvGetter
}

func (b Base) Join(elem ...string) string { return filepath.Join(elem...) }

func (b Base) getAppPrefix() (prefix string) {
	prefix = strings.ToLower(fmt.Sprintf(".%s", strings.TrimPrefix(b.appPrefix, ".")))
	return prefix
}

func (b Base) getUserHome() (home string) {
	if b.getUserHomeFunc != nil {
		home = b.getUserHomeFunc()
		if home != "" {
			return home
		}
	}

	home = b.getenv("HOME")
	if home != "" {
		return home
	}

	home, _ = os.UserHomeDir()
	return home
}

func (b Base) getenv(s string) string {
	if b.getenvFunc != nil {
		return b.getenvFunc(s)
	}
	return os.Getenv(s)
}

// IsDarwin
// ===============
//

type Darwin struct {
	Base
}

func (d *Darwin) UserHome() string {
	return d.getUserHome()
}
func (d *Darwin) Home() string {
	user := d.UserHome()
	return d.Join(user, d.getAppPrefix())
}
func (d *Darwin) BinDir() string {
	return d.Join(d.Home(), "bin")
}

func (d *Darwin) ConfigDir() string {
	return d.Join(d.Home(), "config")
}

func (d *Darwin) DownloadsDir() string {
	return d.Join(d.Home(), "downloads")
}

func (d *Darwin) SDKDir(sdk ...string) string {
	return d.Join(append([]string{d.SDKRoot()}, sdk...)...)
}

func (d *Darwin) SDKRoot() string {
	return d.Join(d.Home(), "sdks")
}

func (d *Darwin) LogDir() string {
	return d.Join(d.Home(), "logs")
}

// UNIX
// ===============
//

type XdgPosix struct {
	Base
}

func (x *XdgPosix) UserHome() string {
	return x.getUserHome()
}

func (x *XdgPosix) Home() string {
	user := x.UserHome()
	home := x.Join(user, ".config", x.getAppPrefix())
	return home
}

func (x *XdgPosix) BinDir() string {
	return x.Join(x.Home(), "bin")
}

func (x *XdgPosix) ConfigDir() string {
	return x.Join(x.Home(), "config")
}

func (x *XdgPosix) DownloadsDir() string {
	return x.Join(x.Home(), "downloads")
}

func (x *XdgPosix) SDKDir(sdk ...string) string {
	return x.Join(append([]string{x.SDKRoot()}, sdk...)...)
}

func (x *XdgPosix) SDKRoot() string {
	return x.Join(x.Home(), "sdks")
}

func (x *XdgPosix) LogDir() string {
	return x.Join(x.Home(), "logs")
}

func NewHomeFinderForOS(goos, appPrefix string, getHome ConfigGetter, getEnv EnvGetter) HomeFinder {
	base := Base{
		appPrefix:       appPrefix,
		getUserHomeFunc: getHome,
		getenvFunc:      getEnv,
	}

	switch goos {
	case system.OsDarwin:
		return &Darwin{Base: base}
	case system.OsWindows:
		panic(errors.Errorf("runtime 'windows' is not yet supported"))
	default:
		return &XdgPosix{Base: base}
	}
}

func DefaultHomeFinder(appPrefix string) HomeFinder {
	return NewHomeFinder(appPrefix, func() string {
		home, _ := os.UserHomeDir()
		return home
	}, os.Getenv)
}

func NewHomeFinder(appPrefix string, getHome ConfigGetter, getEnv EnvGetter) HomeFinder {
	goos := runtime.GOOS
	return NewHomeFinderForOS(goos, appPrefix, getHome, getEnv)
}