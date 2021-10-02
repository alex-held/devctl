package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/pkg/cli/cmds/info"
	"github.com/alex-held/devctl/pkg/cli/cmds/list"
	"github.com/alex-held/devctl/pkg/cli/cmds/plugin"
	cliflag "github.com/alex-held/devctl/pkg/cli/flags"
	"github.com/alex-held/devctl/pkg/cli/options"
	"github.com/alex-held/devctl/pkg/cli/templates"
	"github.com/alex-held/devctl/pkg/cli/util"
	"github.com/alex-held/devctl/pkg/env"
)

// NewDefaultKubectlCommand creates the `kubectl` command with default arguments
func NewDefaultKubectlCommand() *cobra.Command {
	return NewDefaultKubectlCommandWithArgs(util.NewDefaultPluginHandler(), os.Args, os.Stdin, os.Stdout, os.Stderr)
}

// NewDefaultKubectlCommandWithArgs creates the `kubectl` command with arguments
func NewDefaultKubectlCommandWithArgs(pluginHandler util.PluginHandler, args []string, in io.Reader, out, errout io.Writer) *cobra.Command {
	cmd := NewDevctlCommand(in, out, errout)

	if pluginHandler == nil {
		return cmd
	}

	if len(args) > 1 {
		cmdPathPieces := args[1:]

		// only look for suitable extension executables if
		// the specified command does not already exist
		if _, _, err := cmd.Find(cmdPathPieces); err != nil {
			if err := HandlePluginCommand(pluginHandler, cmdPathPieces); err != nil {
				fmt.Fprintf(errout, "Error: %v\n", err)
				os.Exit(1)
			}
		}
	}

	return cmd
}

// HandlePluginCommand receives a pluginHandler and command-line arguments and attempts to find
// a plugin executable on the PATH that satisfies the given arguments.
func HandlePluginCommand(pluginHandler util.PluginHandler, cmdArgs []string) error {
	var remainingArgs []string // all "non-flag" arguments
	for _, arg := range cmdArgs {
		if strings.HasPrefix(arg, "-") {
			break
		}
		remainingArgs = append(remainingArgs, strings.Replace(arg, "-", "_", -1))
	}

	if len(remainingArgs) == 0 {
		// the length of cmdArgs is at least 1
		return fmt.Errorf("flags cannot be placed before plugin name: %s", cmdArgs[0])
	}

	foundBinaryPath := ""

	// attempt to find binary, starting at longest possible name with given cmdArgs
	for len(remainingArgs) > 0 {
		path, found := pluginHandler.Lookup(strings.Join(remainingArgs, "-"))
		if !found {
			remainingArgs = remainingArgs[:len(remainingArgs)-1]
			continue
		}

		foundBinaryPath = path
		break
	}

	if len(foundBinaryPath) == 0 {
		return nil
	}

	// invoke cmd binary relaying the current environment and args given
	if err := pluginHandler.Execute(foundBinaryPath, cmdArgs[len(remainingArgs):], os.Environ()); err != nil {
		return err
	}

	return nil
}

func NewDevctlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	warningsAsErrors := false

	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "devctl",
		Short: "devctl manges your development environment",
		Long: `
      devctl manges your development environment.
      Find more information at:
            https://github.com/alex-held/devctl/docs/overview.md`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	flags := cmds.PersistentFlags()
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	flags.BoolVar(&warningsAsErrors, "warnings-as-errors", warningsAsErrors, "Treat warnings received from the server as errors and exit with a non-zero exit code")

	devctlFlags := options.NewConfigFlags()
	devctlFlags.AddFlags(flags)

	// Updates hooks to add kubectl command headers: SIG CLI KEP 859.
	// addCmdHeaderHooks(cmds, kubeConfigFlags)

	cmds.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// From this point and forward we get warnings on flags that contain "_" separators
	cmds.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	ioStreams := options.IOStreams{In: in, Out: out, ErrOut: err}

	f := env.NewFactory()

	groups := templates.CommandGroups{
		{
			Message: "Basic Commands (Beginner):",
			Commands: []*cobra.Command{
				list.NewCmdCreate(f, ioStreams),
				info.NewCmd(f, ioStreams),
			},
		},
		{
			Message: "Advanced Commands (Plugins):",
			Commands: []*cobra.Command{
				plugin.NewCmd(f, ioStreams),
			},
		},
		// {
		// 	Message: "Deploy Commands:",
		// 	Commands: []*cobra.Command{
		// 		rollout.NewCmdRollout(f, ioStreams),
		// 		scale.NewCmdScale(f, ioStreams),
		// 		autoscale.NewCmdAutoscale(f, ioStreams),
		// 	},
		// },
		// {
		// 	Message: "Cluster Management Commands:",
		// 	Commands: []*cobra.Command{
		// 		certificates.NewCmdCertificate(f, ioStreams),
		// 		clusterinfo.NewCmdClusterInfo(f, ioStreams),
		// 		top.NewCmdTop(f, ioStreams),
		// 		drain.NewCmdCordon(f, ioStreams),
		// 		drain.NewCmdUncordon(f, ioStreams),
		// 		drain.NewCmdDrain(f, ioStreams),
		// 		taint.NewCmdTaint(f, ioStreams),
		// 	},
		// },
		// {
		// 	Message: "Troubleshooting and Debugging Commands:",
		// 	Commands: []*cobra.Command{
		// 		describe.NewCmdDescribe("kubectl", f, ioStreams),
		// 		logs.NewCmdLogs(f, ioStreams),
		// 		attach.NewCmdAttach(f, ioStreams),
		// 		cmdexec.NewCmdExec(f, ioStreams),
		// 		portforward.NewCmdPortForward(f, ioStreams),
		// 		proxyCmd,
		// 		cp.NewCmdCp(f, ioStreams),
		// 		auth.NewCmdAuth(f, ioStreams),
		// 		debug.NewCmdDebug(f, ioStreams),
		// 	},
		// },
		// {
		// 	Message: "Advanced Commands:",
		// 	Commands: []*cobra.Command{
		// 		diff.NewCmdDiff(f, ioStreams),
		// 		apply.NewCmdApply("kubectl", f, ioStreams),
		// 		patch.NewCmdPatch(f, ioStreams),
		// 		replace.NewCmdReplace(f, ioStreams),
		// 		wait.NewCmdWait(f, ioStreams),
		// 		kustomize.NewCmdKustomize(ioStreams),
		// 	},
		// },
		// {
		// 	Message: "Settings Commands:",
		// 	Commands: []*cobra.Command{
		// 		label.NewCmdLabel(f, ioStreams),
		// 		annotate.NewCmdAnnotate("kubectl", f, ioStreams),
		// 		completion.NewCmdCompletion(ioStreams.Out, ""),
		// 	},
		// },
	}
	groups.Add(cmds)

	//	filters := []string{"options"}

	// // Hide the "alpha" subcommand if there are no alpha commands in this build.
	// alpha := NewCmdAlpha(ioStreams)
	// if !alpha.HasSubCommands() {
	// 	filters = append(filters, alpha.Name())
	// }

	//	templates.ActsAsRootCommand(cmds, filters, groups...)

	// TODO: uncommment
	//	util.SetFactoryForCompletion(f)
	//	registerCompletionFuncForGlobalFlags(cmds, f)

	// cmds.AddCommand(alpha)
	// cmds.AddCommand(cmdconfig.NewCmdConfig(clientcmd.NewDefaultPathOptions(), ioStreams))
	// cmds.AddCommand(plugin.NewCmdPlugin(ioStreams))
	// cmds.AddCommand(version.NewCmdVersion(f, ioStreams))
	// cmds.AddCommand(apiresources.NewCmdAPIVersions(f, ioStreams))
	// cmds.AddCommand(apiresources.NewCmdAPIResources(f, ioStreams))
	// cmds.AddCommand(options.NewCmdOptions(ioStreams.Out))

	return cmds
}
