package manifest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

type FormatTest struct {
	t        *testing.T
	Actual   string
	Expected string
}

var manifest = Manifest{
	Version: "3.2.202",
	SDK:     "dotnet",
	Variables: Variables{
		{Key: "url", Value: "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz"},
		{Key: "install-root", Value: "[[_sdks]]/[[sdk]]/[[version]]"},
		{Key: "link-root", Value: "/usr/local/share/dotnet"},
	},
	Instructions: Instructions{
		Step{
			Command: &DevEnvCommand{Command: "mkdir", Args: []string{"-p", "[[install-root]]"}},
		},
		Step{
			Pipe: []DevEnvCommand{
				{
					Command: "curl",
					Args:    []string{"[[url]]"},
				},
				{
					Command: "tar",
					Args:    []string{"-C", "[[install-root]]", "-x"},
				},
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

func TestFormatAsTree(t *testing.T) {

	test := FormatTest{
		t:      t,
		Actual: manifest.Format(Tree),
		Expected: `dotnet-3.2.202
└── variables
│   ├── {Key:[[_home]] Value:/Users/dev/.dev-env}
│   ├── {Key:[[_manifests]] Value:/Users/dev/.dev-env/manifests}
│   ├── {Key:[[_sdks]] Value:/Users/dev/.dev-env/sdk}
│   ├── {Key:[[home]] Value:/Users/dev}
│   ├── {Key:[[install-root]] Value:/Users/dev/.dev-env/sdk/dotnet/3.2.202}
│   ├── {Key:[[link-root]] Value:/usr/local/share/dotnet}
│   ├── {Key:[[sdk]] Value:dotnet}
│   ├── {Key:[[url]] Value:https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz}
│   ├── {Key:[[version]] Value:3.2.202}
└── links
│   ├── {Source:/Users/dev/.dev-env/sdk/dotnet/3.2.202/host/fxr Target:/usr/local/share/dotnet/host/fxr}
│   ├── {Source:/Users/dev/.dev-env/sdk/dotnet/3.2.202/sdk/3.2.202 Target:/usr/local/share/dotnet/sdk/3.2.202}
│   ├── {Source:/Users/dev/.dev-env/sdk/dotnet/3.2.202/shared/Microsoft.NETCore.App Target:/usr/local/share/dotnet/shared/Microsoft.NETCore.App/3.2.202}
│   ├── {Source:/Users/dev/.dev-env/sdk/dotnet/3.2.202/shared/Microsoft.AspNetCore.All Target:/usr/local/share/dotnet/shared/Microsoft.AspNetCore.All/3.2.202}
│   ├── {Source:/Users/dev/.dev-env/sdk/dotnet/3.2.202/shared/Microsoft.AspNetCore.App Target:/usr/local/share/dotnet/shared/Microsoft.AspNetCore.App/3.2.202}
└── instructions
    └── 0
    │   ├── mkdir -p /Users/dev/.dev-env/sdk/dotnet/3.2.202
    └── 1
        └── curl https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz
        └── tar -C /Users/dev/.dev-env/sdk/dotnet/3.2.202 -x
`,
	}

	test.Run()

	_ = ioutil.WriteFile("/Users/dev/temp/table-manifest-expected", []byte(test.Expected), 0644)
	_ = ioutil.WriteFile("/Users/dev/temp/table-manifest-actual", []byte(test.Actual), 0644)
}

func TestFormatWithFormatTypeTable(t *testing.T) {

	test := FormatTest{
		t:      t,
		Actual: manifest.Format(Table),
		Expected: `
Properties
  Property | Value
----------- ----------
  Version  | 3.2.202
  SDK      | dotnet

Variables
  Variables        | Value
------------------- ------------------------------------------------------------------------------------------------------------------------------------------------------------------
  [[_home]]        | /Users/dev/.dev-env
  [[_manifests]]   | /Users/dev/.dev-env/manifests
  [[_sdks]]        | /Users/dev/.dev-env/sdk
  [[home]]         | /Users/dev
  [[install-root]] | /Users/dev/.dev-env/sdk/dotnet/3.2.202
  [[link-root]]    | /usr/local/share/dotnet
  [[sdk]]          | dotnet
  [[url]]          | https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz
  [[version]]      | 3.2.202

Links
  Source                                                                 | Target
------------------------------------------------------------------------- ------------------------------------------------------------------
  /Users/dev/.dev-env/sdk/dotnet/3.2.202/host/fxr                        | /usr/local/share/dotnet/host/fxr
  /Users/dev/.dev-env/sdk/dotnet/3.2.202/sdk/3.2.202                     | /usr/local/share/dotnet/sdk/3.2.202
  /Users/dev/.dev-env/sdk/dotnet/3.2.202/shared/Microsoft.NETCore.App    | /usr/local/share/dotnet/shared/Microsoft.NETCore.App/3.2.202
  /Users/dev/.dev-env/sdk/dotnet/3.2.202/shared/Microsoft.AspNetCore.All | /usr/local/share/dotnet/shared/Microsoft.AspNetCore.All/3.2.202
  /Users/dev/.dev-env/sdk/dotnet/3.2.202/shared/Microsoft.AspNetCore.App | /usr/local/share/dotnet/shared/Microsoft.AspNetCore.App/3.2.202

Instructions
  Order | Command
-------- -----------------------------------------------------------------------------------------------------------------------------------------------------------------------
      0 | mkdir -p /Users/dev/.dev-env/sdk/dotnet/3.2.202
      1 | curl https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz
        | tar -C /Users/dev/.dev-env/sdk/dotnet/3.2.202 -x
`,
	}

	test.Run()

	_ = ioutil.WriteFile("/Users/dev/temp/table-manifest-expected", []byte(test.Expected), 0644)
	_ = ioutil.WriteFile("/Users/dev/temp/table-manifest-actual", []byte(test.Actual), 0644)
}

func (test *FormatTest) Run() {
	actual := test.Actual
	expected := test.Expected

	fmt.Printf("<EXPECTED> %d\n", len(expected))
	fmt.Println(expected)
	fmt.Println()

	fmt.Printf("<ACTUAL> %d\n", len(actual))
	fmt.Println(actual)
	fmt.Println()
	assert.Equal(test.t, expected, actual)
}
