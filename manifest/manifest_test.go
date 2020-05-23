package manifest

import (
	"fmt"
	"github.com/alex-held/dev-env/config"
	. "github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

var configFilePath = "/Users/dev/go/src/github.com/dev-env/testdata/manifest"
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
		Variables: StringSliceStringMap{
			"url":          "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz",
			"install-root": "[[_sdks]]/[[sdk]]/[[version]]",
			"link-root":    "/usr/local/share/dotnet",
		},
		Install: Installer{},
		Links:   nil,
	}

	variables := manifest.ResolveVariables()
	a.Equal("https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz", variables["[[url]]"])
	a.Equal(config.GetSdks(), variables["[[_sdks]]"])
	a.Equal("dotnet", variables["[[sdk]]"])
	a.Equal("3.1.100", variables["[[version]]"])
	a.Equal("/usr/local/share/dotnet", variables["[[link-root]]"])
	a.Equal(config.GetSdks()+"/"+"dotnet"+"/"+"3.1.100", variables["[[install-root]]"])
}

func TestManifest_ResolveCommands(t *testing.T) {
	a, _ := setup(t)
	installRoot := path.Join(config.GetSdks(), "dotnet", "3.1.100")
	url := "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz"

	manifest := Manifest{
		Version: "3.1.100",
		SDK:     "dotnet",
		Variables: StringSliceStringMap{
			"url":          url,
			"install-root": "[[_sdks]]/[[sdk]]/[[version]]",
			"link-root":    "/usr/local/share/dotnet",
		},
		Instructions: Instructions{
			DevEnvCommand{
				Command: "mkdir",
				Args:    []string{"-p", "[[install-root]]"},
			},
			Pipe{
				Commands: []DevEnvCommand{
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
	installRoot := path.Join(config.GetSdks(), "dotnet", "3.1.100")
	linkRoot := "/usr/local/share/dotnet"

	manifest := Manifest{
		Version: "3.1.100",
		SDK:     "dotnet",
		Variables: StringSliceStringMap{
			"url":          "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz",
			"install-root": "[[_sdks]]/[[sdk]]/[[version]]",
			"link-root":    "/usr/local/share/dotnet",
		},
		Install: Installer{},
		Links:   links,
	}

	resolvedLinks := manifest.ResolveLinks()
	manifest.Links = resolvedLinks

	a.Equal(Link{Source: path.Join(installRoot, "host/fxr"), Target: path.Join(linkRoot, "host/fxr")}, resolvedLinks[0])
	a.Equal(Link{Source: path.Join(installRoot, "sdk", "3.1.100"), Target: path.Join(linkRoot, "sdk", "3.1.100")}, resolvedLinks[1])
	a.Equal(Link{Source: path.Join(installRoot, "shared/Microsoft.NETCore.App"), Target: path.Join(linkRoot, "shared/Microsoft.NETCore.App", "3.1.100")}, resolvedLinks[2])
	a.Equal(Link{Source: path.Join(installRoot, "shared/Microsoft.AspNetCore.All"), Target: path.Join(linkRoot, "shared/Microsoft.AspNetCore.All", "3.1.100")}, resolvedLinks[3])
	a.Equal(Link{Source: path.Join(installRoot, "shared/Microsoft.AspNetCore.App"), Target: path.Join(linkRoot, "shared/Microsoft.AspNetCore.App", "3.1.100")}, resolvedLinks[4])
	println(manifest.Format(Table))
}

func TestReadYaml(t *testing.T) {
	a, _ := setup(t)
	file, _ := ReadFile(NewOsFs(), dotnetYAML)
	yaml := string(file)
	m := &Manifest{}
	manifest, err := readYaml(yaml, m)
	a.NoError(err)
	a.Len(manifest.Variables, 3)

	a.Contains(manifest.Variables, "install-root")
	a.Contains(manifest.Variables, "url")
	a.Contains(manifest.Variables, "link-root")
	a.Equal("https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz", manifest.Variables.ToMap()["url"]) //nolint:lll
	a.Equal("[[_sdks]]/[[sdk]]/[[version]]", manifest.Variables.ToMap()["install-root"])
	a.Equal("/usr/local/share/dotnet", manifest.Variables.ToMap()["link-root"])
	a.Equal(Installer{
		Instructions: map[int]Instruction{
			0: DevEnvCommand{
				Command: "mkdir",
				Args:    []string{"-p", "[[install-root]]"},
			},
			1: Pipe{Commands: []DevEnvCommand{
				{
					Command: "curl",
					Args:    []string{"[[url]]"},
				},
				{
					Command: "tar",
					Args:    []string{"-C", "[[install-root]]", "-x"},
				},
			}},
		},
	}, manifest.Install)
	a.ElementsMatch(links, manifest.Links)
}

func setup(t *testing.T) (*assert.Assertions, Fs) {
	fs := NewMemMapFs()
	return assert.New(t), fs
}
