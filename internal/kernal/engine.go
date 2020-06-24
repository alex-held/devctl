package kernal

import (
	"fmt"
	"os"
	
	pipes "github.com/ebuchman/go-shell-pipes"
	"github.com/rs/zerolog/log"
	
	"github.com/alex-held/dev-env/internal/api"
	"github.com/alex-held/dev-env/internal/spec"
	"github.com/alex-held/dev-env/shared"
)

type EngineCore struct {
	API    api.GithubAPI
	DryRun bool
	path   shared.DefaultPathFactory
}

type Installer interface {
	Runnable
	Output() *chan string
	Install(spec spec.Spec)
	Uninstall(spec spec.Spec)
}

type Runnable interface {
	Started() (finished bool)
	Finished() (finished bool, err error)
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
	case spec.Spec:
		log.Info().Msg("> Installing Application")
		installer := NewInstaller(&engine.path, InstallerOptions{dry: engine.DryRun})
		go installer.Install(e)
		for o := range *installer.Output() {
			println(o)
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
