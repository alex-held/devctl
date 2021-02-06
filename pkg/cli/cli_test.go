package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStaticCliConfigOption
func TestStaticCliConfigOption(t *testing.T) {
	tCases := map[string]struct {
		Input    StaticOption
		Expected staticConfig
	}{
		"DefaultStaticCliConfigOption": {Input: DefaultStaticCliConfigOption(), Expected: staticConfig{
			cliName:        "devctl",
			cliDescription: "A lightweight dev-environment manager / bootstrapper",
			configFileName: "",
			configFileType: "",
			envPrefix:      "DEVCTL",
		}},
		"StaticCliConfigOption": {Input: StaticCliConfigOption("a", "b"), Expected: staticConfig{
			cliName:        "a",
			cliDescription: "b",
			configFileName: "",
			configFileType: "",
			envPrefix:      "A",
		}},
		"DefaultStaticConfigFileOption": {Input: DefaultStaticConfigFileOption(), Expected: staticConfig{
			cliName:        "",
			cliDescription: "",
			configFileName: "config",
			configFileType: "yaml",
			envPrefix:      "",
		}},
		"StaticConfigFileOption": {Input: StaticConfigFileOption(".devctl-cfg", "toml"), Expected: staticConfig{
			cliName:        "",
			cliDescription: "",
			configFileName: ".devctl-cfg",
			configFileType: "toml",
			envPrefix:      "",
		}},
	}

	for scenario, tc := range tCases {
		t.Run(scenario, func(scenarioT *testing.T) {
			staticCfg := &staticConfig{}
			actual := *tc.Input(staticCfg)
			assert.Equal(scenarioT, actual, tc.Expected)
		})
	}
}
