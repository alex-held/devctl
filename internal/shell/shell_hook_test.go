package shell

import (
	"html/template"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	assert2 "github.com/stretchr/testify/assert"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

func setup(t *testing.T, content string) (filepath string, file afero.File, fs afero.Fs) {
	fs = afero.NewMemMapFs()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Cannot setup shellrc file for test. error=%v", err)
	}
	filepath = path.Join(home, ".zshrc")
	file, err = fs.Create(filepath)
	if err != nil {
		t.Fatalf("Cannot setup shellrc file for test. error=%v", err)
	}
	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Cannot setup shellrc file for test. error=%v", err)
	}
	return filepath, file, fs
}

const defaultShellRc = "# ZSHRC - START\n\nsource $DOTFILES/zsh/init.zsh\n\n# ZSHRC - END\n"

type ShellRC struct {
	Before        string
	DevCtlSection struct {
		EnvVars map[string]string
	}
	After string
}

func (rc *ShellRC) Render() (rendered string, err error) {
	tmpl, err := template.ParseGlob("templates/shellrc.tmpl")

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

type DevCtlInjector struct {
	RC     *ShellRC
	Pather devctlpath.Pather
}

func (i *DevCtlInjector) AddDevCtlSection() {

}

func AddDevCtlEnvVars() {

}
