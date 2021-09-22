//go:generate pluggen gen -o "./testdata/plugins/sdk-01.so" -p "plugins/sdk-01" --pkg devctl

package sdk

/*

import (
	_ "embed"
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
)

var (
	//go:embed testdata/plugins/sdk-01.so
	Sdk01Plugin  []byte
	ExpectedArgs = []string{"1.16.3"}
)

var _ = Describe("SDKPlugin", func() {
	var fs afero.Fs
	var pather devctlpath.Pather
	var sdkPluginDir, pluginPath string

	BeforeEach(func() {
		fs = afero.NewOsFs()
		pather = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			tmp, _ := afero.TempDir(fs, "", "devctl-SDKPluginTests-*")
			_ = fs.MkdirAll(tmp, os.ModePerm)
			return tmp
		}))
		sdkPluginDir = pather.ConfigRoot("plugins")
		pluginPath = path.Join(sdkPluginDir, "sdk-01.so")
	})

	Context("LoadSDKPlugin", func() {

		When("file system does not contain plugins", func() {
			var loaderSut SDKPluginLoaderFn

			It("returns error when loading", func() {
				_, err :=  loaderSut.LoadSDKPlugin("no-path")
				Expect(err).ShouldNot(Succeed())
			})

			It("returns error when loading", func() {
				unboundPlugin, _ :=  loaderSut.LoadSDKPlugin("no-path")
				Expect(unboundPlugin).Should(BeNil())
			})
		})

	When("file system one SDKPlugin", func() {
			var loaderSut SDKPluginLoaderFn

			BeforeEach(func() {
				err := os.MkdirAll(path.Dir(pluginPath), os.ModePerm)
				if err != nil {
					panic(err)
				}
				err = ioutil.WriteFile(pluginPath, Sdk01Plugin, os.ModePerm)
				if err != nil {
					panic(err)
				}
			})

			It("does not return error", func() {
				_, err := loaderSut.LoadSDKPlugin(pluginPath)
				Expect(err).Should(Succeed())
			})

			It("returns SDKPlugin when loading", func() {
				actual, _ := loaderSut.LoadSDKPlugin(pluginPath)
				Expect(actual).ShouldNot(BeNil())
			})

		})
	})
})

//var _ = Describe("SDKPluginSuite", func() {
//	var pather devctlpath.Pather
//	var tmpRoot, pluginPath string
//
//	BeforeEach(func() {
//		pather = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
//			tmpRoot, _ = os.MkdirTemp("", "devctl.SDKPluginSuite")
//			println("TMP_ROOT_PATH=" + tmpRoot)
//			return tmpRoot
//		}))
//		pluginPath = pather.ConfigRoot("plugins", "sdk", "sdk-01.so")
//	})
//
//	Context("check whether the binding from .so to SDKPlugin works as expected", func() {
//		var output *bytes.Buffer
//		var ExpectedArgs []string
//		var ctx context.Context
//		var fs afero.Fs
//
//		BeforeEach(func() {
//			ctx = context.TODO()
//			fs = afero.NewOsFs()
//			ExpectedArgs = []string{"1.16.3"}
//			output = &bytes.Buffer{}
//		})
//
//		AfterEach(func() {
//			_ = fs.RemoveAll(pluginPath)
//			_ = fs.RemoveAll(tmpRoot)
//		})
//
//		When("PluginName()", func() {
//
//			BeforeEach(func() {
//				data, _ := afero.ReadFile(fs, "testdata/plugins/sdk/sdk-01.so")
//				_ = afero.WriteFile(fs, pluginPath, data, os.ModePerm)
//			})
//			It("returns the name of the loaded plugin", func() {
//				sut, _ := sdk.LoadSDKPlugin(pluginPath)
//				name := sut.PluginName()
//				Expect(name).Should(Equal("sdk-01"))
//			})
//
//			It("returns no error", func() {
//				_, err := sdk.LoadSDKPlugin(pluginPath)
//				Expect(err).Should(Succeed())
//			})
//		})
//
//		When("Current()", func() {
//
//			It("prints the args to stdout", func() {
//				sut, err := sdk.LoadSDKPlugin(pluginPath)
//				Expect(err).Should(Succeed())
//				Expect(sut.SetStdout(output)).Should(Succeed())
//				_ = sut.Current(ctx, ExpectedArgs)
//				println("OUTPUT: " + output.String())
//				Expect(output).Should(ContainSubstring("1.16.3"))
//			})
//
//			It("returns no error", func() {
//				sut, err := sdk.LoadSDKPlugin(pluginPath)
//				Expect(err).Should(Succeed())
//				Expect(sut.SetStdout(output)).Should(Succeed())
//				Expect(sut.Current(ctx, ExpectedArgs)).Should(Succeed())
//				println("OUTPUT: " + output.String())
//			})
//		})
//	})
//
//	Context("SDKPlugin public API", func() {
//
//		When("no plugin at filepath", func() {
//			It("returns an error", func() {
//				_, err := sdk.LoadSDKPlugin(pluginPath)
//				Expect(err).ShouldNot(Succeed())
//			})
//
//			It("returns an error", func() {
//				sut, _ := sdk.LoadSDKPlugin(pluginPath)
//				Expect(sut).Should(BeNil())
//			})
//		})
//
//		When("plugin is at filepath", func() {
//
//		})
//	})
//})
//
//func copyPluginToTmp(pluginPath string) {
//	p, _ := ioutil.ReadFile(path.Join("testdata", "plugins", "sdk", "sdk-01.so"))
//	_ = afero.NewOsFs().MkdirAll(path.Dir(pluginPath), os.ModePerm)
//	_ = ioutil.WriteFile(pluginPath, p, os.ModePerm)
//}



*/
