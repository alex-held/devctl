package sdk

import (
	"bytes"
	_ "embed"
	"fmt"
	"plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// VALID Tests, but they test the SDKBinder..
// ==============================================================================
//

var _ = Describe("GOSDKPlugin", func() {

	var inputPlugin *plugin.Plugin
	BeforeSuite(func() {
		p, err := plugin.Open("testdata/plugins/sdk-01.so")
		Expect(err).Should(Succeed())
		inputPlugin = p
	})

	Describe("GoSDKPluginBinder", func() {
		var sut SDKPluginBinderFn
		Describe("Bind()", func() {

			var out *bytes.Buffer
			var _ error
			var sdkPlugin SDKPlugin

			BeforeEach(func() {
				out = &bytes.Buffer{}
			})

			It("should not return an error", func() {
				_, err := sut.Bind(inputPlugin)
				Expect(err).Should(Succeed())
			})

			It("should bind PluginName()", func() {
				Expect(sdkPlugin.PluginName()).Should(Equal("sdk-01"))
			})

			It("should bind SetStdoutName()", func() {
				sdkPlugin, err := sut.Bind(inputPlugin)
				Expect(err).Should(Succeed())
				sdkPlugin.SetStdout(out)
				actualName := sdkPlugin.PluginName()
				Expect(out.String()).Should(
					And(
						Equal("sdk-01"),
						Equal(actualName),
					),
				)
			})

			It("should bind Download()", func() {
				sdkPlugin, err := sut.Bind(inputPlugin)
				Expect(err).Should(Succeed())
				sdkPlugin.SetStdout(out)
				Expect(sdkPlugin.Download(nil, []string{"1.16.3"})).Should(Succeed())
				Expect(out.String()).Should(ContainSubstring(fmt.Sprintf("Download + %v","1.16.3")))
			})

			It("should bind Use()", func() {
				sdkPlugin, err := sut.Bind(inputPlugin)
				Expect(err).Should(Succeed())
				sdkPlugin.SetStdout(out)
				Expect(sdkPlugin.Use(nil, []string{"1.16.3"})).Should(Succeed())
				Expect(out.String()).Should(ContainSubstring(fmt.Sprintf("Use + %v","1.16.3")))
			})

			It("should bind List()", func() {
				sdkPlugin, err := sut.Bind(inputPlugin)
				Expect(err).Should(Succeed())
				sdkPlugin.SetStdout(out)
				Expect(sdkPlugin.List(nil, []string{"1.16.3"})).Should(Succeed())
				Expect(out.String()).Should(ContainSubstring(fmt.Sprintf("List + %v","1.16.3")))
			})

			It("should bind Current()", func() {
				sdkPlugin, err := sut.Bind(inputPlugin)
				Expect(err).Should(Succeed())
				sdkPlugin.SetStdout(out)
				Expect(sdkPlugin.Current(nil, []string{"1.16.3"})).Should(Succeed())
				Expect(out.String()).Should(ContainSubstring(fmt.Sprintf("Current + %v","1.16.3")))
			})

			It("should bind Install()", func() {
				sdkPlugin, err := sut.Bind(inputPlugin)
				Expect(err).Should(Succeed())
				sdkPlugin.SetStdout(out)
				Expect(sdkPlugin.Install(nil, []string{"1.16.3"})).Should(Succeed())
				Expect(out.String()).Should(ContainSubstring(fmt.Sprintf("Install + %v","1.16.3")))
			})
		})
	})
})
