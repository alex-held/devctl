package zsh

import (
	"bytes"
	"context"
	"testing"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
	"github.com/alex-held/devctl-kit/pkg/plugins"
)

func TestInit(t *testing.T) {
	type vars struct {
		cmd      string
		args     []string
		cfg      Config
		expected string
	}

	tts := []struct {
		name string
		vars vars
	}{
		{
			name: "init",
			vars: vars{
				cmd:  "init",
				args: []string{""},
				cfg: Config{
					Completions: map[string]string{
						"kubectl":       "source <(kubectl completion zsh)",
						"golangci-lint": "source <(golangci-lint completion zsh)",
						"gh":            "source <(gh completion -s zsh)",
						"operator-sdk":  "source <(operator-sdk completion zsh)",
						"kompose":       "source <(kompose completion zsh)",
						"kubebuilder":   "source <(kubebuilder completion zsh)",
						"kustomize":     "source <(kustomize completion zsh)",
					},
					Context: &plugins.Context{
						Out: nil,
						Pather: devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
							return "/home/user/.devctl"
						})),
						Context: context.Background(),
					},
				},
				expected: "typeset -A cmds=(\n    [\"golangci-lint\"]=\"completion zsh\"\n    [gh]=\"completion -s zsh\"\n    [npm]=\"completion\"\n    [netlify]=\"completion:generate --shell=zsh\"\n    [ionic]=\"completion\"\n    [kubectl]=\"completion zsh\"\n    [kompose]=\"completion zsh\"\n    [kubebuilder]=\"completion zsh\"\n    [kustomize]=\"completion zsh\"\n    [\"operator-sdk\"]=\"completion zsh\"\n)\nfor k v in ${(kv)cmds}; do\n    [[ ! -f $ZINIT[COMPLETIONS_DIR]/__tabtab.zsh ]] && \\\n        eval \"$k $v\" > \"$ZINIT[COMPLETIONS_DIR]/_$k\" || true\n    source \"$ZINIT[COMPLETIONS_DIR]/_$k\"\ndone\nunset cmds\n",
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			tt.vars.cfg.Out = out
		})
	}
}
