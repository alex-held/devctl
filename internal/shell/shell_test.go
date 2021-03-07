package shell

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/stretchr/testify/require"

	"github.com/alex-held/devctl/internal/testutils/matchers"
)

type SDKMap map[string]SDK

var TestSdkMap = SDKMap{
	"go": {
		SDK:  "go",
		Path: "$DEVCTL_HOME/sdks/go/current",
	},
	"java": {
		SDK:  "java",
		Path: "$DEVCTL_HOME/sdks/java/current",
	},
	"dotnet": {
		SDK:  "dotnet",
		Path: "$DEVCTL_HOME/sdks/dotnet/current",
	},
	"scala": {
		SDK:  "scala",
		Path: "$DEVCTL_HOME/sdks/scala/current",
	},
}

type TestCase interface {
	//	gotypes.Type
	//	GetMatchers() []gomegatypes.GomegaMatcher
	// GetInputArgs() []interface{}
	//	Givens() []string
	//	Whens() []string
	//	Ands() []string
	Then() string
	Description() string
	Build() struct {
		Config        *ShellHookConfig
		FixtureConfig FixtureConfig
	}
	Bootstrap(fn ...ShellHookApplyFn) TestCase
	Configure(fn ...ShellHookApplyFn) TestCase
	CaptureActual(actual interface{}) TestCase
	CaptureActualWithError(actual interface{}, err error) TestCase
	AssertAndReport()
}

func (t *testCase) String() string { return fmt.Sprintf("TestCase: %v\n\n", *t) }

func (t *testCase) Build() struct {
	Config        *ShellHookConfig
	FixtureConfig FixtureConfig
} {
	return struct {
		Config        *ShellHookConfig
		FixtureConfig FixtureConfig
	}{
		Config:        t.builder.config,
		FixtureConfig: t.builder.fixtureConfig,
	}
}

type Builder struct {
	builder          testInput
	BaseConfig       interface{}
	lazyConfigSetups []func(interface{}) interface{}
	isLocked         bool
}

type FixtureConfig map[string]interface{}

type testInput struct {
	Inputs        interface{}
	config        *ShellHookConfig
	fixtureConfig FixtureConfig
}

type TestInput interface {
	Config() interface{}
	Input() interface{}
}

type testCase struct {
	Common struct {
		Out bytes.Buffer
	}
	ThenDescription, description                   string
	GivenCollection, WhenCollection, AndCollection []string
	builder                                        *testInput
	input                                          TestInput
	Output                                         *TestOutput
	Matchers                                       []gomegatypes.GomegaMatcher
	fixture                                        *ShellHookFixture
	grouping                                       string
}

func (t *testCase) Description() string                       { return t.description }
func (t *testCase) Group() string                             { return t.grouping }
func (t *testCase) GetMatchers() []gomegatypes.GomegaMatcher  { return t.Matchers }
func (t *testCase) Givens() []string                          { return t.GivenCollection }
func (t *testCase) Whens() []string                           { return t.WhenCollection }
func (t *testCase) Ands() []string                            { return t.AndCollection }
func (t *testCase) Then() string                              { return t.ThenDescription }
func (t *testCase) CaptureActual(actual interface{}) TestCase { t.Output.Actual = actual; return t }
func (t *testCase) CaptureActualWithError(actual interface{}, err error) TestCase {
	t.Output = &TestOutput{Error: err, Actual: actual}
	return t
}

func (t *testCase) Bootstrap(bootstrapFn ...ShellHookApplyFn) TestCase {
	for _, bootstrap := range bootstrapFn {
		c := t.builder.config
		bootstrap(c)
	}
	return t
}

func (t *testCase) Configure(fn ...ShellHookApplyFn) TestCase {
	for _, applyFn := range fn {
		t.builder.config = applyFn(t.builder.config)
	}
	return t
}
func (t *testCase) AssertAndReport() {
	t.fixture.G.Helper()
	Expect(t.Output.Error).ShouldNot(HaveOccurred())
	for _, matcher := range t.Matchers {
		Expect(t.Output.Actual).Should(matcher)
	}
	//	return t.fixture.Assert(Expect(t.Output.Actual))
}

