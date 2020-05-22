package install

import (
	"fmt"
	. "github.com/alex-held/dev-env/config"
	"github.com/alex-held/dev-env/manifest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os/exec"
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
	err := Install(*manifest)
	fmt.Errorf(err.Error())
}

func TestCommandExecutor_Execute(t *testing.T) {
	cmd := exec.Command("mkdir", "-p", "/Users/dev/.dev-env/sdk/dotnet/3.1.100")
	err := cmd.Run()

	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}
}

func GetTestManifest() *manifest.Manifest {
	m := manifest.Manifest{
		Version: "3.1.100",
		SDK:     "dotnet",
		Variables: manifest.StringSliceStringMap{
			"url":          "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz",
			"install-root": "[[_sdks]]/[[sdk]]/[[version]]",
			"link-root":    "/usr/local/share/dotnet",
		},
		Install: manifest.Installer{
			Instructions: map[int]manifest.Instruction{
				0: manifest.DevEnvCommand{
					Command: "",
					Args:    nil,
				},
				1: manifest.Pipe{
					Commands: []manifest.DevEnvCommand{
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
		},
		Links: []manifest.Link{
			{Source: "[[install-root]]/host/fxr", Target: "[[link-root]]/host/fxr"},
			{Source: "[[install-root]]/sdk/[[version]]", Target: "[[link-root]]/sdk/[[version]]"},
			{Source: "[[install-root]]/shared/Microsoft.NETCore.App", Target: "[[link-root]]/shared/Microsoft.NETCore.App/[[version]]"},
			{Source: "[[install-root]]/shared/Microsoft.AspNetCore.All", Target: "[[link-root]]/shared/Microsoft.AspNetCore.All/[[version]]"},
			{Source: "[[install-root]]/shared/Microsoft.AspNetCore.App", Target: "[[link-root]]/shared/Microsoft.AspNetCore.App/[[version]]"},
		},
	}

	manifest.PrintYaml(m)
	manifest.PrintJson(m)
	return &m
}
