package shell

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func (t Section) Render() string {
	sb := strings.Builder{}
	gen := t.GeneratorGetter()
	tmpl := gen.GetDefaultTemplate(t.TemplateID)
	err := tmpl.ExecuteTemplate(&sb, t.TemplateID, t.Data)
	if err != nil {
		gen.Fail(err, "failed when rendering template. templateID=%v; data=%v, template=%v\n\n", t.TemplateID, t.Data, *tmpl)
	}
	rendered := sb.String()
	return rendered
}

type Section struct {
	TemplateID      string
	GeneratorGetter func() *ShellHookGenerator
	Data            interface{}
	Title           string
}

func (g *ShellHookGenerator) GetDefaultTemplate(templateID string) *template.Template {
	return g.Config.Templates.Lookup(templateID)
}

func (gen *ShellHookGenerator) NewSection(section string, data interface{}) (out *Section) {
	return NewSection(func() *ShellHookGenerator {
		return gen
	}, section, data)
}

func NewSection(genGetter GeneratorGetter, section string, data interface{}) (out *Section) {
	var Title = strings.ToUpper(section)
	var TemplateID = strings.ToLower(section)

	out = &Section{
		Title:           Title,
		TemplateID:      TemplateID,
		GeneratorGetter: genGetter,
		Data:            data,
	}
	return out
}

// Configurator provides an api to configure ShellHookConfig
type Configurator interface {
	Get() (root *ShellHookConfig)
	ListSections() (sections []string)
	AddOrUpdateSection(name string, section Section) *ShellHookConfig
	AddOrUpdateSections(sections Sections) *ShellHookConfig
	LookupTemplate(templateID string) *template.Template
}

type SDK struct {
	SDK, Path string
}
type SDKList []SDK

type DevCtlSectionConfig struct {
	Prefix string
}
type SDKSectionConfig struct {
	SDKs []SDK
}

type ShellHookConfig struct {
	GeneratorGetter GeneratorGetter
	Templates       *template.Template
	Sections        Sections
}

type GeneratorGetter func() *ShellHookGenerator

func (c *ShellHookGenerator) ListSections() (sections map[string]string) {
	return map[string]string{
		"devctl": "devctl",
		"sdk":    "sdk",
	}
}

// AddOrUpdateSection adds or replaces a section of ShellHookConfig
func (c *ShellHookConfig) AddOrUpdateSection(section *Section) *ShellHookConfig {
	c.Sections = append(c.Sections, *section)
	return c
}

func (g *ShellHookGenerator) AddSection(factory func(gen GeneratorGetter) Section) *ShellHookGenerator {
	section := factory(func() *ShellHookGenerator {
		return g
	})
	g.Config.Sections = append(g.Config.Sections, section)
	return g
}

func (g *ShellHookConfig) AddSection(factory func(gen GeneratorGetter) Section) *ShellHookConfig {
	section := factory(g.GeneratorGetter)
	g.Sections = append(g.Sections, section)
	return g
}

type ShellHookApplyFn func(*ShellHookConfig) *ShellHookConfig
type ShellHookSectionsApplyFn func(*Sections) *Sections

func (c *ShellHookConfig) AddSections(sections Sections) *ShellHookConfig {
	return c.Apply(func(cfg *ShellHookConfig) *ShellHookConfig {
		cfg.Sections = append(cfg.Sections, sections...)
		return cfg
	})
}

func (c *ShellHookConfig) Apply(applyFns ...ShellHookApplyFn) *ShellHookConfig {
	result := c
	for _, fn := range applyFns {
		result = fn(result)
	}
	return result
}

type ShellHookGenerator struct {
	Config *ShellHookConfig
	Out    *bytes.Buffer
	Errors []error
}

func (g *ShellHookGenerator) Write(p []byte) (n int, err error) {
	println("[ShellHookGenerator] Write")
	return g.Out.Write(p)
}

type ShellHook struct {
	ShellScriptString string
	Config            ShellHookConfig
	Buffer            *bytes.Buffer
}

func (s *ShellHook) String() string {
	return s.ShellScriptString
}

func (s *ShellHook) GoString() string {
	return s.ShellScriptString
}

func (g *ShellHookGenerator) Generate() (hook *ShellHook, err error) {
	config := g.Config
	err = config.Templates.ExecuteTemplate(g.Out, "zsh", &config)
	if err != nil {
		fmt.Printf("error: %+v\n\n", err)
		return nil, errors.Wrapf(err, "failed to execute template; template=%s; data=%v\n", "zsh", config)
	}

	sh := g.Out.String()
	println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	println(sh)
	println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	hook = &ShellHook{
		ShellScriptString: sh,
		Buffer:            &bytes.Buffer{},
		Config:            *config,
	}

	written, err := io.Copy(hook.Buffer, g.Out)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to copy shell hook buffer into ShellHook; template=%s; written=%d\n", "zsh", written)
	}

	return hook, nil
}

func (g *ShellHookGenerator) Fail(err error, msg string, args ...interface{}) {
	err = errors.Wrapf(err, msg, args...)
	g.Errors = append(g.Errors, err)
}
