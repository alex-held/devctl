package manifest

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

var (
	configFilePath = "/Users/dev/go/src/github.com/dev-env/testdata/manifest"
	dotnetYAML     = path.Join(configFilePath, "dotnet.yaml")
	dotnetJSON     = path.Join(configFilePath, "dotnet.json")
	variables      = StringSliceStringMap{
		"url":          "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz",
		"install-root": "[[_sdks]]/[[sdk]]/[[version]]",
		"link-root":    "/usr/local/share/dotnet",
	}
	links = []Link{
		{Source: "[[install-root]]/host/fxr", Target: "[[link-root]]/host/fxr"},
		{Source: "[[install-root]]/sdk/[[version]]", Target: "[[link-root]]/sdk/[[version]]"},
		{Source: "[[install-root]]/shared/Microsoft.NETCore.App", Target: "[[link-root]]/shared/Microsoft.NETCore.App/[[version]]"},
		{Source: "[[install-root]]/shared/Microsoft.AspNetCore.All", Target: "[[link-root]]/shared/Microsoft.AspNetCore.All/[[version]]"},
		{Source: "[[install-root]]/shared/Microsoft.AspNetCore.App", Target: "[[link-root]]/shared/Microsoft.AspNetCore.App/[[version]]"},
	}
)

func TestReadYaml(t *testing.T) {
	a, _ := setup(t)
	file, _ := afero.ReadFile(afero.NewOsFs(), dotnetYAML)
	yaml := string(file)
	m := &Manifest{}
	manifest, err := readYaml(yaml, m)
	a.NoError(err)
	a.Len(manifest.Variables, 3)

	a.Contains(manifest.Variables, "install-root")
	a.Contains(manifest.Variables, "url")
	a.Contains(manifest.Variables, "link-root")
	a.Equal("https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz", manifest.Variables.ToMap()["url"])
	a.Equal("[[_sdks]]/[[sdk]]/[[version]]", manifest.Variables.ToMap()["install-root"])
	a.Equal("/usr/local/share/dotnet", manifest.Variables.ToMap()["link-root"])
	a.ElementsMatch(links, manifest.Links)

}

func setup(t *testing.T) (*assert.Assertions, afero.Fs) {
	fs := afero.NewMemMapFs()
	return assert.New(t), fs
}
