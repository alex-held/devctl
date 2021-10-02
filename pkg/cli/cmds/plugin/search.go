package plugin

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/alex-held/devctl/pkg/env"
	"github.com/alex-held/devctl/pkg/index/installation"
	"github.com/alex-held/devctl/pkg/index/scanner"
	"github.com/alex-held/devctl/pkg/index/spec"
)

type pluginEntry struct {
	p         spec.Plugin
	indexName string
}

// newSearchCmd creates the 'devctl index search' commands
func newSearchCmd(f env.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Discover devctl plugins",
		Long: `List devctl plugins available and search among them.
If no arguments are provided, all plugins will be listed.
Examples:
  To list all plugins:
    devctl index search
  To fuzzy search plugins with a keyword:
    devctl index search KEYWORD`,
		RunE: func(cmd *cobra.Command, args []string) error {
			indexes, err := scanner.ListIndexes(f.Paths())
			if err != nil {
				return errors.Wrap(err, "failed to list indexes")
			}

			klog.V(3).Infof("found %d indexes", len(indexes))

			var plugins []pluginEntry
			for _, idx := range indexes {
				ps, err := scanner.LoadPluginsFromFS(f, idx.Name)
				if err != nil {
					return errors.Wrapf(err, "failed to load the list of plugins from the index %q", idx.Name)
				}
				for _, p := range ps {
					plugins = append(plugins, pluginEntry{p, idx.Name})
				}
			}

			pluginCanonicalNames := make([]string, len(plugins))
			pluginCanonicalNameMap := make(map[string]pluginEntry, len(plugins))
			for i, p := range plugins {
				cn := canonicalName(p.p, p.indexName)
				pluginCanonicalNames[i] = cn
				pluginCanonicalNameMap[cn] = p
			}

			installed := make(map[string]bool)
			receipts, err := installation.GetInstalledPluginReceipts(f)
			if err != nil {
				return errors.Wrap(err, "failed to load installed plugins")
			}
			for _, receipt := range receipts {
				cn := canonicalName(receipt.Plugin, indexOf(receipt))
				installed[cn] = true
			}

			var searchResults []string
			if len(args) > 0 {
				matches := fuzzy.Find(strings.Join(args, ""), pluginCanonicalNames)
				for _, m := range matches {
					searchResults = append(searchResults, m.Str)
				}
			} else {
				searchResults = pluginCanonicalNames
			}

			// No plugins found
			if len(searchResults) == 0 {
				return nil
			}

			var rows [][]string
			cols := []string{"NAME", "DESCRIPTION", "INSTALLED"}
			for _, canonicalName := range searchResults {
				v := pluginCanonicalNameMap[canonicalName]
				var status string
				if installed[canonicalName] {
					status = "yes"
				} else if _, ok, err := installation.GetMatchingPlatform(v.p.Spec.Platforms); err != nil {
					return errors.Wrapf(err, "failed to get the matching platform for plugin %s", canonicalName)
				} else if ok {
					status = "no"
				} else {
					status = fmt.Sprintf("unavailable on %v/%v", runtime.GOOS, runtime.GOARCH)
				}

				rows = append(rows, []string{displayName(v.p, v.indexName), limitString(v.p.Spec.ShortDescription, 50), status})
			}
			rows = sortByFirstColumn(rows)
			return printTable(f.Streams().Out, cols, rows)
		},

		PreRunE: func(c *cobra.Command, args []string) error {
			return checkIndex(f, c, args)
		},
	}

	return cmd
}

func limitString(s string, length int) string {
	if len(s) > length && length > 3 {
		s = s[:length-3] + "..."
	}
	return s
}
