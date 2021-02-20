package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/alex-held/devctl/internal/logging"

	"github.com/coreos/etcd/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cliInstance *app
)

type CLI interface {
	GetHomeFinder() HomeFinder
	Name() string
	ConfigFileName() string
	ConfigDir() string
}

type app struct {
	staticConfig *staticConfig
	context      Contextified
}

func (a *app) GetHomeFinder() HomeFinder {
	return a.context.g.Env.HomeFinder
}

func (a *app) Name() string {
	return a.staticConfig.cliName
}

func (a *app) ConfigFileName() (filename string) {
	filename = filepath.Join(
		a.ConfigDir(),
		fmt.Sprintf("%s.%s",
			a.staticConfig.configFileName,
			a.staticConfig.configFileType))
	return filename
}

func (a *app) ConfigDir() string {
	home, err := homedir.Dir()
	if err != nil {
		ExitWithError(1, err)
	}
	dir := filepath.Join(home, fmt.Sprintf(".%s", a.Name()))
	return dir
}

// GetOrCreateCLI e
func GetOrCreateCLI() CLI {
	if cliInstance == nil {
		cliInstance = newApp(DefaultStaticCliConfigOption(), DefaultStaticConfigFileOption())
		cliInstance.configureViper()
	}
	return cliInstance
}

// ExitWithError  prints an error message and exits the application with ErrorCode: code
func ExitWithError(code int, err error) {
	if err == nil {
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
	if cerr, ok := err.(*client.ClusterError); ok {
		_, _ = fmt.Fprintln(os.Stderr, cerr.Detail())
	}
	os.Exit(code)
}

func newApp(option ...StaticOption) (cli *app) {
	c := &staticConfig{}
	for _, o := range []StaticOption{DefaultStaticCliConfigOption(), DefaultStaticConfigFileOption()} {
		c = o(c)
	}

	for _, o := range option {
		c = o(c)
	}

	l := logging.NewLogger(func(l *logrus.Logger) *logrus.Logger {
		l.SetOutput(os.Stdout)
		return l
	})

	cli = &app{
		staticConfig: c,
		context: NewContextified(&GlobalContext{
			Log: l,
			VDL: NewVDebugLog(l),
			Env: &Env{
				RWMutex:    sync.RWMutex{},
				HomeFinder: DefaultHomeFinder(c.cliName),
			},
		}),
	}

	return cli
}

// ConfigureStorage configures the config storage using multiple StaticOption's
func (a *app) configureViper() {
	viper.SetEnvPrefix(a.staticConfig.envPrefix)
	viper.AddConfigPath(a.ConfigDir())
	viper.SetConfigName(a.staticConfig.configFileName)
	viper.SetConfigType(a.staticConfig.configFileType)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		ExitWithError(1, err)
	}
}