package app

import (
	"fmt"
	"os"

	"github.com/alex-held/devctl/internal/config"
	"github.com/alex-held/devctl/pkg/constants"
	"github.com/alex-held/devctl/pkg/devctlpath"
	"github.com/alex-held/devctl/pkg/logging"
)

var (
	cliInstance *app
)

const (
	appName = "devctl"
)

type PatherGetter interface {
	GetPather() devctlpath.Pather
}

type CLI interface {
	ConfigGetter
	ConfigUpdater
	PatherGetter
	Name() string
	Version() string
	ConfigFileName() string
}

type ConfigGetter interface {
	GetConfig() (cfg *config.DevCtlConfig, err error)
	MustGetConfig() (cfg *config.DevCtlConfig)
}

type ConfigUpdater interface {
	UpdateConfig(cfg *config.DevCtlConfig) (err error)
}

type app struct {
	Log         logging.Log // Handles all logging
	GlobalFlags GlobalFlags
	Args        []string
	Pather      devctlpath.Pather
}

func (a *app) GetPather() devctlpath.Pather {
	return a.Pather
}

// Version returns the semver of the devctl binary
func (a *app) Version() string {
	return a.GlobalFlags.AppVersion
}

// ConfigFileName returns the path of the DevCtlConfigFile
// defaults to $HOME/.devctl/config.yaml
func (a *app) ConfigFileName() string {
	return a.GlobalFlags.Config
}

// Name returns the name of the devctl binary
func (a *app) Name() string {
	return appName
}

func (a *app) MustGetConfig() (cfg *config.DevCtlConfig) {
	cfg, err := a.GetConfig()
	ExitWhenError(constants.NoConfigFileDetected, err)
	return cfg
}

func (a *app) GetConfig() (cfg *config.DevCtlConfig, err error) {
	cfg, err = config.ParseConfigFile(a.ConfigFileName())
	return cfg, err
}

func (a *app) UpdateConfig(cfg *config.DevCtlConfig) error {
	return config.WriteDevEnvConfig(a.ConfigFileName(), *cfg)
}

// GetCLI returns the global CLI
func GetCLI() CLI {
	return cliInstance
}

func GetOrCreateCLI(globalFlags GlobalFlags, args []string) CLI {
	if cliInstance == nil {
		cliInstance = newApp(globalFlags, args)
	}
	return cliInstance
}

// ExitWhenError  prints an error message and exits the application with ErrorCode: code
func ExitWhenError(code int, err error) {
	if err == nil {
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(code)
}

func newApp(globalFlags GlobalFlags, args []string) (cli *app) {
	logger := logging.NewLogger()

	if globalFlags.Verbose {
		logger = logging.NewLogger(logging.WithVerbose(true), logging.WithLevel(logging.LogLevelDebug))
	}

	cli = &app{
		Log:         logger,
		GlobalFlags: globalFlags,
		Args:        args,
		Pather:      devctlpath.DefaultPather(),
	}
	return cli
}
