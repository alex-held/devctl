package zsh

import (
	"bytes"
	"testing"

	"github.com/alex-held/gold"
	"github.com/stretchr/testify/assert"

	"github.com/alex-held/devctl-kit/pkg/plugins"

	"github.com/alex-held/devctl-kit/pkg/generation/banner"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
)

var config = Config{
	Context: &plugins.Context{},
	Vars: map[string]string{
		"version": "v1.0.0",
	},
	Exports: map[string]string{
		"GOPATH":        "$HOME/go",
		"EXPORT_STRING": "abc",
		"EXPORT_BOOL":   "true",
		"EXPORT_NUMBER": "1",
	},
	Aliases: map[string]string{
		"k":   "kubectl",
		"bq":  "bazel query '...' | fzf",
		"cdg": "cd $GOPATH",
		"cdr": "cd $HOME/source/repos",
	},
	Completions: CompletionsSpec{
		CLI: map[string]string{
			"ionic":         "completion",
			"npm":           "completion",
			"gh":            "completion -s zsh",
			"netlify":       "completion:generate --shell=zsh",
			"golangci-lint": "completion zsh",
			"kubectl":       "completion zsh",
			"kompose":       "completion zsh",
			"kubebuilder":   "completion zsh",
			"operator-sdk":  "completion zsh",
		},
		Command: map[string]string{
			"gobuffalo": "if [[ ! -f \"${ZINIT[COMPLETIONS_DIR]}/_buffalo\" ]]; then\n  if [[ -f \"$ZSH/custom/gobuffalo.zsh/buffalo.plugin.zsh\" ]]; then\n    ln -s $ZSH/custom/gobuffalo.zsh/buffalo.plugin.zsh $ZINIT[COMPLETIONS_DIR]/_buffalo\n    source $ZINIT[COMPLETIONS_DIR]/_buffalo\n  fi\nfi",
			"nvm":       "[ -s \"/usr/local/opt/nvm/nvm.sh\" ] && . \"/usr/local/opt/nvm/nvm.sh\"\n[ -s \"/usr/local/opt/nvm/etc/bash_completion.d/nvm\" ] && . \"/usr/local/opt/nvm/etc/bash_completion.d/nvm\"",
		},
	},
}

func TestGenerateExports(t *testing.T) {
	gen := NewGenerator()

	w := &bytes.Buffer{}
	err := gen.Exports(w)
	assert.NoError(t, err)
}

func TestCompletionsTempl(t *testing.T) {
	completionsDir := devctlpath.DevCtlConfigRoot("configs", "zsh", "completions")
	configFile := devctlpath.DevCtlConfigRoot("configs", "zsh", "completions.yaml")

	in := CompletionsTmplData{
		FileHeader: FileHeaderTmplData{
			CONFIGFILE: configFile,
		},
		Header:          banner.GenerateBanner("Completions", banner.KIND_SHELL),
		COMPLETIONS_DIR: completionsDir,
		Completions: CompletionsSpec{
			CLI: map[string]string{
				"ionic":         "completion",
				"npm":           "completion",
				"gh":            "completion -s zsh",
				"netlify":       "completion:generate --shell=zsh",
				"golangci-lint": "completion zsh",
				"kubectl":       "completion zsh",
				"kompose":       "completion zsh",
				"kubebuilder":   "completion zsh",
				"operator-sdk":  "completion zsh",
			},
			Command: map[string]string{
				"gobuffalo": "if [[ ! -f \"${ZINIT[COMPLETIONS_DIR]}/_buffalo\" ]]; then\n  if [[ -f \"$ZSH/custom/gobuffalo.zsh/buffalo.plugin.zsh\" ]]; then\n    ln -s $ZSH/custom/gobuffalo.zsh/buffalo.plugin.zsh $ZINIT[COMPLETIONS_DIR]/_buffalo\n    source $ZINIT[COMPLETIONS_DIR]/_buffalo\n  fi\nfi",
				"nvm":       "[ -s \"/usr/local/opt/nvm/nvm.sh\" ] && . \"/usr/local/opt/nvm/nvm.sh\"\n[ -s \"/usr/local/opt/nvm/etc/bash_completion.d/nvm\" ] && . \"/usr/local/opt/nvm/etc/bash_completion.d/nvm\"",
			},
		},
	}

	actual, err := in.Render()
	assert.NoError(t, err)

	g := gold.New(t)
	g.Assert(t, "render", []byte(actual))
}
