package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	. "github.com/alex-held/dev-env/config"
)

func TestName(t *testing.T) {
	spec := Spec{
		Package: SpecPackage{
			Name:        "dotnet",
			Version:     "3.1.202",
			Tags:        []string{"dotnet", "sdk", "core"},
			Repo:        "https://github.com/dotnet/sdk",
			Description: "dotnet sdk",
		},
		Variables: map[string]string{
			"link_root": "/usr/local/share/dotnet",
		},
		Install: struct{ Instructions []string }{
			Instructions: []string{
				"curl https://download.visualstudio.microsoft.com | tar -x -C {{ .InstallLocation }}",
				"ln -s {{ .InstallLocation }}/host/fxr /usr/local/share/dotnet/host/fxr'",
			}},
		UninstallCmds: []string{
			"rm /usr/local/share/dotnet",
			"rm {{ .InstallLocation }}",
		},
	}

	bytes, _ := yaml.Marshal(spec)
	yml := string(bytes)
	println(yml)
}

func TestSpec_GetInstallInstructions(t *testing.T) {
	spec := Spec{
		Package: SpecPackage{
			Name:        "dotnet",
			Version:     "3.1.202",
			Tags:        []string{"dotnet", "sdk", "core"},
			Repo:        "https://github.com/dotnet/sdk",
			Description: "dotnet sdk",
		}, Variables: map[string]string{
			"link_root": "/usr/local/share/dotnet",
		},
		Install: struct{ Instructions []string }{
			Instructions: []string{
				"curl https://download.visualstudio.microsoft.com | tar -x -C {{ .InstallLocation }}",
				"ln -s {{ .InstallLocation }}/host/fxr {{ $link_root }}/host/fxr'",
			}},
		UninstallCmds: nil,
	}
	path := NewTestPathFactory()
	instructions, err := spec.GetInstallInstructions(path)
	require.NoError(t, err)
	require.Len(t, instructions, 2)
	require.Equal(t, "curl https://download.visualstudio.microsoft."+
		"com | tar -x -C /home/sdk/dotnet/3.1.202", instructions[0])
	require.Equal(t, "ln -s /home/sdk/dotnet/3.1."+
		"202/host/fxr /usr/local/share/dotnet/host/fxr'", instructions[1])
}

func NewTestPathFactory() PathFactory {
	home := "/home"
	path := DefaultPathFactory{
		UserHomeOverride: &home,
	}
	return &path
}
