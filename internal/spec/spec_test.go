package spec_test

import (
	"io/ioutil"
	"testing"
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	
	spec2 "github.com/alex-held/dev-env/internal/spec"
	"github.com/alex-held/dev-env/shared"
)

func init() {
	shared.BootstrapLogger(zerolog.TraceLevel)
}

func TestName(t *testing.T) {
	spec := spec2.Spec{
		Package: spec2.SpecPackage{
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
	spec := spec2.Spec{
		Package: spec2.SpecPackage{
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
	log.Trace().
		Interface("Instructions", instructions).
		Msg("Replaced template variables in Spec-Instructions")
	require.NoError(t, err)
	require.Len(t, instructions, 2)
	require.Equal(t, "curl https://download.visualstudio.microsoft."+
		"com | tar -x -C /home/sdk/dotnet/3.1.202", instructions[0])
	require.Equal(t, "ln -s /home/sdk/dotnet/3.1."+
		"202/host/fxr /usr/local/share/dotnet/host/fxr'", instructions[1])
}

func TestSpec_Serialze(t *testing.T) {
	expectedBytes, err := ioutil.ReadFile("../testdata/spec/dotnet-3.1.202.yaml")
	expected := string(expectedBytes)
	require.NoError(t, err)
	spec := spec2.Spec{
		Package: spec2.SpecPackage{
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
	bytes, err := yaml.Marshal(spec)
	serialized := string(bytes)
	println(serialized)
	require.NoError(t, err)
	require.Equal(t, expected, serialized)
}

func NewTestPathFactory() shared.PathFactory {
	home := "/home"
	path := shared.DefaultPathFactory{
		UserHomeOverride: &home,
	}
	return &path
}
