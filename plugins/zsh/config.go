package zsh

type ConfigFile struct {
	Exports     map[string]string `yaml:"exports,omitempty"`
	Completions CompletionsSpec   `yaml:"completions,omitempty"`
}

type CompletionsSpec struct {
	CLI     map[string]string `yaml:"cli,omitempty"`
	Command map[string]string `yaml:"command,omitempty"`
}
