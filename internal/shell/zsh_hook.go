package shell

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

type Section struct {
	*node
	Title      string
	TemplateID string
	Template   *template.Template
	Data       interface{}
}

func (s *Section) GetDefaultTemplate() *template.Template {
	cfg := s.getRootConfig()
	return cfg.LookupTemplate(s.TemplateID)
}

type UninitializedSection func(cfgNode RootCfgNode) *Section

func CreateUninitializedSection(section string, data interface{}) UninitializedSection {
	return func(root RootCfgNode) *Section {
		return newSection(section, data).InitializeWithRootNode(root)
	}
}

func NewSection(section string, data interface{}) (out *Section) {
	out = &Section{
		Title:      strings.ToUpper(section),
		TemplateID: strings.ToLower(section),
		Template:   template.Must(template.ParseGlob("templates/*.tmpl")),
		Data:       data,
	}
	out.node = nil
	return out
}

func (s *Section) InitializeWithRootNode(r RootCfgNode) *Section {
	root := r.(*rootNode)
	s.node = &node{
		parent:   root,
		isRooted: false,
	}
	return s.initialize(root)
}

func (s *Section) initialize(root *rootNode) *Section {
	s.Template = root.config.GetTemplateForSection(*s)
	s.TemplateID = strings.ToLower(s.TemplateID)
	s.Title = strings.ToUpper(s.Title)
	return s
}

