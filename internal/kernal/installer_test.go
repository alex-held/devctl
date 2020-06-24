package kernal

import (
	"testing"
	
	"github.com/stretchr/testify/require"
	
	spec2 "github.com/alex-held/dev-env/internal/spec"
	"github.com/alex-held/dev-env/shared"
)

func TestInstaller_Install(t *testing.T) {
	path := shared.NewTestPathFactory()
	opt := InstallerOptions{dry: true}
	var output []string
	spec := spec2.Spec{
		Package: spec2.SpecPackage{
			Name:        "dotnet",
			Version:     "3.1.202",
			Tags:        []string{"dotnet", "sdk", "core"},
			Repo:        "https://github.com/dotnet/sdk",
			Description: "dotnet sdk",
		},
		Install: struct{ Instructions []string }{
			Instructions: []string{
				"curl https://download.visualstudio.microsoft.com | tar -x -C {{ .InstallLocation }}",
				"ln -s {{ .InstallLocation }}/host/fxr /usr/local/share/dotnet/host/fxr'",
			}},
		UninstallCmds: nil,
	}
	installer := NewInstaller(path, opt)
	go installer.Install(spec)
	for o := range *installer.Output() {
		println(o)
		output = append(output, o)
	}
	finished, err := installer.Finished()
	require.True(t, finished)
	require.NoError(t, err)
	require.Equal(t, "Installing dotnet into directory /home/.devenv/sdk/dotnet/3.1.202", output[0])
	require.Equal(t, "[STEP 1/2]curl https://download.visualstudio.microsoft."+
		"com | tar -x -C /home/.devenv/sdk/dotnet/3.1.202", output[1])
	require.Equal(t, "[STEP 2/2]ln -s /home/.devenv/sdk/dotnet/3.1."+
		"202/host/fxr /usr/local/share/dotnet/host/fxr'", output[2])
	require.Equal(t, len(output), 3)
}
