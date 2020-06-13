package kernal

import (
	"fmt"
	"os"

	pipes "github.com/ebuchman/go-shell-pipes"
	"github.com/rs/zerolog/log"

	"github.com/alex-held/dev-env/api"
	"github.com/alex-held/dev-env/meta"
)

type EngineCore struct {
	API    api.GithubAPI
	DryRun bool
}

type Engine interface {
	Execute(executable interface{}) error
}

func (engine *EngineCore) Execute(executable interface{}) error {
	switch e := executable.(type) {
	case []string:
		log.Info().Msg("> Executing Commands")
		for i, cmd := range e {
			log.Trace().
				Str("Command", cmd).
				Int("Step", i).
				Msg("")

			if !engine.DryRun {
				out, err := pipes.RunString(cmd)
				if err != nil {
					log.Err(err)
					return err
				}
				_, _ = os.Stdout.WriteString(out + "\n")
			}
		}
		return nil
	case meta.Meta:
		log.Info().Msg("> Installing Application")
		for i, cmd := range e.Install {
			log.Trace().
				Str("Command", cmd).
				Int("Step", i).
				Msg("")
		}

		log.Info().Msg("> Linking Application")
		for i, ln := range e.Link {
			log.Trace().Int("Step", i).Msg(ln)
		}
		return nil
	default:
		return fmt.Errorf("Executable does not have valid type. ")
	}
}

func NewEngine() Engine {
	engine := EngineCore{
		API: api.NewGithubAPI(nil),
	}
	return &engine
}
