package banner_test

import (
	"testing"

	"github.com/sebdah/goldie/v2"

	"github.com/alex-held/devctl/pkg/generation/banner"
)

func TestGenerateHeader(t *testing.T) {
	tts := []string{"Completions", "ZSH", "Exports", "Aliases"}

	for _, tt := range tts {
		t.Run(tt, func(t *testing.T) {
			actual := banner.GenerateBanner(tt, banner.KIND_SHELL)

			g := goldie.New(t, goldie.WithTestNameForDir(true), goldie.WithDiffEngine(goldie.ColoredDiff))
			g.Assert(t, tt, []byte(actual))
		})
	}
}
