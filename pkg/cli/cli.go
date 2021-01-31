package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/coreos/etcd/client"
)

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
