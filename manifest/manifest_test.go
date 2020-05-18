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
	variables      = StringMap{
		"url":          "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz",
		"install-root": "[[dev-env-sdks]]/[[sdk]]/[[version]]",
		"link-root":    "/usr/local/share/dotnet",
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
	// vars := manifest.Variables.ToMap()

	a.Subset(manifest.Variables, variables)
	a.Contains(manifest.Variables, "install-root")
	a.Contains(manifest.Variables, "url")
	a.Contains(manifest.Variables, "link-root")
	//    a.ElementsMatch(manifest.Variables.ToMap(), variables)
}

func setup(t *testing.T) (*assert.Assertions, afero.Fs) {
	fs := afero.NewMemMapFs()
	return assert.New(t), fs
}
