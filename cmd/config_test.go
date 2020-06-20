package cmd

import (
	"fmt"
	"path"
	"testing"

	"github.com/gosuri/uitable"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/alex-held/dev-env/meta"
)

func TestReadConfig(t *testing.T) {
	config, err := readConfig("config", map[string]interface{}{
		meta.DEVENV_HOME: "$HOME/.devenv",
	})
	require.NoError(t, err)

	config.SetEnvPrefix("DEVENV")
	config.Set("DEVENV_CONFIG", config.ConfigFileUsed())
	// err = config.WriteConfig()

	settings := config.AllSettings()
	PrintTable("all", settings)

	environment := config.Sub("environment").Sub("vars")
	PrintTable("environment.vars", environment.AllSettings())
}

type globalConfig struct {
	ErrorFile  string
	DevEnvHome string
}

var GlobalConfig = globalConfig{}

func TestReadPrefixedEnvironmentVariables(t *testing.T) {
	config := viper.New()
	config.SetConfigType("json")
	config.SetConfigFile("devenv")
	config.SetEnvPrefix("DEVENV")
	home, err := homedir.Dir()
	require.NoError(t, err)
	config.AddConfigPath(path.Join(home, ".devenv"))

	if err = config.ReadInConfig(); err != nil {
		_ = fmt.Errorf("Failed to read the configuration file: %s ", err)
		require.NoError(t, err)
	}

	if err = config.Unmarshal(&GlobalConfig); err != nil {
		_ = fmt.Errorf("Failed to unmarshal the configuration file: %s ", err)
		require.NoError(t, err)
	}

}

func PrintTable(title string, settings map[string]interface{}) {
	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("KEY", "VALUE")
	for key, value := range settings {
		table.AddRow(key, value)
	}
	fmt.Println(title)
	fmt.Println()
	fmt.Println(table)
}
