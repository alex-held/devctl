package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
	
	config2 "github.com/alex-held/dev-env/config"
)

func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("config called")
		},
	}
	
	cmd.AddCommand(newConfigViewCommand())
	
	return cmd
}

func newConfigViewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Displays the current Configuration",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			devEnvConfig := config2.LoadViperConfig()
			configString := fmt.Sprintf("%+v\n", *devEnvConfig)
			fmt.Println(configString)
		},
	}
}
