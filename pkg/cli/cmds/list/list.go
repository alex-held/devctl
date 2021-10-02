package list

import (
	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/pkg/cli/options"
	"github.com/alex-held/devctl/pkg/cli/util"
	"github.com/alex-held/devctl/pkg/env"
)

type ListOptions struct {
	options.IOStreams

	Remote     bool
	Upgradable bool
	All        bool
}

// NewListOptions returns an initialized ListOptions instance
func NewListOptions(ioStreams options.IOStreams) *ListOptions {
	return &ListOptions{
		IOStreams: ioStreams,
	}
}

// ValidateArgs makes sure there is no discrepancy in command options
func (o *ListOptions) ValidateArgs(cmd *cobra.Command, args []string) error {
	return nil
}

// Complete completes all the required options
func (o *ListOptions) Complete(f env.Factory, cmd *cobra.Command) error {
	// implement for completions
	// https://pkg.go.dev/k8s.io/cli-runtime/pkg/genericclioptions#RecordFlags.Complete
	return nil
}

// Run performs the listing of plugins
func (o *ListOptions) Run(f env.Factory, cmd *cobra.Command) error {
	return nil
}

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdCreate(f env.Factory, ioStreams options.IOStreams) *cobra.Command {
	o := NewListOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "create -f FILENAME",
		DisableFlagsInUseLine: true,
		Short:                 "lists devctl plugins",
		Long:                  "lists devctl plugins",
		Example: `
		To list installed plugins:
			devctl list
		To list remote plugins from the registry:
			devctl list --remote
		To list all plugins:
			devctl list --all
		To list upgradable plugins
			devctl list --upgradeable`,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Complete(f, cmd))
			util.CheckErr(o.ValidateArgs(cmd, args))
			util.CheckErr(o.Run(f, cmd))
		},
	}

	// bind flag structs
	// uncomment for completions
	// o.RecordFlags.AddFlags(cmd)

	cmd.Flags().BoolVar(&o.Remote, "remote", o.Remote, "List plugins present on remote")
	cmd.Flags().BoolVarP(&o.All, "all", "a", o.All, "List all plugins (remote and local)")
	cmd.Flags().BoolVar(&o.Upgradable, "upgradable", o.Upgradable, "List installed plugins which can be upgraded")

	// usage := "to use to create the resource"
	// cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)
	// cmdutil.AddValidateFlags(cmd)
	// cmd.Flags().BoolVar(&o.EditBeforeCreate, "edit", o.EditBeforeCreate, "Edit the API resource before creating")
	// cmd.Flags().Bool("windows-line-endings", runtime.GOOS == "windows",
	// 	"Only relevant if --edit=true. Defaults to the line ending native to your platform.")
	// cmdutil.AddApplyAnnotationFlags(cmd)
	// cmdutil.AddDryRunFlag(cmd)
	// cmd.Flags().StringVarP(&o.Selector, "selector", "l", o.Selector, "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	// cmd.Flags().StringVar(&o.Raw, "raw", o.Raw, "Raw URI to POST to the server.  Uses the transport specified by the kubeconfig file.")
	// cmdutil.AddFieldManagerFlagVar(cmd, &o.fieldManager, "kubectl-create")
	//
	// o.PrintFlags.AddFlags(cmd)

	// create subcommands
	// cmd.AddCommand(NewCmdCreateNamespace(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateQuota(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateSecret(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateConfigMap(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateServiceAccount(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateService(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateDeployment(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateClusterRole(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateClusterRoleBinding(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateRole(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateRoleBinding(f, ioStreams))
	// cmd.AddCommand(NewCmdCreatePodDisruptionBudget(f, ioStreams))
	// cmd.AddCommand(NewCmdCreatePriorityClass(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateJob(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateCronJob(f, ioStreams))
	// cmd.AddCommand(NewCmdCreateIngress(f, ioStreams))
	return cmd
}
