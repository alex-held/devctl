package main

import (
	"fmt"
	"github.com/alex-held/dev-env/manifest"
	. "github.com/ganbarodigital/go_scriptish"
)

func Main() {
	cmd := ExecPipeline(manifest.Symlink("/source/a/b", "/source/a/b"))
	fmt.Printf("%+v", cmd)
}
