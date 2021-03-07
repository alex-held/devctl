package shell

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"

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
	Bootstrap(fn ...ShellHookApplyFn) TestCase
	Configure(fn ...ShellHookApplyFn) TestCase
	CaptureActual(actual interface{}) TestCase
	CaptureActualWithError(actual interface{}, err error) TestCase
	GetSubjectUnderTest() *ShellHookGenerator
	AssertAndReport()
}

func (t *testCase) String() string { return fmt.Sprintf("TestCase: %v\n\n", *t) }

type testInput struct {
	Inputs    interface{}
	Generator *ShellHookGenerator
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
	SubjectUnderTest                               *ShellHookGenerator
	input                                          TestInput
	Output                                         *TestOutput
	Matchers                                       []gomegatypes.GomegaMatcher
	fixture                                        *ShellHookFixture
	grouping                                       string
}

func (t *testCase) GetSubjectUnderTest() *ShellHookGenerator {
	return t.SubjectUnderTest
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
		t.SubjectUnderTest.Config = bootstrap(t.SubjectUnderTest.Config)
	}
	return t
}

func (t *testCase) Configure(fn ...ShellHookApplyFn) TestCase {
	for _, applyFn := range fn {
		t.SubjectUnderTest.Config = applyFn(t.SubjectUnderTest.Config)
	}
	return t
}

func (t *testCase) AssertAndReport() {
	t.fixture.G.Helper()
	Expect(t.Output.Error).ShouldNot(HaveOccurred())

	result, ok := t.Output.Actual.(*ShellHook)
	Expect(ok).Should(BeTrue(), "actual result should be of type *ShellHook\n")

	fmt.Printf("Output: \n%s\n\n", result.ShellScriptString)

	for _, matcher := range t.Matchers {
		Expect(t.Output.Actual).Should(matcher)
	}
	//	return t.fixture.Assert(Expect(t.Output.Actual))
}

func NewTestCase(description, then string, matchers ...gomegatypes.GomegaMatcher) *testCase {
	tt := template.Must(template.ParseGlob("templates/*.tmpl"))

	gen := &ShellHookGenerator{
		Config: NewShellHookConfig(WithTemplates(tt)),
		Out:    &bytes.Buffer{},
		Errors: []error{},
	}
	gen.Config.GeneratorGetter = func() *ShellHookGenerator {
		return gen
	}

	tc := &testCase{
		ThenDescription:  then,
		description:      description,
		SubjectUnderTest: gen,
		Output:           &TestOutput{},
		Matchers:         matchers,
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
	TestCases       TestCaseCollection
	CurrentTestCase TestCase
	Out             *bytes.Buffer
	G               *goblin.G
}

func (f *ShellHookFixture) CaptureActualWithError(actual interface{}, err error) *ShellHookFixture {
	f.CurrentTestCase.CaptureActualWithError(actual, err)
	return f
}

func (f *ShellHookFixture) CaptureActual(actual interface{}) *ShellHookFixture {
	f.CurrentTestCase.CaptureActual(actual)
	return f
}

func (f *ShellHookFixture) Next() { f.TestCases.Next() }

type ShellHookGenerateAssertion struct {
	Matchers []gomegatypes.GomegaMatcher
	Config   *ShellHookConfig
}
type FixtureAsserter interface {
	Assert(accept GomegaAssertion) FixtureAsserter
}

func (a *ShellHookGenerateAssertion) Apply(applyFn func(*ShellHookGenerateAssertion) *ShellHookGenerateAssertion) *ShellHookGenerateAssertion {
	return applyFn(a)
}
func (a *ShellHookGenerateAssertion) WithSetup(applyFn ShellHookApplyFn) *ShellHookGenerateAssertion {
	a.Config = applyFn(a.Config)
	return a
}

func assert(asserts ...gomegatypes.GomegaMatcher) ShellHookGenerateAssertion {
	tt := template.Must(template.ParseGlob("templates/*.tmpl"))
	cfg := NewShellHookConfig(WithTemplates(tt))
	return ShellHookGenerateAssertion{
		Matchers: asserts,
		Config:   cfg,
	}
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
	fixture.CurrentTestCase = tests.Next()
	fixture.TestCases = tests
	return fixture
}

func (f *ShellHookFixture) CreateSubjectUnderTest() *ShellHookGenerator {
	return f.CurrentTestCase.GetSubjectUnderTest()
}

func CreateShellHooksTestData() TestCaseCollection {
	return *NewTestCases(
		// devctl
		NewTestCase(
			"devctl",
			"THEN generates DEVCTL exports",
			matchers.ContainExport("DEVCTL_PREFIX", ".devctl"),
			matchers.ContainExport("DEVCTL_HOME", "$HOME/$DEVCTL_PREFIX")).
			Configure(func(cfg *ShellHookConfig) *ShellHookConfig {
				cfg.CreateDevctlSection()
				return cfg
			}),
		NewTestCase(
			"devctl",
			"THEN generates DEVCTL section header",
			matchers.HaveSectionHeader("DEVCTL")).
			Configure(func(cfg *ShellHookConfig) *ShellHookConfig {
				cfg.CreateDevctlSection()
				return cfg
			}),
		NewTestCase(
			"devctl",
			"THEN generates newlines after exports",
			HaveSuffix("\n\n")).
			Configure(func(cfg *ShellHookConfig) *ShellHookConfig {
				cfg.CreateDevctlSection()
				return cfg
			}),
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
					CreateSDKSection(cfg.GeneratorGetter, "go", "java", "dotnet", "scala"),
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
