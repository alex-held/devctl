//go:generate pluggen gen -o "$PWD/testdata/plugins/sdk-01.so" -p "./plugins/sdk-01" --pkg devctl

package plugins

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/plugins/sdk"
	"github.com/alex-held/devctl/pkg/devctlpath"
)

var _ = Describe("Registry", func() {
	var fs afero.Fs
	var pather devctlpath.Pather
	var sut Registry
	var sdkPluginDir string

	BeforeEach(func() {
		fs = afero.NewOsFs()
		pather = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			tmp, _ := os.MkdirTemp("devctl", "PluginRegistryTests")
			return path.Join(tmp)
		}))
		sut = NewRegistry(pather, fs)
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
				_, err := sut.ReloadPlugins()
				Expect(err).Should(Succeed())
			})
			It("returns an empty list of plugins", func() {
				plugins, _ := sut.ReloadPlugins()
				Expect(plugins).Should(BeEmpty())
			})
		})

		When("one sdk plugin exists in the /tmp/devctl/<prefix>/plugins", func() {

			BeforeEach(func() {
				data, _ := afero.ReadFile(afero.NewOsFs(), "testdata/plugins/sdk-01.so")
				_ = afero.WriteFile(fs, path.Join(sdkPluginDir, "sdk-01.so"), data, os.ModePerm)
			})

			AfterEach(func() {
				_ = fs.RemoveAll(path.Join(sdkPluginDir, "sdk-01.so"))
			})

			It("doesn't return an error", func() {
				_, err := sut.ReloadPlugins()
				Expect(err).Should(Succeed())
			})

			It("returns a list with one SDKPlugin", func() {
				plugins, _ := sut.ReloadPlugins()
				Expect(plugins).Should(HaveLen(1))
				sdkPlugin, ok := plugins[0].(sdk.SDKPlugin)
				Expect(ok).Should(BeTrue())
				Expect(sdkPlugin).ShouldNot(BeNil())
			})

			When("calling PluginName() on first SDKPlugin", func() {
				It("returns the current version", func() {
					plugins, _ := sut.ReloadPlugins()
					sdkPlugin, _ := plugins[0].(sdk.SDKPlugin)
					Expect(sdkPlugin.PluginName()).Should(Equal("sdk-01"))
				})
			})

			When("calling Current() on first SDKPlugin", func() {
				var stdoutLogPath string
				BeforeEach(func() {
					stdoutLogPath = pather.ConfigRoot("stdout.log")
					stdoutLog, _ := os.Create(stdoutLogPath)
					os.Stdout = stdoutLog
				})
				AfterEach(func() {
					_ = fs.RemoveAll(stdoutLogPath)
				})

				It("returns the current version", func() {
					plugins, _ := sut.ReloadPlugins()
					sdkPlugin, _ := plugins[0].(sdk.SDKPlugin)
					Expect(sdkPlugin.Current(nil, []string{"1.16.3"})).Should(Succeed())
					b, _ := ioutil.ReadFile(stdoutLogPath)
					Expect(string(b)).Should(ContainSubstring("1.16.3"))
				})
			})

		})

	})

})
