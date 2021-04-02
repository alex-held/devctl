package cmd

import (
	"github.com/spf13/cobra"
)

func (e *Executor) initCompletion() {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Output completion script",
	}
	e.rootCmd.AddCommand(completionCmd)
	zshCmd := &cobra.Command{
		Use:   "zsh",
		Short: "Output zsh completion script",
		RunE:  e.executeZshCompletion,
	}
	completionCmd.AddCommand(zshCmd)
}
