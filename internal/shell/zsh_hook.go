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

type UninitializedSection func(cfgNode ShellRootCfgNode) *Section

func CreateUninitializedSection(section string, data interface{}) UninitializedSection {
	return func(root ShellRootCfgNode) *Section {
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
	// todo: this node does not get initialized with a parent. maybe this cases issues at some point
	out.node = nil
	return out
}

func (s *Section) InitializeWithRootNode(r ShellRootCfgNode) *Section {
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
	root      *rootNode
	Templates *template.Template
	Sections  Sections
}

func (c *ShellHookConfig) GetTemplateForSection(section Section) *template.Template {
	return c.Templates.Lookup(section.TemplateID)
}

func (c *ShellHookConfig) ApplyWithSelector(inSelector InSelectorApplyFn, mapFn MapFn, updateSelector UpdateConfigSelector) *ShellHookConfig { //nolint:lll
	return c.Apply(func(c *ShellHookConfig) *ShellHookConfig {
		selectInput := inSelector(c)
		mapped := mapFn(selectInput)
		updateSelector(c, mapped)
		return c
	})
}

func (c *ShellHookConfig) LookupTemplate(templateID string) *template.Template {
	return c.Templates.Lookup(templateID)
}

func (*ShellHookConfig) ListSections() (sections map[string]string) {
	return map[string]string{
		"devctl": "devctl",
		"sdk":    "sdk",
	}
}

func (c *ShellHookConfig) GetSection(title string) (section Section, ok bool) {
	v, o := c.Sections[title]
	return v, o
}

func (c *ShellHookConfig) Get() (root *ShellHookConfig) { return c }

// AddOrUpdateSection adds or replaces a section of ShellHookConfig
func (c *ShellHookConfig) AddOrUpdateSection(section *Section) *ShellHookConfig {
	return c.AddOrUpdateSectionForKey(section.Title, *section)
}

// AddOrUpdateSectionForKey adds or replaces a section of ShellHookConfig
func (c *ShellHookConfig) AddOrUpdateSectionForKey(name string, section Section) *ShellHookConfig {
	return c.Apply(func(c *ShellHookConfig) *ShellHookConfig {
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
			in.Template = c.Templates.Lookup(in.TemplateId)
			out = in
			return out
		})*/
	})
}

type ShellRootCfgNode interface {
	ShellCfgNode
}

type ShellCfgNode interface {
	IsRooted() bool
	Parent() ShellCfgNode
	RootNode() ShellRootCfgNode
}

type rootNode struct {
	*node
	config *ShellHookConfig
}

func (r *rootNode) NewChildNode() *node {
	child := &node{
		parent:   ShellRootCfgNode(r),
		isRooted: false,
	}
	return child
}

func (r *rootNode) Parent() ShellCfgNode       { return r }
func (r *rootNode) RootNode() ShellRootCfgNode { return r }
func (r rootNode) IsRooted() bool              { return true }

type node struct {
	parent   ShellCfgNode
	isRooted bool
}

func (c *node) getRootConfig() (root *ShellHookConfig) {
	rootNode := c.RootNode().(*rootNode)
	return rootNode.config
}

func (c *node) Map(selector func(*ShellHookConfig) (interface{}, error)) (out interface{}, err error) {
	cfg := c.getRootConfig()
	out, err = selector(cfg)
	return out, err
}

func (c *node) RootNode() (root ShellRootCfgNode) {
	parent := c.Parent()

	for !parent.IsRooted() {
		parent = parent.Parent()
	}

	rNode := parent.(*rootNode)
	return rNode.RootNode()
}

func (c *node) Parent() ShellCfgNode { return c.parent }
func (c *node) IsRooted() bool       { return c.isRooted }

type ShellHookApplyFn func(*ShellHookConfig) *ShellHookConfig

type InSelectorApplyFn func(cfg *ShellHookConfig) (applyOn interface{})
type MapFn func(in interface{}) (mappedOut interface{})
type UpdateConfigSelector func(cfg *ShellHookConfig, mappedOut interface{}) (updatedConfig *ShellHookConfig)
type ShellHookApplySelectorFn func(inSelector InSelectorApplyFn, mapFn MapFn, updateSelector UpdateConfigSelector)

type ShellHookApplySelectorFn2 func(
	inSelector func(cfg *ShellHookConfig) (applyOn interface{}),
	mapFn func(in interface{}) (mappedOut interface{}),
	updateConfigSelector func(cfg *ShellHookConfig, mappedOut interface{}) (updatedConfig *ShellHookConfig))

type ShellHookSubApplyFn func(config *ShellHookConfig, applySelectorFn ...ShellHookApplySelectorFn)
type ShellHookSectionsApplyFn func(*Sections) *Sections

func (c *ShellHookConfig) AddOrUpdateSections(sections Sections) *ShellHookConfig {
	return c.Apply(func(c *ShellHookConfig) *ShellHookConfig {
		for k, v := range sections {
			c = c.AddOrUpdateSectionForKey(k, v)
		}
		return c
	})
}

func (c *ShellHookConfig) Apply(applyFns ...ShellHookApplyFn) *ShellHookConfig {
	result := c
	for _, fn := range applyFns {
		result = fn(result)
	}
	return result
}

func (c *ShellHookConfig) Root() (root *ShellHookConfig) { return c }

func (c *node) NextParent() (node ShellCfgNode, ok bool) {
	if c.parent == nil {
		return c, true
	}
	return c.parent, false
}

type ShellHookGenerator struct {
	Out       *bytes.Buffer
	Templates *template.Template
}

func (g *ShellHookGenerator) Write(p []byte) (n int, err error) {
	println("[ShellHookGenerator] Write")
	return g.Out.Write(p)
}

type ShellHook struct {
	string
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

func (g *ShellHookGenerator) Generate(config ShellHookConfig) (hook *ShellHook, err error) {
	err = g.Templates.ExecuteTemplate(g.Out, "zsh", &config)

	if err != nil {
		fmt.Printf("error: %+v\n\n", err)
		return nil, errors.Wrapf(err, "failed to execute template; template=%s; data=%v\n", "zsh", config)
	}

	sh := g.Out.String()
	fmt.Println("Output: " + sh)
	hook = &ShellHook{
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
