package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/internal/app"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dev-env",
	Short: "A weight dev-environment manager/bootstrapper.",
	Long: `dev-env can manage all kinds of sdks, runtime dependencies, plugins, EnvVars and directories.
Examples and usage of using your application. For example:

dev-env install java 14.0.1

dev-env current java

dev-env config view

dev-env list

dev-env list dotnet

dev-env use go 1.15.x
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//nolint:gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		NewCompletionCommand(),
		NewConfigCommand(),
		NewSdkCommand(),
		NewSdkManCommand(),
		NewPrefixCommand(),
		NewInitCommand(),
	)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&config, "config", "$DEVCTL_ROOT", "config file (default is $DEVCTL_ROOT=$HOME/.devctl/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "--verbose")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	gFlags := app.GlobalFlags{
		Verbose: verbose,
		Config:  config,
	}

	var args []string
	copy(args, os.Args)
	_ = app.GetOrCreateCLI(gFlags, args)
}
