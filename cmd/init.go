package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/internal/shell"
)

func NewInitCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:       "init",
		Long:      "init",
		Example:   "devctl init zsh",
		ValidArgs: []string{"zsh", "bash", "fish"},
	}

	cmd.AddCommand(newZshInitCommand())

	return cmd
}

func newZshInitCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:  "zsh",
		Long: "zsh",
		Example:   "devctl init zsh",
		Run: func(c *cobra.Command, args []string) {
			source := shell.ShellSource()
			fmt.Println(source)
		},
	}

	return cmd
}
