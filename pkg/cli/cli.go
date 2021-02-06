// Package cli
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coreos/etcd/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Once the Application starts, following values get configured exactly once.
// Those values cannot be changed while the application is running
// cliDescription string
// configFileName string
// configFileType string
// envPrefix       string
type staticConfig struct {
	cliName        string
	cliDescription string
	configFileName string
	configFileType string
	envPrefix      string
}

type CLI interface {
	Name() string
	ConfigFileName() string
	ConfigDir() string
}

var staticConfigAfterLoad *staticConfig

func (c *staticConfig) Name() string {
	return c.cliName
}

func (c *staticConfig) ConfigFileName() string {
	filename := filepath.Join(c.ConfigDir(), fmt.Sprintf("%s.%s", c.configFileName, c.configFileType))
	return filename
}

func (c *staticConfig) ConfigDir() string {
	home, err := homedir.Dir()
	if err != nil {
		ExitWithError(1, err)
	}
	dir := filepath.Join(home, fmt.Sprintf(".%s", c.Name()))
	return dir
}

// GetCLI e
func GetCLI() CLI {
	if staticConfigAfterLoad == nil {
		ConfigureStorage(DefaultStaticCliConfigOption(), DefaultStaticConfigFileOption())
	}
	return staticConfigAfterLoad
}

// ExitWithError  prints an error message and exits the application with ErrorCode: code
func ExitWithError(code int, err error) {
	_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
	if cerr, ok := err.(*client.ClusterError); ok {
		_, _ = fmt.Fprintln(os.Stderr, cerr.Detail())
	}
	os.Exit(code)
}

// StaticOption configures an instance ofg staticConfig lazily
type StaticOption func(config *staticConfig) *staticConfig

// DefaultStaticConfigFileOption configures defaults for ConfigFileName and ConfigFileType
func DefaultStaticConfigFileOption() StaticOption {
	return StaticConfigFileOption("config", "yaml")
}

// StaticConfigFileOption configures the ConfigFileName and ConfigFileType
func StaticConfigFileOption(configName, configType string) StaticOption {
	return func(c *staticConfig) *staticConfig {
		c.configFileName = configName
		c.configFileType = configType
		return c
	}
}

// DefaultStaticCliConfigOption Configures default CliName, CliDescription and CliEnvPrefix
func DefaultStaticCliConfigOption() StaticOption {
	return StaticCliConfigOption("devctl", "A lightweight dev-environment manager / bootstrapper")
}

// StaticCliConfigOption Configures the CliName, CliDescription and CliEnvPrefix of this CLI application
func StaticCliConfigOption(cliName, cliDescription string) StaticOption {
	return func(c *staticConfig) *staticConfig {
		c.cliName = cliName
		c.cliDescription = cliDescription
		c.envPrefix = strings.ToUpper(c.cliName)
		return c
	}
}

// ConfigureStorage configures the config storage using multiple StaticOption's
func ConfigureStorage(option ...StaticOption) {
	c := newStaticConfig(option...)
	staticConfigAfterLoad = c

	viper.SetEnvPrefix(c.envPrefix)

	viper.AddConfigPath(c.ConfigDir())
	viper.SetConfigName(c.configFileName)
	viper.SetConfigType(c.configFileType)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		ExitWithError(1, err)
	}
}

func newStaticConfig(option ...StaticOption) (c *staticConfig) {
	c = &staticConfig{}
	for _, o := range option {
		c = o(c)
	}
	return c
}