// Configurator provides an api to configure ShellHookConfig
type Configurator interface {
	Get() (root *HookConfig)
	ListSections() (sections []string)
	AddOrUpdateSection(name string, section Section) *HookConfig
	AddOrUpdateSections(sections Sections) *HookConfig
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

type HookConfig struct {
	root      *rootNode
	Templates *template.Template
	Sections  Sections
}

func (c *HookConfig) GetTemplateForSection(section Section) *template.Template {
	return c.Templates.Lookup(section.TemplateID)
}

func (c *HookConfig) ApplyWithSelector(inSelector InSelectorApplyFn, mapFn MapFn, updateSelector UpdateConfigSelector) *HookConfig { //nolint:lll
	return c.Apply(func(c *HookConfig) *HookConfig {
		selectInput := inSelector(c)
		mapped := mapFn(selectInput)
		updateSelector(c, mapped)
		return c
	})
}

func (c *HookConfig) LookupTemplate(templateID string) *template.Template {
	return c.Templates.Lookup(templateID)
}

func (*HookConfig) ListSections() (sections map[string]string) {
	return map[string]string{
		"devctl": "devctl",
		"sdk":    "sdk",
	}
}

func (c *HookConfig) GetSection(title string) (section Section, ok bool) {
	v, o := c.Sections[title]
	return v, o
}

func (c *HookConfig) Get() (root *HookConfig) { return c }

// AddOrUpdateSection adds or replaces a section of ShellHookConfig
func (c *HookConfig) AddOrUpdateSection(section *Section) *HookConfig {
	return c.AddOrUpdateSectionForKey(section.Title, *section)
}

// AddOrUpdateSectionForKey adds or replaces a section of ShellHookConfig
func (c *HookConfig) AddOrUpdateSectionForKey(name string, section Section) *HookConfig {
	return c.Apply(func(c *HookConfig) *HookConfig {
		initializedSection := section.InitializeWithRootNode(c.root)
		_ = fmt.Sprintf("intiialized section 1=\n%+v\n", *initializedSection)
		c.Sections[name] = *initializedSection
		return c

		/*	configuredSection := initializedSection.Apply(func(in *Section) (out *Section) {
			in.Title = strings.ToUpper(name)
			in.TemplateId = strings.ToLower(name)
			in.Data = section.Data
			in.node = &node{
				parent:   c.root.node,
				isRooted: false,
			}
			in.Template = c.Templates.lookupSymbol(in.TemplateId)
			out = in
			return out
		})*/
	})
}

type RootCfgNode interface {
	CfgNode
}

type CfgNode interface {
	IsRooted() bool
	Parent() CfgNode
	RootNode() RootCfgNode
}

type rootNode struct {
	*node
	config *HookConfig
}

func (r *rootNode) NewChildNode() *node {
	child := &node{
		parent:   RootCfgNode(r),
		isRooted: false,
	}
	return child
}

func (r *rootNode) Parent() CfgNode       { return r }
func (r *rootNode) RootNode() RootCfgNode { return r }
func (r rootNode) IsRooted() bool         { return true }

type node struct {
	parent   CfgNode
	isRooted bool
}

func (c *node) getRootConfig() (root *HookConfig) {
	rootNode := c.RootNode().(*rootNode)
	return rootNode.config
}

func (c *node) Map(selector func(*HookConfig) (interface{}, error)) (out interface{}, err error) {
	cfg := c.getRootConfig()
	out, err = selector(cfg)
	return out, err
}

func (c *node) RootNode() (root RootCfgNode) {
	parent := c.Parent()

	for !parent.IsRooted() {
		parent = parent.Parent()
	}

	rNode := parent.(*rootNode)
	return rNode.RootNode()
}

func (c *node) Parent() CfgNode { return c.parent }
func (c *node) IsRooted() bool  { return c.isRooted }

type HookApplyFn func(*HookConfig) *HookConfig

type InSelectorApplyFn func(cfg *HookConfig) (applyOn interface{})
type MapFn func(in interface{}) (mappedOut interface{})
type UpdateConfigSelector func(cfg *HookConfig, mappedOut interface{}) (updatedConfig *HookConfig)
type HookApplySelectorFn func(inSelector InSelectorApplyFn, mapFn MapFn, updateSelector UpdateConfigSelector)

type HookApplySelectorFn2 func(
	inSelector func(cfg *HookConfig) (applyOn interface{}),
	mapFn func(in interface{}) (mappedOut interface{}),
	updateConfigSelector func(cfg *HookConfig, mappedOut interface{}) (updatedConfig *HookConfig))

type HookSubApplyFn func(config *HookConfig, applySelectorFn ...HookApplySelectorFn)
type HookSectionsApplyFn func(*Sections) *Sections

func (c *HookConfig) AddOrUpdateSections(sections Sections) *HookConfig {
	return c.Apply(func(c *HookConfig) *HookConfig {
		for k, v := range sections {
			c = c.AddOrUpdateSectionForKey(k, v)
		}
		return c
	})
}

func (c *HookConfig) Apply(applyFns ...HookApplyFn) *HookConfig {
	result := c
	for _, fn := range applyFns {
		result = fn(result)
	}
	return result
}

func (c *HookConfig) Root() (root *HookConfig) { return c }

func (c *node) NextParent() (node CfgNode, ok bool) {
	if c.parent == nil {
		return c, true
	}
	return c.parent, false
}

type HookGenerator struct {
	Out       *bytes.Buffer
	Templates *template.Template
}

func (g *HookGenerator) Write(p []byte) (n int, err error) {
	println("[ShellHookGenerator] Write")
	return g.Out.Write(p)
}

type Hook struct {
	string
	ShellScriptString string
	Config            HookConfig
	Buffer            *bytes.Buffer
}

func (s *Hook) String() string {
	return s.ShellScriptString
}

func (s *Hook) GoString() string {
	return s.ShellScriptString
}

func (g *HookGenerator) Generate(config HookConfig) (hook *Hook, err error) {
	err = g.Templates.ExecuteTemplate(g.Out, "zsh", &config)

	if err != nil {
		fmt.Printf("error: %+v\n\n", err)
		return nil, errors.Wrapf(err, "failed to execute template; template=%s; data=%v\n", "zsh", config)
	}

	sh := g.Out.String()
	fmt.Println("Output: " + sh)
	hook = &Hook{
		string:            sh,
		ShellScriptString: sh,
		Buffer:            &bytes.Buffer{},
		Config:            config,
	}

	written, err := io.Copy(hook.Buffer, g.Out)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to copy shell hook buffer into ShellHook; template=%s; written=%d\n", "zsh", written)
	}

	return hook, nil
}