func NewTestCase(description, then string, matchers ...gomegatypes.GomegaMatcher) *testCase {
	tc := &testCase{
		ThenDescription: then,
		description:     description,
		builder: &testInput{
			fixtureConfig: FixtureConfig{},
			config:        NewShellHookConfig(),
		},
		input:    nil,
		Output:   &TestOutput{},
		Matchers: matchers,
	}
	return tc
}

type TestOutput struct {
	Error  error
	Actual interface{}
}

type TestCases []TestCase

func (c TestCaseCollection) Next() TestCase {
	c.CurrentPosition++
	idx := c.CurrentPosition
	var next = (c.TestCases)[idx]
	return next
}

type TestCaseCollection struct {
	TestCases
	CurrentPosition int
	fixtureGetter   func() *ShellHookFixture
}

func NewTestCases(tcs ...TestCase) *TestCaseCollection {
	// tc1 := tcs[0]
	// tct := tc1.Underlying()
	tcSlice := append(TestCases{}, tcs...)

	tcc := &TestCaseCollection{
		// Slice:           gotypes.NewSlice(gotypes.Id(gotypes.NewPackage(""), "")),
		TestCases:       tcSlice,
		CurrentPosition: -1,
	}
	return tcc
}

/*
func (tcc *TestCaseCollection) Underlying() gotypes.Type {
	var testCaseInstance TestCase = &testCase{}
	var testCaseType = testCaseInstance.Underlying()
	return gotypes.NewSlice(testCaseType)
}

func (tcc *TestCaseCollection) Len() int { return len(*tcc.TestCases) }
func (tcc *TestCaseCollection) Cap() int { return cap(*tcc.TestCases) }

func (tc *TestCases) append(others ...TestCase) *TestCases {
	*tc = append(*tc, others...)
	return tc
}

func (tcc *TestCaseCollection) Get(id int) TestCase {
	var tc TestCases
	_ = copy(tc, *tcc.TestCases)
	//goland:noinspection GoNilness
	return (tc)[id]
}
func (tcc *TestCaseCollection) Set(id int, testCase TestCase) {
	(*tcc.TestCases)[id] = testCase
}

func (tcc *TestCaseCollection) GetNext() (tc TestCase, done bool) {
	currentIdx := tcc.CurrentPosition
	nextIdx := currentIdx + 1
	tcc.CurrentPosition = nextIdx
	tc = (*tcc.TestCases)[nextIdx]
	return tc, tc == nil
}
*/

type ShellHookFixture struct {
	TestCases TestCaseCollection
	// Cfg              *ShellHookConfig
	CurrentTestCase TestCase
	// SubjectUnderTest *ShellHookGenerator
	Out *bytes.Buffer
	G   *goblin.G
}

func (f *ShellHookFixture) CaptureActualWithError(actual interface{}, err error) *ShellHookFixture {
	f.CurrentTestCase.CaptureActualWithError(actual, err)
	return f
}

func (f *ShellHookFixture) CaptureActual(actual interface{}) *ShellHookFixture {
	f.CurrentTestCase.CaptureActual(actual)
	return f
}

type FixtureConfigFn func(fixture *FixtureConfig) *FixtureConfig

type ConfigurableFixtureAsserter interface {
	// Configure(applyFn ShellHookApplyFn) ConfigurableFixtureAsserter
	Assert(accept GomegaAssertion) ConfigurableFixtureAsserter
}

func (f *ShellHookFixture) Next() { f.TestCases.Next() }

type ShellHookGenerateAssertion struct {
	Matchers []gomegatypes.GomegaMatcher
	Config   *ShellHookConfig
}
type FixtureAsserter interface {
	Assert(accept GomegaAssertion) FixtureAsserter
}

type ShellHookFixtureAssertion interface {
	GetConfig() ShellHookConfig
	Assert(accept GomegaAssertion)
}

