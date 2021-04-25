//go:generate pluggen gen -o "$PWD/testdata/plugins/sdk-01.so" -p "./plugins/sdk-01" --pkg devctl

package plugins

import (
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

var _ = Describe("Manager", func() {
	var fs afero.Fs
	var pather devctlpath.Pather
	var sut pluginManager
	var sdkPluginDir string

	BeforeEach(func() {
		fs = afero.NewOsFs()
		pather = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			tmp, _ := os.MkdirTemp("devctl", "PluginRegistryTests")
			return path.Join(tmp)
		}))
		sut = pluginManager{
			fs:       fs,
			pather:   pather,
			registry: GlobalRegistry.(registry),
		}
		sdkPluginDir = pather.ConfigRoot("plugins")
	})

	Context("ReloadPlugins", func() {

		When("NewRegistry", func() {
			It("Returns the registry", func() {
				Expect(sut).ShouldNot(BeNil())
			})
		})

		When("no plugin exists in the plugin search paths", func() {
			BeforeEach(func() {
				_ = fs.MkdirAll(sdkPluginDir, os.ModePerm)
			})

			It("doesn't return an error", func() {
				err := sut.LoadSDKPlugins()
				Expect(err).Should(Succeed())
			})
		})

		Describe("one sdk plugin exists in the /tmp/devctl/<prefix>/plugins", func() {

			BeforeEach(func() {
				data, _ := afero.ReadFile(afero.NewOsFs(), "testdata/plugins/sdk-01.so")
				_ = afero.WriteFile(fs, path.Join(sdkPluginDir, "sdk-01.so"), data, os.ModePerm)
			})

			AfterEach(func() {
				_ = fs.RemoveAll(path.Join(sdkPluginDir, "sdk-01.so"))
			})

			When("before the  pluginManager has started to load the plugins", func() {

				It("returns a provider with no plugins", func() {
					provider := sut.GetProvider()
					plugins := provider()
					Expect(plugins).Should(HaveLen(0))
				})

				It("returns not nil", func() {
					Expect(sut.GetProvider()).ShouldNot(BeNil())
				})
			})

			When("after the  pluginManager has loaded the plugins", func() {

				BeforeEach(func() {
					_ = sut.LoadSDKPlugins()
				})
				It("returns a provider with one plugin", func() {
					provider := sut.GetProvider()
					plugins := provider()
					Expect(plugins).Should(HaveLen(1))
				})

			})
		})
	})

})
