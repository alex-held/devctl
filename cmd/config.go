package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/internal/config"

	"github.com/alex-held/devctl/internal/cli"
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
			c := cli.GetOrCreateCLI()
			fmt.Println(c.ConfigFileName())
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
			cfg := config.LoadViperConfig()
			err := printConfig(cfg)
			if err != nil {
				cli.ExitWithError(1, err)
			}
			os.Exit(0)
		},
	}
}

func printConfig(cfg *config.DevEnvConfig) (err error) {
	_, err = fmt.Fprintf(os.Stdout, "%#v", cfg)
	return err
}
