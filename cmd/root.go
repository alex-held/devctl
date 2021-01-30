package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/coreos/etcd/client"
	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	cliName        = "devenv"
	cliDescription = "A lightweight dev-environment manager / bootstrapper "

	// The name of our config file, without the file extension because viper supports many different config file languages.
	defaultConfigFilename = "devenv"
	defaultConfigFileType = "yaml"

	envPrefix = "DEVENV"
)

func ExitWithError(code int, err error) {
	_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
	if cerr, ok := err.(*client.ClusterError); ok {
		_, _ = fmt.Fprintln(os.Stderr, cerr.Detail())
	}
	os.Exit(code)
}

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
	// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
	//	return initializeConfig()
	// },
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		NewConfigCommand(),
		NewSdkCommand(),
	)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devenv/viper/devenv.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initializeConfig() error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		devenv := path.Join(home, ".devenv", "debug")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in $HOME/.devenv/debug directory with name "devenv.yaml"
		viper.AddConfigPath(devenv)
		viper.SetEnvPrefix(envPrefix)
		viper.SetConfigName(defaultConfigFilename)
		viper.SetConfigType(defaultConfigFileType)
	}

	//	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		return err
	}

	return nil
}

func initConfig() {
	_ = initializeConfig()
}
