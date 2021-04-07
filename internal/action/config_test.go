package action

import (
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"

	"github.com/alex-held/devctl/internal/config"
)

func (f *ActionTestFixture) AssertConfig(assertionsFn configAssertionFn) {
	cfg, err := f.actions.Config.Load()
	if err != nil {
		f.g.Fatalf("failed to load config from fs; error=%v\n", err)
	}
	assertions := assertionsFn(cfg)
	for _, a := range assertions {
		Expect(a.actual).Should(a.matcher)
	}
}

func (f *ActionTestFixture) SetupConfig(configFn func(*config.Config)) {
	c := config.NewBlankConfig()
	configFn(c)

	b, err := yaml.Marshal(*c)
	if err != nil {
		f.g.Fatalf("failed to marshal config file while setup; error=%v", err)
	}
	err = afero.WriteFile(f.fs, f.pather.ConfigFilePath(), b, 0700)
	if err != nil {
		f.g.Fatalf("failed to write marshaled config file while setup; error=%v", err)
	}
}

type configAssertionFn func(c *config.Config) []configAssertion

type configAssertion struct {
	actual  interface{}
	matcher types.GomegaMatcher
}

func TestConfig_SetCurrentSDK(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Config", func() {
		const sdkPath = "sdks/scala/2.13.4"
		const sdk = "scala"
		const version = "2.13.4"

		var fixture *ActionTestFixture
		var sut *Config

		g.JustBeforeEach(func() {
			fixture = SetupFixture(g)
			sut = fixture.actions.Config
		})

		g.Describe("GIVEN no config file", func() {
			g.It("THEN creates config file with sdk and set installation + current", func() {
				err := sut.SetCurrentSdk(sdk, version, sdkPath)
				Expect(err).Should(BeNil())

				fixture.AssertConfig(func(c *config.Config) []configAssertion {
					return []configAssertion{
						{c.Sdks[sdk], Not(BeNil())},
						{c.Sdks[sdk].Current, Equal(version)},
						{c.Sdks[sdk].Installations[version], Equal(sdkPath)},
						{c.Sdks[sdk].Installations, HaveLen(1)},
						{c.Sdks[sdk].SDK, Equal(sdk)},
					}
				})
			})
		})

		g.Describe("GIVEN blank config file", func() {
			g.It("THEN add sdk with installation and current", func() {
				fixture.SetupConfig(func(c *config.Config) {})
				err := sut.SetCurrentSdk(sdk, version, sdkPath)
				Expect(err).Should(BeNil())

				fixture.AssertConfig(func(c *config.Config) []configAssertion {
					return []configAssertion{
						{c.Sdks[sdk], Not(BeNil())},
						{c.Sdks[sdk].Current, Equal(version)},
						{c.Sdks[sdk].Installations[version], Equal(sdkPath)},
						{c.Sdks[sdk].Installations, HaveLen(1)},
						{c.Sdks[sdk].SDK, Equal(sdk)},
					}
				})
			})
		})

		g.Describe("GIVEN config file contains sdk", func() {
			const otherVersion = "2.13.0"
			const otherSdkPath = "sdks/scala/2.13.0"

			g.Describe("AND version not installed", func() {
				g.It("THEN adds installation + sets current for sdk", func() {
					fixture.SetupConfig(func(c *config.Config) {
						c.Sdks[sdk] = config.SdkConfig{
							SDK:     sdk,
							Current: otherVersion,
							Installations: map[string]string{
								otherVersion: otherSdkPath,
							},
						}
					})

					err := sut.SetCurrentSdk(sdk, version, sdkPath)
					Expect(err).Should(BeNil())
					fixture.AssertConfig(func(c *config.Config) []configAssertion {
						return []configAssertion{
							{*c, Not(BeNil())},
							{c.Sdks, HaveLen(1)},
							{c.Sdks[sdk].Installations, HaveLen(2)},
							{c.Sdks[sdk].Installations[version], Equal(sdkPath)},
							{c.Sdks[sdk].SDK, Equal(sdk)},
							{c.Sdks[sdk].Current, Equal(version)},
						}
					})
				})
			})

			g.Describe("AND version already installed", func() {
				g.Describe("AND current set to different version", func() {
					g.It("THEN", func() {
						fixture.SetupConfig(func(c *config.Config) {
							c.Sdks[sdk] = config.SdkConfig{
								SDK:     sdk,
								Current: otherVersion,
								Installations: map[string]string{
									otherVersion: otherSdkPath,
									version:      sdkPath,
								},
							}
						})

						err := sut.SetCurrentSdk(sdk, version, sdkPath)
						Expect(err).Should(BeNil())
						fixture.AssertConfig(func(c *config.Config) []configAssertion {
							return []configAssertion{
								{*c, Not(BeNil())},
								{c.Sdks, HaveLen(1)},
								{c.Sdks[sdk].Installations, HaveLen(2)},
								{c.Sdks[sdk].Installations[version], Equal(sdkPath)},
								{c.Sdks[sdk].SDK, Equal(sdk)},
								{c.Sdks[sdk].Current, Equal(version)},
							}
						})
					})
				})

				g.Describe("AND current set to same version with different path", func() {
					const otherSdkPath = "other/path/to/sdks/scala/2.13.0"

					g.It("THEN update installation path", func() {
						fixture.SetupConfig(func(c *config.Config) {
							c.Sdks[sdk] = config.SdkConfig{
								SDK:     sdk,
								Current: version,
								Installations: map[string]string{
									version: otherSdkPath,
								},
							}
						})

						err := sut.SetCurrentSdk(sdk, version, sdkPath)
						Expect(err).Should(BeNil())
						fixture.AssertConfig(func(c *config.Config) []configAssertion {
							return []configAssertion{
								{*c, Not(BeNil())},
								{c.Sdks, HaveLen(1)},
								{c.Sdks[sdk].Installations, HaveLen(1)},
								{c.Sdks[sdk].Installations[version], Equal(sdkPath)},
								{c.Sdks[sdk].SDK, Equal(sdk)},
								{c.Sdks[sdk].Current, Equal(version)},
							}
						})
					})
				})
			})
		})
	})
}
