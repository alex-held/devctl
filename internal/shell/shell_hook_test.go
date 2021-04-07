package shell

import (
	"html/template"
	"strings"
	"testing"

	"github.com/pkg/errors"
	assert2 "github.com/stretchr/testify/assert"
)

type ShellRC struct {
	Before        string
	DevCtlSection struct {
		EnvVars map[string]string
	}
	After string
}

func (rc *ShellRC) Render() (rendered string, err error) {
	tmpl, err := template.ParseGlob("templates/shellrc.tmpl")
	if err != nil {
		return "", err
	}
	sb := strings.Builder{}

	err = tmpl.ExecuteTemplate(&sb, "shellrc", *rc)
	if err != nil {
		return "", errors.Wrapf(err, "failed to render shellrc template. error=%v", err)
	}
	rendered = sb.String()
	return rendered, nil
}

func TestRender(t *testing.T) {
	sut := &ShellRC{
		Before: "",
		DevCtlSection: struct {
			EnvVars map[string]string
		}{
			EnvVars: map[string]string{
				"DEVCTL_PATH": "/homedir/.devctl",
			},
		},
		After: "",
	}

	expected := `

# --------------------------------------------------
# DEVCTL
# --------------------------------------------------
export DEVCTL_PATH=/homedir/.devctl


`
	actual, err := sut.Render()
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, expected, actual)
}
