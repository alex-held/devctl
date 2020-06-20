package kernal

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alex-held/dev-env/config"
)

func NewTestPathFactory() config.PathFactory {
	homeOverride := "/home"
	return &config.DefaultPathFactory{
		UserHomeOverride: &homeOverride,
		DevEnvDirectory:  ".devenv",
	}
}

func TestGitInstaller_Install(t *testing.T) {
	mockPathFactory := NewTestPathFactory()
	var output []string
	gitInstaller := NewGitInstaller(mockPathFactory)
	spec := config.Spec{
		Name:          "dotnet",
		Version:       "3.1.202",
		Type:          "git",
		Tags:          []string{"sdk", "dotnet"},
		Repo:          "https://github.com/dotnet/sdk",
		Desc:          "the dotnet sdk",
		InstallCmds:   []string{"echo Installing dotnet sdk"},
		UninstallCmds: []string{"echo Uninstalling dotnet sdk"},
	}

	go gitInstaller.Install(spec)
	for o := range *gitInstaller.Output() {
		println(o)
		output = append(output, o)
	}
	finished, err := gitInstaller.Finished()
	require.True(t, finished)
	require.NoError(t, err)
	require.Equal(t, output[0], "Cloning https://github.com/dotnet/sdk into /home/.devenv/sdk/dotnet/3.1.202")
	require.Equal(t, output[1], "Installing dotnet")
	require.Equal(t, len(output), 2)
}

func TestGitInstaller_Install_HasError_For_Not_GitType_Spec(t *testing.T) {
	mockPathFactory := NewTestPathFactory()
	var output []string
	gitInstaller := NewGitInstaller(mockPathFactory)

	spec := config.Spec{
		Name:          "dotnet",
		Version:       "3.1.202",
		Type:          "tar",
		Tags:          []string{"sdk", "dotnet"},
		Repo:          "https://github.com/dotnet/sdk",
		Desc:          "the dotnet sdk",
		InstallCmds:   []string{"echo Installing dotnet sdk"},
		UninstallCmds: []string{"echo Uninstalling dotnet sdk"},
	}
	go gitInstaller.Install(spec)
	for o := range *gitInstaller.Output() {
		println(o)
		output = append(output, o)
	}
	finished, err := gitInstaller.Finished()
	require.Equal(t, output[0], "Cloning https://github.com/dotnet/sdk into /home/.devenv/sdk/dotnet/3.1.202")
	require.Equal(t, len(output), 1)
	require.True(t, finished)
	require.Error(t, err)
}
