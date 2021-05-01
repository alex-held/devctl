// +build tools

package tools

import (
	_ "github.com/onsi/ginkgo/ginkgo"


	_ "github.com/axw/gocov/gocov"
	_ "github.com/mattn/goveralls"
	_ "github.com/modocache/gover"
	_ "golang.org/x/tools/cmd/cover"
)
