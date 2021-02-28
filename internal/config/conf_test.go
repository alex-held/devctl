package config

import (
	"path"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/alex-held/devctl/internal/testutils"
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
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("ViperConfig", func() {
		const testdataPath = "testdata/devenv.yaml"
		dir := path.Dir(testdataPath)
		config := path.Base(testdataPath)

		g.It("WHEN Loading Config from file => THEN config contains all configurations of the config-file", func() {
			viper.AddConfigPath(dir)
			viper.SetConfigName(config)
			viper.SetConfigType("yaml")

			cfg := LoadViperConfig()
			logger := testutils.NewLogger(nil)
			fields := logrus.Fields{
				"testdata-path": testdataPath,
				"global-config": cfg.GlobalConfig,
				"SDK-config":    cfg.SDKConfig,
				"config":        *cfg,
			}

			logger.WithFields(fields).Traceln("loaded viper config from testdata/devenv.yaml")
			Expect(cfg.GlobalConfig).To(Equal(DevEnvGlobalConfig{Version: "v1"}))
			Expect(cfg.SDKConfig.SDKS[0].SDK).To(Equal("java"))
			Expect(cfg.SDKConfig.SDKS).To(HaveLen(2))
			Expect(cfg.SDKConfig.SDKS[0].Installations).To(HaveLen(2))
			Expect(cfg).To(Equal(expectedConfig))
		})
	})
}
