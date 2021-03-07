package shell

import "text/template"

var defaults []Option

type Option func(*ShellHookConfig) *ShellHookConfig

func WithTemplates(t *template.Template) Option {
	return func(cfg *ShellHookConfig) *ShellHookConfig {
		cfg.Templates = t
		return cfg
	}
}

func NewShellHookConfig(opts ...Option) *ShellHookConfig {
	cfg := &ShellHookConfig{
		Templates: nil,
		Sections:  Sections{},
	}

	for _, opt := range defaults {
		cfg = opt(cfg)
	}

	for _, opt := range opts {
		cfg = opt(cfg)
	}

	return cfg
}
