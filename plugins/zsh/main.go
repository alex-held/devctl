package zsh

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/alex-held/devctl-kit/pkg/plugins"
)

func CreateConfig() *Config {
	return &Config{
		Context: &plugins.Context{},
		// Vars:        map[string]string{},
		// Exports:     map[string]string{},
		// Aliases:     map[string]string{},
		// Completions: CompletionsSpec{},
	}
}

var ErrWrongArgumentsProvided = errors.New("number of arguments is invalid")

func Exec(cfg *Config, args []string) (err error) {
	if len(args) == 0 {
		return ErrWrongArgumentsProvided
	}

	args = args[0:]
	z := NewZSH(cfg)

	return z.handlers[args[0]](args[1:])
}

type commandHandler func(args []string) (err error)

type ZSH struct {
	Config    *Config
	Generator Generator
	handlers  map[string]commandHandler
}

func NewZSH(cfg *Config) *ZSH {
	z := &ZSH{
		Config: cfg,
		Generator: NewGenerator(func(gc *GeneratorConfig) *GeneratorConfig {
			gc.HttpClient = http.DefaultClient
			gc.Templates = templates
			gc.TemplateConfigs = cfg.TemplateConfigs()
			// fmt.Printf("TemplateConfigs=%v\n", gc.TemplateConfigs)
			return gc
		}),
		handlers: map[string]commandHandler{},
	}
	z.initHandlers()
	return z
}

func (z *ZSH) initHandlers() {
	z.handlers = map[string]commandHandler{
		"init": func(args []string) (err error) {
			fmt.Printf("init called. args=%v\n", args)
			return nil
		},
		"gen": func(args []string) (err error) {
			usage := fmt.Sprint(`
			USAGE
				devctl zsh gen <TYPE>

			EXAMPLES
				devctl zsh gen completions
				devctl zsh gen exports
			`)

			if len(args) == 0 {
				print(usage)
				return errors.New("required argument not provided")
			}
			if len(args) > 2 {
				print(usage)
				return errors.New("too many arguments provided")
			}

			var filepath string
			logGen := func() {
				fmt.Printf("generation %s at '%s'\n", args[0], filepath)
			}

			g := z.Generator
			var file *os.File

			switch args[0] {
			case "completions":
				filepath = z.Config.Pather.Config("zsh", "init.d", "05-completions.zsh")
				logGen()
				if file, err = os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
					return err
				}

				defer file.Close()
				return g.Completions(file)
			case "exports":
				filepath = z.Config.Pather.Config("zsh", "init.d", "03-exports.zsh")
				logGen()
				if file, err = os.Open(filepath); err != nil {
					return err
				}

				defer file.Close()
				return g.Exports(file)
			default:
				print(usage)
				return nil
			}
		},
	}
}
