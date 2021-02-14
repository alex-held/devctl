package cmd

import (
	"context"
	"fmt"
	
	"github.com/alex-held/devctl/internal/sdkman"
	"github.com/spf13/cobra"
	
	"github.com/alex-held/devctl/pkg/cli"
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
			ctx := context.Background()
			client := sdkman.NewSdkManClient()
			sdks, resp, err := client.ListSdks.ListAllSDK(ctx)
			if err != nil {
				cli.ExitWithError(1, err)
			}
			defer resp.Body.Close()
			fmt.Printf("%#v", sdks)
		},
	}
}
