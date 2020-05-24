package manifest

import (
	"fmt"
	. "github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"testing"
)

var configFilePath = "/Users/dev/.go/src/github.com/alex-held/dev-env/testdata/manifest"
var dotnetYAML = path.Join(configFilePath, "dotnet.yaml")
var _ = path.Join(configFilePath, "dotnet.json")
var _ = StringSliceStringMap{
	"url":          "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz", //nolint:lll,gofmt
	"install-root": "[[_sdks]]/[[sdk]]/[[version]]",
	"link-root":    "/usr/local/share/dotnet",
}
var links = []Link{
	{Source: "[[install-root]]/host/fxr", Target: "[[link-root]]/host/fxr"},
	{Source: "[[install-root]]/sdk/[[version]]", Target: "[[link-root]]/sdk/[[version]]"},
	{Source: "[[install-root]]/shared/Microsoft.NETCore.App", Target: "[[link-root]]/shared/Microsoft.NETCore.App/[[version]]"},
	{Source: "[[install-root]]/shared/Microsoft.AspNetCore.All", Target: "[[link-root]]/shared/Microsoft.AspNetCore.All/[[version]]"},
	{Source: "[[install-root]]/shared/Microsoft.AspNetCore.App", Target: "[[link-root]]/shared/Microsoft.AspNetCore.App/[[version]]"},
}

func TestResolveVariables(t *testing.T) {
	a, _ := setup(t)
	manifest := Manifest{
		Version: "3.1.100",
		SDK:     "dotnet",
		Variables: Variables{
			{Key: "url", Value: "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz"},
			{Key: "install-root", Value: "[[_sdks]]/[[sdk]]/[[version]]"},
			{Key: "link-root", Value: "/usr/local/share/dotnet"},
		},
		Links: nil,
	}

	variables := manifest.ResolveVariables()
	a.Equal("https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz", variables["[[url]]"])
	a.Equal(DefaultPaths.GetSdks(), variables["[[_sdks]]"])
	a.Equal("dotnet", variables["[[sdk]]"])
	a.Equal("3.1.100", variables["[[version]]"])
	a.Equal("/usr/local/share/dotnet", variables["[[link-root]]"])
	a.Equal(DefaultPaths.GetSdks()+"/"+"dotnet"+"/"+"3.1.100", variables["[[install-root]]"])
}

func TestManifest_ResolveCommands(t *testing.T) {
	a, _ := setup(t)
	installRoot := path.Join(DefaultPaths.GetSdks(), "dotnet", "3.1.100")
	url := "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz"

	manifest := Manifest{
		Version: "3.1.100",
		SDK:     "dotnet",
		Variables: Variables{
			{Key: "url", Value: url},
			{Key: "install-root", Value: installRoot},
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
		Links: links,
	}

	resolvedCommands := manifest.ResolveInstructions()
	a.Equal(fmt.Sprintf("mkdir -p %s", installRoot), resolvedCommands[0].Format())
	a.Equal(fmt.Sprintf("curl %s | tar -C %s -x", url, installRoot), resolvedCommands[1].Format())
}

func TestManifest_ResolveLinks(t *testing.T) {
	a, _ := setup(t)
	installRoot := path.Join(DefaultPaths.GetSdks(), "dotnet", "3.1.100")
	linkRoot := "/usr/local/share/dotnet"

	manifest := Manifest{
		Version: "3.1.100",
		SDK:     "dotnet",
		Variables: Variables{
			{Key: "url", Value: "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz"},
			{Key: "install-root", Value: "[[_sdks]]/[[sdk]]/[[version]]"},
			{Key: "link-root", Value: "/usr/local/share/dotnet"},
		},
		Links: links,
	}

	resolvedLinks := manifest.resolveLinks()
	manifest.Links = resolvedLinks

	a.Equal(Link{Source: path.Join(installRoot, "host/fxr"), Target: path.Join(linkRoot, "host/fxr")}, resolvedLinks[0])
	a.Equal(Link{Source: path.Join(installRoot, "sdk", "3.1.100"), Target: path.Join(linkRoot, "sdk", "3.1.100")}, resolvedLinks[1])
	a.Equal(Link{Source: path.Join(installRoot, "shared/Microsoft.NETCore.App"), Target: path.Join(linkRoot, "shared/Microsoft.NETCore.App", "3.1.100")}, resolvedLinks[2])
	a.Equal(Link{Source: path.Join(installRoot, "shared/Microsoft.AspNetCore.All"), Target: path.Join(linkRoot, "shared/Microsoft.AspNetCore.All", "3.1.100")}, resolvedLinks[3])
	a.Equal(Link{Source: path.Join(installRoot, "shared/Microsoft.AspNetCore.App"), Target: path.Join(linkRoot, "shared/Microsoft.AspNetCore.App", "3.1.100")}, resolvedLinks[4])
	println(manifest.Format(Table))
}

func TestWriteYaml(t *testing.T) {
	PrintYaml(manifest)
}

func TestReadYaml(t *testing.T) {
	a, _ := setup(t)
	file, _ := ioutil.ReadFile(dotnetYAML)
	yaml := string(file)
	m := &Manifest{}
	manifest, err := readYaml(yaml, m)
	a.NoError(err)
	a.Len(manifest.Variables, 3)

	// Variables
	a.Equal("[[_sdks]]/[[sdk]]/[[version]]", manifest.MustGetVariable("install-root"))
	a.Equal("dotnet", manifest.MustGetVariable("sdk"))
	a.Equal("3.2.202", manifest.MustGetVariable("version"))
	a.Equal("/usr/local/share/dotnet", manifest.MustGetVariable("link-root"))
	a.Equal("https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz", manifest.MustGetVariable("url"))

	// Instructions
	a.Equal(Step{Command: &DevEnvCommand{
		Command: "mkdir",
		Args:    []string{"-p", "[[install-root]]"},
	}}, manifest.Instructions[0])

	a.Equal(Step{Pipe: []DevEnvCommand{
		{Command: "curl", Args: []string{"[[url]]"}},
		{Command: "tar", Args: []string{"-C", "[[install-root]]", "-x"}},
	}}, manifest.Instructions[1])

	// Links
	a.Equal(Link{Source: "[[install-root]]/host/fxr", Target: "[[link-root]]/host/fxr"}, manifest.Links[0])
	a.Equal(Link{Source: "[[install-root]]/sdk/[[version]]", Target: "[[link-root]]/sdk/[[version]]"}, manifest.Links[1])
	a.Equal(Link{Source: "[[install-root]]/shared/Microsoft.NETCore.App", Target: "[[link-root]]/shared/Microsoft.NETCore.App/[[version]]"}, manifest.Links[2])
	a.Equal(Link{Source: "[[install-root]]/shared/Microsoft.AspNetCore.All", Target: "[[link-root]]/shared/Microsoft.AspNetCore.All/[[version]]"}, manifest.Links[3])
	a.Equal(Link{Source: "[[install-root]]/shared/Microsoft.AspNetCore.App", Target: "[[link-root]]/shared/Microsoft.AspNetCore.App/[[version]]"}, manifest.Links[4])
}

func setup(t *testing.T) (*assert.Assertions, Fs) {
	fs := NewMemMapFs()
	return assert.New(t), fs
}
