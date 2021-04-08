package meta

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func IsDevctl(mod string) bool {
	if !strings.HasPrefix(mod, "go.mod") {
		mod = filepath.Join(mod, "go.mod")
	}

	b, err := ioutil.ReadFile(mod)
	if err != nil {
		return false
	}

	return bytes.Contains(b, []byte("github.com/alex-held/devctl"))
}
