package zsh

import (
	"io"
	"net/http"
	"text/template"
)

type Generator interface {
	Completions(w io.Writer) (err error)
	Exports(w io.Writer) (err error)
}

type generator struct {
	*GeneratorConfig
}

type Option func(*GeneratorConfig) *GeneratorConfig

type GeneratorConfig struct {
	HttpClient      *http.Client
	Templates       *template.Template
	TemplateConfigs map[string]interface{}
}

func NewGenerator(opts ...Option) Generator {
	g := &generator{
		&GeneratorConfig{
			HttpClient:      http.DefaultClient,
			Templates:       templates,
			TemplateConfigs: map[string]interface{}{},
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
