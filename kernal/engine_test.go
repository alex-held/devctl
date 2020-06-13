package kernal_test

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/alex-held/dev-env/kernal"
	"github.com/alex-held/dev-env/registry"
	meta2 "github.com/alex-held/dev-env/testdata/meta"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: zerolog.TimeFormatUnix,
	})

	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	lvl := zerolog.GlobalLevel()
	println(lvl.String())
}

func TestPrettyPrint(t *testing.T) {
	meta := meta2.NewDotnetMeta()
	engine := NewTestEngine(true)
	err := engine.Execute(meta)

	if err != nil {
		println(err.Error())
		assert.NoError(t, err)
	}
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
		API:    registry.NewRegistryAPI(),
		DryRun: dry,
	}
	return &engine
}
