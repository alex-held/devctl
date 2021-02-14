package cli

import (
	"strings"
)

// StaticOption configures an instance ofg staticConfig lazily
type StaticOption func(config *staticConfig) *staticConfig

// DefaultStaticConfigFileOption configures defaults for ConfigFileName and ConfigFileType
func DefaultStaticConfigFileOption() StaticOption {
	return StaticConfigFileOption("config", "yaml")
}

// DefaultStaticCliConfigOption Configures default CliName, CliDescription and CliEnvPrefix
func DefaultStaticCliConfigOption() StaticOption {
	return StaticCliConfigOption("devctl", "A lightweight dev-environment manager / bootstrapper")
}

// StaticConfigFileOption configures the ConfigFileName and ConfigFileType
func StaticConfigFileOption(configName, configType string) StaticOption {
	return func(c *staticConfig) *staticConfig {
		c.configFileName = configName
		c.configFileType = configType
		return c
	}
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
