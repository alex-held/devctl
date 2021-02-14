package cmd

import (
	"fmt"
	"os"
	
	"github.com/spf13/cobra"
	
	"github.com/alex-held/devctl/pkg/cli"
)

var cfgFile string


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dev-env",
	Short: "A lightweight dev-environment manager/bootstrapper.",
	Long: `dev-env can manage all kinds of sdks, runtime dependencies, plugins, EnvVars and directories.
Examples and usage of using your application. For example:

dev-env install java 14.0.1

dev-env current java

dev-env config view

dev-env list

dev-env list dotnet

dev-env use go 1.15.x
`,
	/*	PersistentPreRun: funcutil(cmd *cobra.Command, args []string) {
		initConfig()
	},*/
	
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: funcutil(cmd *cobra.Command, args []string) { },
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
	)
	
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devctl/config.yaml)")
	
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	_ = cli.GetOrCreateCLI()
}
