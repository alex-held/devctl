package printutils

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/alex-held/devctl/pkg/index/installation"
	"github.com/alex-held/devctl/pkg/index/spec"
)

func printPluginInfo(out io.Writer, indexName string, plugin spec.Plugin) {
	fmt.Fprintf(out, "NAME: %s\n", plugin.Name)
	fmt.Fprintf(out, "INDEX: %s\n", indexName)
	if platform, ok, err := installation.GetMatchingPlatform(plugin.Spec.Platforms); err == nil && ok {
		if platform.URI != "" {
			fmt.Fprintf(out, "URI: %s\n", platform.URI)
			fmt.Fprintf(out, "SHA256: %s\n", platform.Sha256)
		}
	}
	if plugin.Spec.Version != "" {
		fmt.Fprintf(out, "VERSION: %s\n", plugin.Spec.Version)
	}
	if plugin.Spec.Homepage != "" {
		fmt.Fprintf(out, "HOMEPAGE: %s\n", plugin.Spec.Homepage)
	}
	if plugin.Spec.Description != "" {
		fmt.Fprintf(out, "DESCRIPTION: \n%s\n", plugin.Spec.Description)
	}
	if plugin.Spec.Caveats != "" {
		fmt.Fprintf(out, "CAVEATS:\n%s\n", Indent(plugin.Spec.Caveats))
	}
}

// Indent converts strings to an indented format ready for printing.
// Example:
//
//     \
//      | This plugin is great, use it with great care.
//      | Also, plugin will require the following programs to run:
//      |  * jq
//      |  * base64
//     /
func Indent(s string) string {
	out := "\\\n"
	s = strings.TrimRightFunc(s, unicode.IsSpace)
	out += regexp.MustCompile("(?m)^").ReplaceAllString(s, " | ")
	out += "\n/"
	return out
}
