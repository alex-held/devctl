package shell

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestGenerate(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	// var sut *ShellHookGenerator
	var fixture = NewShellHookFixture(g)

	for _, tc := range fixture.TestCases.TestCases {
		g.Describe("GIVEN valid ShellHookConfig", func() {
			g.Before(func() {
				fixture.CurrentTestCase.Bootstrap(func(cfg *ShellHookConfig) *ShellHookConfig {
					templates, err := template.ParseGlob("templates/*.tmpl")
					if err != nil {
						g.Fatalf("failed to setup template files for tests; error=%v", err)
					}
					fmt.Printf("templates: %v\n", templates)
					return cfg
				})
			})

			g.Describe(fmt.Sprintf("GIVEN section == %s", tc.Description()), func() {
				g.BeforeEach(func() {
					fixture.CurrentTestCase.Bootstrap(func(cfg *ShellHookConfig) *ShellHookConfig {
						cfg = cfg.AddSection(func(gen GeneratorGetter) Section {
							return *CreateSDKSection(cfg.GeneratorGetter, "go", "java", "dotnet", "scala")
						})
						cfg = cfg.AddSection(func(getter GeneratorGetter) Section {
							return *NewSection(getter, "devctl", DevCtlSectionConfig{
								Prefix: ".devctl",
							})
						})
						return cfg
					})
				})
			})

			g.JustBeforeEach(func() {
				fixture.CurrentTestCase = tc
			})

			g.AfterEach(func() {
				fixture.Next()
			})

			g.It(tc.Then(), func() {
				sut := fixture.CreateSubjectUnderTest()
				tc.CaptureActualWithError(
					sut.Generate(),
				)
			})
		})
	}
}

/*
func NewTestCaseGroup(section string , tests  ...*testCase) TestCase {
	for _, test := range tests {
		test.grouping = section
	}
	return
}*/

func CreateSDKSection(getter GeneratorGetter, IDs ...string) *Section {
	var sdks []SDK
	for _, sdk := range IDs {
		sdks = append(sdks, TestSdkMap[sdk])
	}
	return NewSection(getter, "sdk", SDKSectionConfig{
		SDKs: sdks,
	})
}

func (cfg *ShellHookConfig) CreateSDKSection(IDs ...string) *Section {
	return CreateSDKSection(cfg.GeneratorGetter, IDs...)
}

func CreateDevctlSection(getter GeneratorGetter) *Section {
	return NewSection(getter, "devctl", DevCtlSectionConfig{
		Prefix: ".devctl",
	})
}

func (cfg *ShellHookConfig) CreateDevctlSection() *Section {
	return CreateDevctlSection(cfg.GeneratorGetter)
}
