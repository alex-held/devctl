package shell

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	matchers2 "github.com/alex-held/devctl/pkg/testutils/matchers"
)

type ShellHookFixture struct {
	Cfg       *HookConfig
	Templates *template.Template
	Gen       *HookGenerator
	Out       *bytes.Buffer
	G         *goblin.G
}

type SetupFixtureFn func(fixture *ShellHookFixture)
type ShellHookGenerateAssertion struct {
	Matchers []types.GomegaMatcher
	Config   *HookConfig
}

func assert(asserts ...types.GomegaMatcher) ShellHookGenerateAssertion {
	cfg := NewShellHookConfig()
	return ShellHookGenerateAssertion{
		Matchers: asserts,
		Config:   cfg,
	}
}

type Option func(*HookConfig) *HookConfig

func WithTemplates(t *template.Template) Option {
	return func(cfg *HookConfig) *HookConfig {
		cfg.Templates = t
		return cfg
	}
}

func WithSection(initializeSection UninitializedSection) Option {
	return func(cfg *HookConfig) *HookConfig {
		initializedSection := initializeSection(cfg.root)
		return cfg.AddOrUpdateSection(initializedSection)
	}
}

func WithSections(sections ...UninitializedSection) Option {
	return func(cfg *HookConfig) *HookConfig {
		for _, uninitializedSection := range sections {
			cfg = WithSection(uninitializedSection)(cfg)
		}
		return cfg
	}
}

var defaults = []Option{
	WithTemplates(template.Must(template.ParseGlob("templates/*.tmpl"))),
}

func NewShellHookConfig(opts ...Option) *HookConfig {
	cfg := &HookConfig{
		Templates: nil,
		Sections:  Sections{},
	}
	cfg.root = NewShellHookRootNode(cfg)

	for _, opt := range defaults {
		cfg = opt(cfg)
	}

	for _, opt := range opts {
		cfg = opt(cfg)
	}

	return cfg
}

func NewShellHookRootNode(cfg *HookConfig) *rootNode { // nolint: golint
	root := &rootNode{config: cfg}
	root.node = &node{
		parent:   root.node,
		isRooted: true,
	}
	cfg.root = root
	root.config = cfg
	return root
}

// type ShellHookGenerateAssertionFn func(fix *ShellHookFixture) []ShellHookGenerateAssertion

type ShellHookGenerateTestCase struct {
	Description           string
	FixtureBeforeAllSetup SetupFixtureFn
	Expected              []ShellHookGenerateAssertion
	TestCases             map[string]ShellHookGenerateAssertion
}

func (f *ShellHookFixture) Assert(actual interface{}, err error, tc ShellHookGenerateAssertion) {
	Expect(err).ShouldNot(HaveOccurred())
	for _, assertion := range tc.Matchers {
		Expect(actual).Should(assertion)
	}
}

func (f *ShellHookFixture) ConfigureFixture(setupFn SetupFixtureFn) *ShellHookFixture {
	setupFn(f)
	return f
}

func Setup(g *goblin.G) *ShellHookFixture {
	return NewShellHookFixture(g)
}

func NewShellHookFixture(g *goblin.G) *ShellHookFixture {
	templates, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		g.Fatalf("failed to setup template files for tests; error=%v", err)
	}

	fixture := &ShellHookFixture{
		G:         g,
		Cfg:       NewShellHookConfig(WithTemplates(templates), WithSections()),
		Templates: templates,
		Out:       &bytes.Buffer{},
	}

	fixture.Gen = &HookGenerator{
		Out:       fixture.Out,
		Templates: fixture.Templates,
		//		OutputPath: "$DEVCTL_HOME/shell/hooks/zsh_hook",
	}

	return fixture
}

func TestGenerate(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	var sut *HookGenerator
	var fixture *ShellHookFixture

	var data = map[string]ShellHookGenerateTestCase{
		"devctl": {
			FixtureBeforeAllSetup: func(f *ShellHookFixture) {
				f.Cfg.AddOrUpdateSectionForKey("devctl", *NewSection("devctl", DevCtlSectionConfig{
					Prefix: ".devctl",
				}))
			},
			TestCases: map[string]ShellHookGenerateAssertion{
				"THEN generates DEVCTL section header": assert(
					matchers2.HaveSectionHeader("DEVCTL"),
				),
				"THEN generates DEVCTL exports": assert(
					matchers2.ContainExport("DEVCTL_PREFIX", ".devctl"),
					matchers2.ContainExport("DEVCTL_HOME", "$HOME/$DEVCTL_PREFIX"),
				),
				"THEN generates newlines after exports": assert(
					HaveSuffix("\n\n"),
				),
			},
		},
	}

	g.Describe("GIVEN valid ShellHookConfig", func() {
		for section, testCase := range data {
			g.Describe(fmt.Sprintf("GIVEN section == %s", section), func() {
				for description, test := range testCase.TestCases {
					g.BeforeEach(func() {
						fixture = Setup(g)
						fixture.ConfigureFixture(testCase.FixtureBeforeAllSetup)
						sut = fixture.Gen
					})

					g.It(description, func() {
						actual, err := sut.Generate(*fixture.Cfg)
						fixture.Assert(actual, err, test)
					})
				}
			})
		}
	})
}
