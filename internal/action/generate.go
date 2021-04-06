package action

import (
	"bytes"
	"os"
	"text/template"
)

type Generate action

const tmpl = `#!/usr/bin/env zsh
antibody() {
	case "$1" in
	bundle)
		source <( {{ . }} $@ ) || {{ . }} $@
		;;
	*)
		{{ . }} $@
		;;
	esac
}
_antibody() {
	IFS=' ' read -A reply <<< "help bundle update home purge list init"
}
compctl -K _antibody antibody
`

func Init() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	var tmpl = template.Must(template.New("init").Parse(tmpl))
	var out bytes.Buffer
	err = tmpl.Execute(&out, executable)
	return out.String(), err
}
