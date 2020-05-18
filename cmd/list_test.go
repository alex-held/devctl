package cmd

import (
	config2 "github.com/alex-held/dev-env/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ListCmdTest struct {
	args     []string
	config   config2.Config
	expected []config2.SDK
}

var (
	sdks = []config2.SDK{
		{
			Name:    "dotnet",
			Version: "3.1.100",
			Path:    "dotnet-3.1.100",
		},
		{
			Name:    "dotnet",
			Version: "2.1.100",
			Path:    "dotnet-2.1.100",
		},
		{
			Name:    "java",
			Version: "1.8",
			Path:    "java-1.8",
		},
		{
			Name:    "go",
			Version: "1.41",
			Path:    "go-1.41",
		},
	}
	dotnet = []config2.SDK{
		{
			Name:    "dotnet",
			Version: "3.1.100",
			Path:    "dotnet-3.1.100",
		},
		{
			Name:    "dotnet",
			Version: "2.1.100",
			Path:    "dotnet-2.1.100",
		}}
)

func TestExecuteList_WithNoArguments_ListAllSDKs(t *testing.T) {

	test := ListCmdTest{
		args: []string{},
		config: config2.Config{
			Sdks: sdks,
		},
		expected: sdks,
	}
	test.run(t)
}

func Test_ExecuteList_WithMultipleSDKArguments_ListMultipleSDKs(t *testing.T) {
	goAndDotnetSdks := append(dotnet, config2.SDK{
		Name:    "go",
		Version: "1.41",
		Path:    "go-1.41",
	})
	test := ListCmdTest{
		args: []string{"dotnet", "go"},
		config: config2.Config{
			Sdks: goAndDotnetSdks,
		},
		expected: goAndDotnetSdks,
	}
	test.run(t)
}

func TestExecuteList_WithSDKArgument_ListMatchingSDKs(t *testing.T) {
	test := ListCmdTest{
		args: []string{"dotnet"},
		config: config2.Config{
			Sdks: sdks,
		},
		expected: dotnet,
	}
	test.run(t)
}

func (test ListCmdTest) run(t *testing.T) {
	result := executeList(test.config, test.args)
	assert.ElementsMatch(t, test.expected, result)
}
