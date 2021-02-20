package cmd

import (
	"os"

	"github.com/alex-held/devctl/internal/sdkman"
	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/pkg/cli"
)

var fmtFlag string

type OutputFormat string

const (
	Text  OutputFormat = "text"
	Table OutputFormat = "table"
)

// NewSdkManCommand creates the sdkman command
func NewSdkManCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "sdkman",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	cmd.AddCommand(NewSdkManListCommand())
	return cmd
}

// NewSdkManCommand creates the sdkman command
func NewSdkManListCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "list",
		Short: "lists all installable sdks",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			client := sdkman.NewSdkManClient()
			sdks, resp, err := client.ListSdks.ListAllSDK(ctx)
			if err != nil {
				cli.ExitWithError(1, err)
			}
			defer resp.Body.Close()

			formatFlag := cmd.Flag("format")
			format := formatFlag.Value.String()

			switch OutputFormat(format) {
			case Table:
				println(sdks.String())
				os.Exit(0)
			case Text:
				for _, sdk := range sdks {
					println(sdk)
				}
				os.Exit(0)
			default:
				println(sdks)
				os.Exit(0)
			}
		},
	}

	cmd.Flags().StringVarP(&fmtFlag, "format", "f", string(Table), "the output format of the cli app. -format=table")
	return cmd
}
