package plugin

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/pkg/constants"
	"github.com/alex-held/devctl/pkg/env"
	"github.com/alex-held/devctl/pkg/index/installation"
	"github.com/alex-held/devctl/pkg/index/scanner"
)

var (
	forceIndexDelete    *bool
	errInvalidIndexName = errors.New("invalid index name")
)

func NewIndexCommand(f env.Factory) (cmd *cobra.Command) {

	cmd = &cobra.Command{
		Use:   "index",
		Short: "Manage custom plugin indexes",
		Long:  "Manage which repositories are used to discover and install plugins from.",
		Args:  cobra.NoArgs,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List configured indexes",
		Long: `Print a list of configured indexes.
This command prints a list of indexes. It shows the name and the remote URL for
each configured index in table format.`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			indexes, err := scanner.ListIndexes(f.Paths())
			if err != nil {
				return errors.Wrap(err, "failed to list indexes")
			}

			var rows [][]string
			for _, index := range indexes {
				rows = append(rows, []string{index.Name, index.URL})
			}
			return printTable(os.Stdout, []string{"INDEX", "URL"}, rows)
		},
	}

	var addCmd = &cobra.Command{
		Use:     "add",
		Short:   "Add a new index",
		Long:    "Configure a new index to install plugins from.",
		Example: "kubectl krew index add default " + constants.DefaultIndexURI,
		Args:    cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			if !scanner.IsValidIndexName(name) {
				return errInvalidIndexName
			}
			err := scanner.AddIndex(f.Paths(), name, args[1])
			if err != nil {
				return err
			}
			f.Logger().Warnf(`You have added a new index from %q
The plugins in this index are not audited for security by the Krew maintainers.
Install them at your own risk.
`, args[1])
			return nil
		},
	}

	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a configured index",
		Long: `Removes a configured plugin index
It is only safe to remove indexes without installed plugins. Removing an index
while there are plugins installed will result in an error, unless the --force
option is used (not recommended).`,

		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			if !scanner.IsValidIndexName(name) {
				return errInvalidIndexName
			}

			ps, err := installation.InstalledPluginsFromIndex(f, name)
			if err != nil {
				return errors.Wrap(err, "failed querying plugins installed from the index")
			}
			f.Logger().Infof("Found %d plugins from index", len(ps))

			if len(ps) > 0 && !*forceIndexDelete {
				var names []string
				for _, pl := range ps {
					names = append(names, pl.Name)
				}

				f.Logger().Warnf(`Plugins [%s] are still installed from index %q!
Removing indexes while there are plugins installed from is not recommended
(you can use --force to ignore this check).`+"\n", strings.Join(names, ", "), name)
				return errors.Errorf("there are still plugins installed from this index")
			}

			err = scanner.DeleteIndex(f.Paths(), name)
			if os.IsNotExist(err) {
				if *forceIndexDelete {
					f.Logger().Infof("Index not found, but --force is used, so not returning an error")
					return nil // success if --force specified and index does not exist.
				}
				return errors.Errorf("index %q does not exist", name)
			}
			return errors.Wrap(err, "error while removing the plugin index")
		},
	}

	forceIndexDelete = removeCmd.Flags().Bool("force", false,
		"Remove index even if it has plugins currently installed (may result in unsupported behavior)")

	cmd.AddCommand(addCmd)
	cmd.AddCommand(listCmd)
	cmd.AddCommand(removeCmd)

	return cmd
}
