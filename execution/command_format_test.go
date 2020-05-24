package execution

import (
	"fmt"
	. "github.com/alex-held/dev-env/manifest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"path"
	"strings"
	"testing"
)

type InstructingFormatTest struct {
	Instructing interface{}
	Expected    string
	Output      *strings.Builder
	Executor    *CommandExecutor
	FS          *afero.Fs
	Vars        map[string]interface{}
}

func (test *InstructingFormatTest) MkDir(dirs ...string) {
	for _, dir := range dirs {
		fs := test.Executor.FS
		_ = (*fs).MkdirAll(dir, 0644)
	}
}

func NewInstructingFormatTest(i interface{}, expected string, after func(t InstructingFormatTest)) InstructingFormatTest {
	sb := strings.Builder{}

	m := Manifest{}
	switch instruction := i.(type) {
	case DevEnvCommand:
		m.Instructions = Instructions{
			Step{
				Command: &instruction,
			},
		}
	case Pipe:
		m.Instructions = Instructions{
			Step{
				Pipe: instruction.Commands,
			},
		}
	case Link:
		m.Links = []Link{instruction}
	}

	executor := NewCommandExecutor(&m, func(str string) {
		sb.WriteString(str)
	})
	executor.Options.DryRun = true
	fs := afero.NewMemMapFs()
	executor.FS = &fs

	test := InstructingFormatTest{
		Instructing: i,
		Expected:    expected,
		Output:      &sb,
		Executor:    executor,
		FS:          &fs,
		Vars:        map[string]interface{}{},
	}
	after(test)
	return test
}

func (test *InstructingFormatTest) Run(t *testing.T) {
	out, err := test.Executor.Execute()
	assert.NoError(t, err)
	actual := test.Output.String()
	fmt.Printf("<EXPECTED>\n%s-\n", test.Expected)
	fmt.Printf("<ACTUAL>\n%s-\n", actual)
	fmt.Printf("<OUT>\n%s-\n", out)
	assert.Equal(t, test.Expected, actual)
}

func TestLinkCommand_Format_When_Target_Directory_Not_Exists(t *testing.T) {

	link := Link{
		Source: "/source/dotnet/host/fxr",
		Target: "/target/dotnet/host/fxr",
	}

	test := NewInstructingFormatTest(link, fmt.Sprintf(`[Link]    
Creating target directory %s
Creating link %s -> %s

`, path.Dir(link.Target), link.Source, link.Target), func(tst InstructingFormatTest) {})

	test.Run(t)
}

func TestLinkCommand_Format_When_Target_Directory_Already_Exists(t *testing.T) {

	link := Link{
		Source: "/source/dotnet/host/fxr",
		Target: "/target/dotnet/host/fxr",
	}

	test := NewInstructingFormatTest(link, fmt.Sprintf(`[Link]    
Creating link %s -> %s

`, link.Source, link.Target), func(tst InstructingFormatTest) {
		tst.MkDir("/target/dotnet/host")
	})

	test.Run(t)
}

func TestCommandExecutor_Execute(t *testing.T) {
	manifest := GetTestManifest()

	expected := `[Command]
Executing command: 'ls -a /Users/dev/temp/usr/local/share/dotnet'

[Command]
Executing command: 'rm -rdf /Users/dev/temp/usr/local/share/dotnet'

[Command]
Executing command: 'mkdir -p /Users/dev/temp/usr/local/share/dotnet/host'

[Command]
Executing command: 'rm -rdf /Users/dev/temp/usr/local/share/dotnet/host'

[Link]    
Creating target directory /Users/dev/temp/usr/local/share/dotnet/host
Creating link /Users/dev/.dev-env/sdk/dotnet/3.1.202/host/fxr -> /Users/dev/temp/usr/local/share/dotnet/host/fxr

[Link]    
Creating target directory /Users/dev/temp/usr/local/share/dotnet/sdk
Creating link /Users/dev/.dev-env/sdk/dotnet/3.1.202/sdk/3.1.202 -> /Users/dev/temp/usr/local/share/dotnet/sdk/3.1.202

[Link]    
Creating target directory /Users/dev/temp/usr/local/share/dotnet/shared/Microsoft.NETCore.App
Creating link /Users/dev/.dev-env/sdk/dotnet/3.1.202/shared/Microsoft.NETCore.App -> /Users/dev/temp/usr/local/share/dotnet/shared/Microsoft.NETCore.App/3.1.202

[Link]    
Creating target directory /Users/dev/temp/usr/local/share/dotnet/shared/Microsoft.AspNetCore.All
Creating link /Users/dev/.dev-env/sdk/dotnet/3.1.202/shared/Microsoft.AspNetCore.All -> /Users/dev/temp/usr/local/share/dotnet/shared/Microsoft.AspNetCore.All/3.1.202

[Link]    
Creating target directory /Users/dev/temp/usr/local/share/dotnet/shared/Microsoft.AspNetCore.App
Creating link /Users/dev/.dev-env/sdk/dotnet/3.1.202/shared/Microsoft.AspNetCore.App -> /Users/dev/temp/usr/local/share/dotnet/shared/Microsoft.AspNetCore.App/3.1.202

`
	var output = ""

	executor := NewCommandExecutor(manifest, func(str string) {
		output = output + str
	})

	_, _ = executor.Execute()
	assert.Equal(t, expected, output)
}

func GetTestManifest() *Manifest {
	m := Manifest{
		Version: "3.1.202",
		SDK:     "dotnet",
		Variables: Variables{
			{Key: "url", Value: "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz"},
			{Key: "install-root", Value: "[[_sdks]]/[[sdk]]/[[version]]"},
			{Key: "link-root", Value: "/Users/dev/temp/usr/local/share/dotnet"},
		},
		Instructions: Instructions{
			Step{
				Command: &DevEnvCommand{
					Command: "ls",
					Args:    []string{"-a", "/Users/dev/temp/usr/local/share/dotnet"},
				},
			},
			Step{
				Command: &DevEnvCommand{
					Command: "rm",
					Args:    []string{"-rdf", "/Users/dev/temp/usr/local/share/dotnet"},
				},
			},
			Step{
				Command: &DevEnvCommand{
					Command: "mkdir",
					Args:    []string{"-p", "/Users/dev/temp/usr/local/share/dotnet/host"},
				},
			},
			Step{
				Command: &DevEnvCommand{
					Command: "rm",
					Args:    []string{"-rdf", "/Users/dev/temp/usr/local/share/dotnet/host"},
				},
			},
		},
		Links: []Link{
			{Source: "[[install-root]]/host/fxr", Target: "[[link-root]]/host/fxr"},
			{Source: "[[install-root]]/sdk/[[version]]", Target: "[[link-root]]/sdk/[[version]]"},
			{Source: "[[install-root]]/shared/Microsoft.NETCore.App", Target: "[[link-root]]/shared/Microsoft.NETCore.App/[[version]]"},
			{Source: "[[install-root]]/shared/Microsoft.AspNetCore.All", Target: "[[link-root]]/shared/Microsoft.AspNetCore.All/[[version]]"},
			{Source: "[[install-root]]/shared/Microsoft.AspNetCore.App", Target: "[[link-root]]/shared/Microsoft.AspNetCore.App/[[version]]"},
		},
	}

	return &m
}