func (a *ShellHookGenerateAssertion) Apply(applyFn func(*ShellHookGenerateAssertion) *ShellHookGenerateAssertion) *ShellHookGenerateAssertion {
	return applyFn(a)
}
func (a *ShellHookGenerateAssertion) WithSetup(applyFn ShellHookApplyFn) *ShellHookGenerateAssertion {
	a.Config = applyFn(a.Config)
	return a
}

func assert(asserts ...gomegatypes.GomegaMatcher) ShellHookGenerateAssertion {
	cfg := NewShellHookConfig()
	return ShellHookGenerateAssertion{
		Matchers: asserts,
		Config:   cfg,
	}
}

type ShellHookGenerateTestCase struct {
	Description string
	Expected    []ShellHookGenerateAssertion
	TestCases   map[string]ShellHookGenerateAssertion
}

func NewShellHookFixture(g *goblin.G) *ShellHookFixture {
	rep := &goblin.DetailedReporter{}
	g.SetReporter(rep)
	rep.SetTextFancier(goblin.TextFancier(&goblin.TerminalFancier{}))
	tests := CreateShellHooksTestData()
	fixture := &ShellHookFixture{
		//	Reporter:        goblin.Reporter(rep),
		// CurrentTestCase:  tests.Next(),
		Out: &bytes.Buffer{},
		G:   g,
	}
	fixture.TestCases.fixtureGetter = func() *ShellHookFixture {
		return fixture
	}
	fixture.CurrentTestCase = tests.Next()
	fixture.TestCases = tests
	return fixture
}

func (f *ShellHookFixture) CreateSubjectUnderTest() *ShellHookGenerator {
	cfg := f.CurrentTestCase.Build()
	return &ShellHookGenerator{
		Out:       f.Out,
		Templates: cfg.Config.Templates,
		// Templates: f.FixtureConfig["templates"].(*template.Template),
	}
}

// Iterate calls the f function with n = 1, 2, and 3.
func Iterate(f func(next TestCase)) {
	i := 0
	tc := CreateShellHooksTestData().TestCases

	max := len(tc) - 1
	for i = 1; i <= max; i++ {
		next := tc[i]
		f(next)
	}
}

func TestCreateIterator(t *testing.T) {
	var i = 0
	Iterate(func(next TestCase) {
		i++
		fmt.Printf("[%d] <%T> | testcase=%#v\n\n", i, next, next)
	})
	require.Equal(t, 6, i)
}

func CreateShellHooksTestData() TestCaseCollection {
	return *NewTestCases(
		// devctl
		NewTestCase(
			"devctl",
			"THEN generates DEVCTL exports",
			matchers.ContainExport("DEVCTL_PREFIX", ".devctl"),
			matchers.ContainExport("DEVCTL_HOME", "$HOME/$DEVCTL_PREFIX")),
		NewTestCase(
			"devctl",
			"THEN generates DEVCTL section header",
			matchers.HaveSectionHeader("DEVCTL")),
		NewTestCase(
			"devctl",
			"THEN generates newlines after exports",
			HaveSuffix("\n\n")),

		// sdk
		NewTestCase(
			"sdk",
			"THEN generates exports for all sdks",
			matchers.HaveSectionHeader("sdk"),
			matchers.ContainExport("GO_HOME", "$HOME/sdks/go/current"),
			matchers.ContainExport("JAVA_HOME", "$HOME/sdks/java/current"),
			matchers.ContainExport("DOTNET_HOME", "$HOME/sdks/dotnet/current"),
			matchers.ContainExport("SCALA_HOME", "$HOME/sdks/scala/current")).
			Configure(func(cfg *ShellHookConfig) *ShellHookConfig {
				return cfg.AddOrUpdateSection(
					CreateSDKSection("go", "java", "dotnet", "scala"),
				)
			}),

		NewTestCase(
			"sdk",
			"THEN generates DEVCTL section header",
			matchers.HaveSectionHeader("sdk")),
		NewTestCase(
			"sdk",
			"THEN generates newlines after exports",
			HaveSuffix("\n\n")),
	)
}
