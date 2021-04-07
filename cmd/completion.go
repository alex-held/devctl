// Package cmd
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/internal/app"
)

// NewCompletionCommand Creates new Completion cobra.Command
func NewCompletionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(yourprogram completion bash)

# To load completions for each session, execute once:
Linux:
  $ yourprogram completion bash > /etc/bash_completion.d/yourprogram
MacOS:
  $ yourprogram completion bash > /usr/local/etc/bash_completion.d/yourprogram

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ dev-env completion zsh > "${fpath[1]}/_dev-env"

# You will need to start a new shell for this setup to take effect.

Powershell:

PS> yourprogram completion powershell | Out-String | Invoke-Expression

# To load completions for every new session, run:
PS> yourprogram completion powershell > yourprogram.ps1
# and source this file from your powershell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				err := cmd.Root().GenBashCompletion(os.Stdout)
				app.ExitWhenError(1, err)
			case "zsh":
				err := cmd.Root().GenZshCompletion(os.Stdout)
				app.ExitWhenError(1, err)
			case "powershell":
				err := cmd.Root().GenPowerShellCompletion(os.Stdout)
				app.ExitWhenError(1, err)
			}
		},
	}
}
