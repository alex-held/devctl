// +build tools

// Package tools
package tools

import (
	_ "github.com/axw/gocov"
	_ "github.com/mattn/goveralls/tester"
	_ "github.com/modocache/gover/gover"
	_ "github.com/onsi/ginkgo"
	_ "golang.org/x/tools/cover"
)
