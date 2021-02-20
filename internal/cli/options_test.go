package cli

import (
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

// TestStaticCliConfigOption
func TestStaticCliConfigOption(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("StaticCliConfigOption", func() {
		staticCfg := &staticConfig{}

		g.JustBeforeEach(func() {
			staticCfg = &staticConfig{}
		})

		g.It("DefaultStaticCliConfigOption", func() {
			in := DefaultStaticCliConfigOption()
			expected := staticConfig{
				cliName:        "devctl",
				cliDescription: "A lightweight dev-environment manager / bootstrapper",
				configFileName: "",
				configFileType: "",
				envPrefix:      "DEVCTL",
			}
			actual := in(staticCfg)
			Expect(*actual).To(Equal(expected))
		})

		g.It("StaticCliConfigOption", func() {
			in := StaticCliConfigOption("a", "b")
			expected := staticConfig{
				cliName:        "a",
				cliDescription: "b",
				configFileName: "",
				configFileType: "",
				envPrefix:      "A",
			}
			actual := in(staticCfg)
			Expect(*actual).To(Equal(expected))
		})

		g.It("DefaultStaticConfigFileOption", func() {
			in := DefaultStaticConfigFileOption()
			expected := staticConfig{
				cliName:        "",
				cliDescription: "",
				configFileName: "config",
				configFileType: "yaml",
				envPrefix:      "",
			}
			actual := in(staticCfg)
			Expect(*actual).To(Equal(expected))
		})

		g.It("StaticConfigFileOption", func() {
			in := StaticConfigFileOption(".devctl-cfg", "toml")
			expected := staticConfig{
				cliName:        "",
				cliDescription: "",
				configFileName: ".devctl-cfg",
				configFileType: "toml",
				envPrefix:      "",
			}
			actual := in(staticCfg)
			Expect(*actual).To(Equal(expected))
		})
	})
}
