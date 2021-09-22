package zsh

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed "testdata/config.yaml"
var testConfigFile string

func TestReadConfigFile(t *testing.T) {

	expected := Config{
		Exports: map[string]string{
			"DEVCTL_ROOT": "$HOME/.devctl",
			"GOPATH":      "$HOME/go",
			"GOROOT":      "$DEVCTL_ROOT/sdks/go/current",
		},
		Completions: CompletionsSpec{
			CLI: map[string]string{
				"npm":           "completion",
				"gh":            "completion -s zsh",
				"kubectl":       "completion zsh",
				"kustomize":     "completion zsh",
				"operator-sdk":  "completion zsh",
				"golangci-lint": "completion zsh",
			},
			Command: map[string]string{
				"gobuffalo": "if [[ ! -f \"${ZINIT[COMPLETIONS_DIR]}/_buffalo\" ]]; then\n  if [[ -f \"$ZSH/custom/gobuffalo.zsh/buffalo.plugin.zsh\" ]]; then\n    ln -s $ZSH/custom/gobuffalo.zsh/buffalo.plugin.zsh $ZINIT[COMPLETIONS_DIR]/_buffalo\n    source $ZINIT[COMPLETIONS_DIR]/_buffalo\n  fi\nfi",
				"nvm":       "[ -s \"/usr/local/opt/nvm/nvm.sh\" ] && . \"/usr/local/opt/nvm/nvm.sh\"\n[ -s \"/usr/local/opt/nvm/etc/bash_completion.d/nvm\" ] && . \"/usr/local/opt/nvm/etc/bash_completion.d/nvm\"",
			},
		},
	}

	actual, err := ReadConfigFile(strings.NewReader(testConfigFile))
	assert.NoError(t, err)

	assert.Equal(t, expected, *actual)
}
