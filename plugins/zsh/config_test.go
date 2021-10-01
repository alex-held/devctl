package zsh

import (
	"strings"
	"testing"

	"github.com/alex-held/gold"
	"github.com/stretchr/testify/assert"
)

func TestReadConfigFile(t *testing.T) {
	g := gold.New(t)
	cfgYaml, _ := g.Get(t, "config")
	actual, err := ReadConfigFile(strings.NewReader(cfgYaml))
	assert.NoError(t, err)

	g.AssertYaml(t, "config", actual, 4)
}
