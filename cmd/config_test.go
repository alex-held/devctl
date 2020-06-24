package cmd

import (
	"fmt"
	"testing"

	"github.com/gosuri/uitable"
	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	config, err := readConfig("config", map[string]interface{}{
		"DEVENV_HOME": "$HOME/.devenv",
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
