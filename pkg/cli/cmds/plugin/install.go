package plugin

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/alex-held/devctl-kit/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/alex-held/devctl/pkg/constants"
	"github.com/alex-held/devctl/pkg/env"
	"github.com/alex-held/devctl/pkg/index/installation"
	"github.com/alex-held/devctl/pkg/index/pathutil"
	"github.com/alex-held/devctl/pkg/index/printutils"
	"github.com/alex-held/devctl/pkg/index/scanner"
	"github.com/alex-held/devctl/pkg/index/spec"
	"github.com/alex-held/devctl/pkg/index/validate"
)

var (
	manifest, manifestURL, archiveFileOverride *string
	noUpdateIndex                              *bool
)

func NewInstallCmd(f env.Factory) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "install",
		Short: "Install kubectl plugins",
		Long: `Install one or multiple kubectl plugins.
Examples:
  To install one or multiple plugins from the default index, run:
    kubectl krew install NAME [NAME...]
  To install plugins from a newline-delimited file, run:
    kubectl krew install < file.txt
  To install one or multiple plugins from a custom index, run:
    kubectl krew install INDEX/NAME [INDEX/NAME...]
  (For developers) To provide a custom plugin manifest, use the --manifest or
  --manifest-url arguments. Similarly, instead of downloading files from a URL,
  you can specify a local --archive file:
    kubectl krew install --manifest=FILE [--archive=FILE]
Remarks:
  If a plugin is already installed, it will be skipped.
  Failure to install a plugin will not stop the installation of other plugins.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var pluginNames = make([]string, len(args))
			copy(pluginNames, args)

			if !isTerminal(os.Stdin) && (len(pluginNames) != 0 || *manifest != "") {
				fmt.Fprintln(os.Stderr, "WARNING: Detected stdin, but discarding it because of --manifest or args")
			}

			if !isTerminal(os.Stdin) && (len(pluginNames) == 0 && *manifest == "") {
				fmt.Fprintln(os.Stderr, "Reading plugin names via stdin")
				s := bufio.NewScanner(os.Stdin)
				s.Split(bufio.ScanLines)
				for s.Scan() {
					if name := s.Text(); name != "" {
						pluginNames = append(pluginNames, name)
					}
				}
			}

			if *manifest != "" && *manifestURL != "" {
				return errors.New("cannot specify --manifest and --manifest-url at the same time")
			}

			if len(pluginNames) != 0 && (*manifest != "" || *manifestURL != "") {
				return errors.New("must specify either specify either plugin names (via positional arguments or STDIN), or --manifest/--manifest-url; not both")
			}

			if *archiveFileOverride != "" && *manifest == "" && *manifestURL == "" {
				return errors.New("--archive can be specified only with --manifest or --manifest-url")
			}

			var install []pluginEntry
			for _, name := range pluginNames {
				indexName, pluginName := pathutil.CanonicalPluginName(name)
				if !validate.IsSafePluginName(pluginName) {
					return unsafePluginNameErr(pluginName)
				}

				plugin, err := scanner.LoadPluginByName(f, f.Paths().IndexPluginsPath(indexName), pluginName)
				if err != nil {
					if os.IsNotExist(err) {
						return errors.Errorf("plugin %q does not exist in the plugin index", name)
					}
					return errors.Wrapf(err, "failed to load plugin %q from the index", name)
				}
				install = append(install, pluginEntry{
					p:         plugin,
					indexName: indexName,
				})
			}

			if *manifest != "" {
				plugin, err := scanner.ReadPluginFromFile(f.Fs(), *manifest)
				if err != nil {
					return errors.Wrap(err, "failed to load plugin manifest from file")
				}
				install = append(install, pluginEntry{
					p:         plugin,
					indexName: "detached",
				})
			} else if *manifestURL != "" {
				plugin, err := readPluginFromURL(*manifestURL)
				if err != nil {
					return errors.Wrap(err, "failed to read plugin manifest file from url")
				}
				install = append(install, pluginEntry{
					p:         plugin,
					indexName: "detached",
				})
			}

			if len(install) == 0 {
				return cmd.Help()
			}

			for _, pluginEntry := range install {
				klog.V(2).Infof("Will install plugin: %s/%s\n", pluginEntry.indexName, pluginEntry.p.Name)
			}

			var failed []string
			var returnErr error
			for _, entry := range install {
				plugin := entry.p
				f.Logger().Errorf("Installing plugin: %s\n", plugin.Name)
				err := installation.Install(f, plugin, entry.indexName, installation.InstallOpts{
					ArchiveFileOverride: *archiveFileOverride,
				})
				if err == installation.ErrIsAlreadyInstalled {
					klog.Warningf("Skipping plugin %q, it is already installed", plugin.Name)
					continue
				}
				if err != nil {
					klog.Warningf("failed to install plugin %q: %v", plugin.Name, err)
					if returnErr == nil {
						returnErr = err
					}
					failed = append(failed, plugin.Name)
					continue
				}
				fmt.Fprintf(os.Stderr, "Installed plugin: %s\n", plugin.Name)
				output := fmt.Sprintf("Use this plugin:\n\tkubectl %s\n", plugin.Name)
				if plugin.Spec.Homepage != "" {
					output += fmt.Sprintf("Documentation:\n\t%s\n", plugin.Spec.Homepage)
				}
				if plugin.Spec.Caveats != "" {
					output += fmt.Sprintf("Caveats:\n%s\n", printutils.Indent(plugin.Spec.Caveats))
				}
				fmt.Fprintln(os.Stderr, printutils.Indent(output))
				if entry.indexName == constants.DefaultIndexName {
					PrintSecurityNotice(plugin.Name)
				}
			}
			if len(failed) > 0 {
				return errors.Wrapf(returnErr, "failed to install some plugins: %+v", failed)
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if *manifest != "" {
				log.Warnf("--manifest specified, not ensuring plugin index")
				return nil
			}
			if *noUpdateIndex {
				log.Warnf("--no-update-index specified, skipping updating local copy of plugin index")
				return nil
			}

			return ensureIndexes(f, cmd, args)
		},
	}

	manifest = cmd.Flags().String("manifest", "", "(Development-only) specify local plugin manifest file")
	manifestURL = cmd.Flags().String("manifest-url", "", "(Development-only) specify plugin manifest file from url")
	archiveFileOverride = cmd.Flags().String("archive", "", "(Development-only) force all downloads to use the specified file")
	noUpdateIndex = cmd.Flags().Bool("no-update-index", false, "(Experimental) do not update local copy of plugin index before installing")

	return cmd
}

func PrintSecurityNotice(name string) {
	const securityNoticeFmt = `You installed plugin %q from the krew-index plugin repository.
   These plugins are not audited for security by the Krew maintainers.
   Run them at your own risk.`

	fmt.Printf(fmt.Sprintf(securityNoticeFmt, name))
}

func readPluginFromURL(url string) (spec.Plugin, error) {
	klog.V(4).Infof("downloading manifest from url %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return spec.Plugin{}, errors.Wrapf(err, "request to url failed (%s)", url)
	}
	klog.V(4).Infof("manifest downloaded from url, status=%v headers=%v", resp.Status, resp.Header)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return spec.Plugin{}, errors.Errorf("unexpected status code (http %d) from url", resp.StatusCode)
	}
	return scanner.ReadPlugin(resp.Body)
}
