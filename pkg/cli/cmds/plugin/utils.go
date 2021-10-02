package plugin

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/alex-held/devctl/internal/git"
	"github.com/alex-held/devctl/pkg/constants"
	"github.com/alex-held/devctl/pkg/env"
	"github.com/alex-held/devctl/pkg/index/spec"
)

var canonicalNameRegex = regexp.MustCompile(`^[\w-]+/[\w-]+$`)

// indexOf returns the index name of a receipt.
func indexOf(r spec.Receipt) string {
	if r.Status.Source.Name == "" {
		return constants.DefaultIndexName
	}
	return r.Status.Source.Name
}

// displayName returns the display name of a Plugin.
// The index name is omitted if it is the default index.
func displayName(p spec.Plugin, indexName string) string {
	if isDefaultIndex(indexName) {
		return p.Name
	}
	return indexName + "/" + p.Name
}

func isDefaultIndex(name string) bool {
	return name == "" || name == constants.DefaultIndexName
}

// canonicalName returns INDEX/NAME value for a plugin, even if
// it is in the default index.
func canonicalName(p spec.Plugin, indexName string) string {
	if isDefaultIndex(indexName) {
		indexName = constants.DefaultIndexName
	}
	return indexName + "/" + p.Name
}

func isCanonicalName(s string) bool {
	return canonicalNameRegex.MatchString(s)
}

func printTable(out io.Writer, columns []string, rows [][]string) error {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, strings.Join(columns, "\t"))
	fmt.Fprintln(w)
	for _, values := range rows {
		fmt.Fprint(w, strings.Join(values, "\t"))
		fmt.Fprintln(w)
	}
	return w.Flush()
}

func sortByFirstColumn(rows [][]string) [][]string {
	sort.Slice(rows, func(a, b int) bool {
		return rows[a][0] < rows[b][0]
	})
	return rows
}

func checkIndex(f env.Factory, _ *cobra.Command, _ []string) error {
	if ok, err := git.IsGitCloned(f.Paths().IndexPath(constants.DefaultIndexName)); err != nil {
		return errors.Wrap(err, "failed to check local index git repository")
	} else if !ok {
		return errors.New(`krew local plugin index is not initialized (run "kubectl krew update")`)
	}
	return nil
}

func ensureDirs(paths ...string) error {
	for _, p := range paths {
		klog.V(4).Infof("Ensure creating dir: %q", p)
		if err := os.MkdirAll(p, 0755); err != nil {
			return errors.Wrapf(err, "failed to ensure create directory %q", p)
		}
	}
	return nil
}

func isTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}
