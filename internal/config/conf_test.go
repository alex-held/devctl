package config

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"

	"github.com/alex-held/devctl/pkg/logging"
)

var expectedConfig = &DevEnvConfig{
	GlobalConfig: DevEnvGlobalConfig{Version: "v1"},
	Sdks: map[string]DevEnvSDKConfig{
		"java": {
			Current: "openjdk-11",
			Candidates: []SDKCandidate{
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
		"haskell": {},
	},
}

func TestSerializeConfig(t *testing.T) {
	b := &bytes.Buffer{}
	logger := logging.NewLogger(
		logging.WithBuffer(b),
		logging.WithLevel(logging.LogLevelDebug),
		logging.WithName("Serialize"),
	)

	yaml := expectedConfig.GoString()
	logger.Print(yaml)
}

type pather struct {
	DevEnvConfigPath string
}

func (p *pather) ConfigFilePath() string           { return p.DevEnvConfigPath }
func (p *pather) ConfigRoot(elem ...string) string { return "" }
func (p *pather) Config(elem ...string) string     { return "" }
func (p *pather) Bin(elem ...string) string        { return "" }
func (p *pather) Download(elem ...string) string   { return "" }
func (p *pather) SDK(elem ...string) string        { return "" }
func (p *pather) Cache(elem ...string) string      { return "" }

func TestViperConfig(t *testing.T) {
	g := goblin.Goblin(t)

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("ViperConfig", func() {
		const testdataPath = "testdata/config.yaml"
		var b *bytes.Buffer
		var logger logging.Log
		var cfg *DevEnvConfig

		g.BeforeEach(func() {
			b = &bytes.Buffer{}
			logger = logging.NewLogger(
				logging.WithLevel(logging.LogLevelDebug),
				logging.WithBuffer(b),
				logging.WithName("ViperConfigTest"),
				logging.WithOutputs(io.MultiWriter(os.Stderr, os.Stdout)),
			)
		})

		g.It("WHEN Loading Config from file => THEN config contains all configurations of the config-file", func() {
			var err error
			cfg, err = parseConfigFile(testdataPath)
			if err != nil {
				g.Failf(err.Error())
			}

			if cfg == nil {
				g.Failf("could not read config file. ")
			}
			logger.Warnf("%+v", cfg)

			//			Expect(cfg.GlobalConfig).To(Equal(DevEnvGlobalConfig{Version: "v1"}))
			//	Expect(cfg.Sdks["java"]).Should(Not(BeEmpty()))
			//		Expect(cfg.Sdks["haskell"]).Should(Not(BeEmpty()))
			Expect(*cfg).To(Equal(*expectedConfig))
		})
	})
}
