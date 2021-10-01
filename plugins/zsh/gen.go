package zsh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/a8m/envsubst/parse"
	"golang.org/x/sync/syncmap"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
	"github.com/alex-held/devctl-kit/pkg/plugins"
)

type Generator interface {
	Completions(w io.Writer) (err error)
	Exports(w io.Writer) (err error)
}

type generator struct {
	*GeneratorConfig
}

var ErrGeneratingTemplate = errors.New("failed to execute template")

type Option func(*GeneratorConfig) *GeneratorConfig

type GeneratorConfig struct {
	HttpClient      *http.Client
	Templates       *template.Template
	TemplateConfigs map[string]interface{}
	Context         *plugins.Context
}

type ResolveContext struct {
	env syncmap.Map

	resolvedC  chan resolvable
	cancelC    <-chan struct{}
	unresolved map[string]resolvable
	context    context.Context
}

// func (c *ResolveContext) ResolveAll() {
// 	for _, unresolved := range c.unresolved {
// 		fullyResolved, resolved, unresolved := unresolved.resolve(c)
// 		if fullyResolved {
// 			c.env.Store(resolved.ID, resolved)
// 		}
// 		for _, unresolved := range unresolved {
// 			c.resolveRecursive(unresolved)
// 		}
// 	}
// }

func (c *ResolveContext) Get(key string) (val string, ok bool) {
	value, ok := c.env.Load(key)
	if v, ok := value.(*string); ok {
		return *v, true
	}
	if v, ok := value.(string); ok {
		return v, true
	}
	return "", false
}

func (c ResolveContext) ProjectEnv() (env parse.Env) {
	c.env.Range(func(key, value interface{}) bool {
		env = append(env, fmt.Sprintf("%v=%v", key, value))
		return true
	})
	return env
}

func (c ResolveContext) Env() (env map[string]string) {
	env = map[string]string{}
	c.env.Range(func(key, value interface{}) bool {
		env[fmt.Sprintf("%v", key)] = fmt.Sprintf("%s", value)
		return true
	})
	return env
}

func (c *ResolveContext) Add(resolvables ...resolvable) {
	for _, r := range resolvables {
		if r.IsResolved() {
			c.env.Store(r.ID, r.Value)
			continue
		}
		c.unresolved[r.ID] = r
	}
}

func (c *ResolveContext) getResolved(id string) (resolved resolvable, ok bool) {
	c.env.Range(func(key, value interface{}) bool {
		if r, ok := value.(*resolvable); ok {
			resolved = *r
			return false
		}
		return true
	})

	return resolved, ok
}

func (c *ResolveContext) resolveRecursive(unresolved string) {
	_, ok := c.getResolved(unresolved)
	if ok {
		return
	}
	u := c.unresolved[unresolved]
	needs := u.Needs()
	for _, need := range needs {
		c.resolveRecursive(need)
	}
}

type Resolvable interface {
	Name() string
	IsResolved() bool
	Unresolved() (unresolved []string)
	Resolve(ctx *ResolveContext) (ok bool)
}

type Export struct {
	Key   string
	Value Resolvable
}

func (g *generator) GenerateExports(w io.Writer, config *Config, ctx *plugins.Context) (err error) {
	type tmpl struct {
		Vars    map[string]string
		Exports map[string]string
	}

	// context.WithTimeout(ctx.Context, time.Second)
	resolverCtx := NewContext(ctx.Context)

	for k, v := range config.Vars {
		r := resolvable{
			ID:    k,
			Value: &v,
		}
		if !r.Resolve(resolverCtx) {
			return errors.New(fmt.Sprintf("unable to resolve '%s'=%s\n", k, v))
		}
	}

	data := tmpl{
		Vars:    config.Vars,
		Exports: config.Exports,
	}

	err = g.Templates.ExecuteTemplate(w, "exports", data)
	return err
}

func NewGenerator(opts ...Option) Generator {
	g := &generator{
		&GeneratorConfig{
			HttpClient:      http.DefaultClient,
			Templates:       templates,
			TemplateConfigs: map[string]interface{}{},
			Context: &plugins.Context{
				Out:     os.Stdout,
				Pather:  devctlpath.DefaultPather(),
				Context: context.Background(),
			},
		},
	}

	for _, opt := range opts {
		opt(g.GeneratorConfig)
	}

	return g
}

func (g *generator) Completions(w io.Writer) (err error) {
	const name = "completions"
	data := g.TemplateConfigs[name]
	return g.Templates.ExecuteTemplate(w, name, data)
}

func (g *generator) Exports(w io.Writer) (err error) {
	const name = "exports"
	data := g.TemplateConfigs[name]
	return g.Templates.ExecuteTemplate(w, name, data)
}
