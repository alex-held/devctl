package install

import (
	. "github.com/alex-held/dev-env/config"
	"github.com/alex-held/dev-env/execution"
	. "github.com/alex-held/dev-env/manifest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func setup(t *testing.T, config Config) (*assert.Assertions, afero.Fs) {
	fs := afero.NewMemMapFs()
	_ = config.WriteToFile(fs, "config.json")
	return assert.New(t), fs
}

func TestInstallAddsSDKToConfig(t *testing.T) {
	a, fs := setup(t, Config{Sdks: []SDK{
		{
			Name:    "java",
			Version: "1.8",
			Path:    "java-1.8",
		},
	}})

	executeInstall(fs, []string{"dotnet", "3.1.100"})
	config, _ := ReadConfigFromFile(fs, "config.json")
	a.Len(config.Sdks, 2)
	dotnetSDK := config.Sdks[1]
	a.Equal("dotnet", dotnetSDK.Name)
	a.Equal("3.1.100", dotnetSDK.Version)
	a.Equal("dotnet-3.1.100", dotnetSDK.Path)
}

func TestInstallManifest(t *testing.T) {
	manifest := GetTestManifest()
	sb := strings.Builder{}
	fs := afero.NewMemMapFs()
	executor := execution.NewCommandExecutor(manifest, func(str string) {
		sb.WriteString(str)
	})
	executor.FS = &fs
	executor.Options.DryRun = true

	out, err := executor.Execute()

	println(out)
	assert.NoError(t, err)

	err = Install(*manifest)
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
