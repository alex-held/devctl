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
					fixture.CurrentTestCase.Configure(func(cfg *ShellHookConfig) *ShellHookConfig {
						return cfg.AddOrUpdateSectionForKey("devctl", *NewSection("devctl", DevCtlSectionConfig{
							Prefix: ".devctl",
						}))
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
					input := fixture.CurrentTestCase.Build().Config
					tc.CaptureActualWithError(
						sut.Generate(*input),
					)
				})
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

func CreateSDKSection(IDs ...string) *Section {
	var sdks []SDK
	for _, sdk := range IDs {
		sdks = append(sdks, TestSdkMap[sdk])
	}
	return NewSection("sdk", SDKSectionConfig{
		SDKs: sdks,
	})
}
