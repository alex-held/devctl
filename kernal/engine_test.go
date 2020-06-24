package kernal_test

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"testing"

	"github.com/alex-held/dev-env/api"
	"github.com/alex-held/dev-env/config"
	"github.com/alex-held/dev-env/kernal"
	"github.com/alex-held/dev-env/shared"
)

func init() {
	shared.BootstrapLogger(zerolog.TraceLevel)
}

func TestPrettyPrint(t *testing.T) {
	spec := NewDotnetSpec()
	engine := NewTestEngine(true)
	err := engine.Execute(spec)

	if err != nil {
		println(err.Error())
		assert.NoError(t, err)
	}
}

func NewDotnetSpec() config.Spec {
	spec := config.Spec{
		Package: config.SpecPackage{
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
	return spec
}
func TestExecuteCommands(t *testing.T) {
	commands := []string{
		"curl https://api.github.com/users/alex-held/followers | grep wonderbird",
		"ls -a",
	}
	engine := NewTestEngine(false)
	err := engine.Execute(commands)

	if err != nil {
		println(err.Error())
		assert.NoError(t, err)
	}
}

func NewTestEngine(dry bool) kernal.Engine {
	engine := kernal.EngineCore{
		API:    api.NewGithubAPI(nil),
		DryRun: dry,
	}
	return &engine
}
