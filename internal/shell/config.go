package shell

import "text/template"

var defaults = []Option{
	WithTemplates(template.Must(template.ParseGlob("templates/*.tmpl"))),
}

type Option func(*ShellHookConfig) *ShellHookConfig

func WithTemplates(t *template.Template) Option {
	return func(cfg *ShellHookConfig) *ShellHookConfig {
		cfg.Templates = t
		return cfg
	}
}

func WithSection(initializeSection UninitializedSection) Option {
	return func(cfg *ShellHookConfig) *ShellHookConfig {
		initializedSection := initializeSection(cfg.root)
		return cfg.AddOrUpdateSection(initializedSection)
	}
}

func WithSections(sections ...UninitializedSection) Option {
	return func(cfg *ShellHookConfig) *ShellHookConfig {
		for _, uninitializedSection := range sections {
			cfg = WithSection(uninitializedSection)(cfg)
		}
		return cfg
	}
}

func NewShellHookConfig(opts ...Option) *ShellHookConfig {
	cfg := &ShellHookConfig{
		Templates: nil,
		Sections:  Sections{},
	}
	cfg.root = newShellHookRootNode(cfg)

	for _, opt := range defaults {
		cfg = opt(cfg)
	}

	for _, opt := range opts {
		cfg = opt(cfg)
	}

	return cfg
}

func newShellHookRootNode(cfg *ShellHookConfig) *rootNode {
	root := &rootNode{config: cfg}
	root.node = &node{
		parent:   root.node,
		isRooted: true,
	}
	cfg.root = root
	root.config = cfg
	return root
}
