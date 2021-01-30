package cmd_backup

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
	
dev-env use go 1.14.x
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

func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dev-env/config.yaml")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dev-env.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(currentCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	/*
		if cfgFile != "" {

			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Search config in  "~/.dev-env/config.yaml"
			viper.AddConfigPath(path.Join(home, ".dev-env"))
			viper.SetConfigName("config.yaml")

			viper.Debug()
		}

		// viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}*/
}
