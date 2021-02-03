package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedConfig = &DevEnvConfig{
	GlobalConfig: DevEnvGlobalConfig{Version: "v1"},
	SDKConfig: DevEnvSDKSConfig{SDKS: []DevEnvSDKConfig{
		{
			SDK:     "java",
			Current: "openjdk-11",
			Installations: []DevEnvSDKInstallationConfig{
				{
					Path:    "/Library/Java/VirtualMachines/OpenJDK15/Contents/Home",
					Version: "openjdk-15",
				},
				{
					Path:    "/Library/Java/VirtualMachines/OpenJDK11/Contents/Home",
					Version: "openjdk-11",
				},
			},
		},
		{
			SDK: "haskell",
		},
	}},
}

func TestViperConfig(t *testing.T) {
	InitViper("testdata/devenv.yaml")
	cfg := LoadViperConfig()

	msg := fmt.Sprintf("Loaded 'testdata/devenv.yaml': \n%+v\n", *cfg)
	fmt.Println(msg)
	require.Equal(t, expectedConfig, cfg)
}
