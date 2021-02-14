package cmd

import (
	"os"
	
	"github.com/alex-held/devctl/internal/sdkman"
	"github.com/spf13/cobra"
	
	"github.com/alex-held/devctl/pkg/cli"
)

type OutputFormat string

const (
	Text  OutputFormat = "text"
	Table OutputFormat = "table"
)

// NewSdkManCommand creates the sdkman command
func NewSdkManCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sdkman",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			// c := cli.GetOrCreateCLI()
			ctx :=	cmd.Context()
			client := sdkman.NewSdkManClient()
			sdks, resp, err := client.ListSdks.ListAllSDK(ctx)
			if err != nil {
				cli.ExitWithError(1, err)
			}
			defer resp.Body.Close()
			
			outputFlag := cmd.Flag("format")
			switch OutputFormat(outputFlag.Value.String()) {
			case Table:
				println(sdks.String())
				goto exit
			case Text:
				for _, sdk := range sdks {
					println(sdk)
				}
				goto exit
			default:
				println(sdks)
				goto exit
			}
			exit:
			os.Exit(0)
		},
	}
}
