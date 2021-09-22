package zsh

import (
	"os"
	"testing"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
	"github.com/stretchr/testify/assert"
)

func TestCompletionsTempl(t *testing.T) {
	completionsDir := devctlpath.DevCtlConfigRoot("configs", "zsh", "completions")
	configFile := devctlpath.DevCtlConfigRoot("configs", "zsh", "completions.yaml")

	in := CompletionsTmplData{
		FileHeader: FileHeaderTmplData{
			CONFIGFILE: configFile,
		},
		Header:          GenerateBanner("Completions"),
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

	if *genGoldenMaster {
		_ = os.WriteFile("testdata/completions_test.golden", []byte(actual), os.ModePerm)
		return
	}

	t.Log(actual)
	assert.Equal(t, completionsExpectedString, actual)
}
