package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
	
	"github.com/alex-held/dev-env/config"
)

// NewSdkCommand creates the `devenv sdk` commands
func NewSdkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sdk",
		Short: "Configure SDK's",
	}
	
	cmd.AddCommand(newSdkListCommand())
	cmd.AddCommand(newSdkAddCommand())
	return cmd
}

// newSdkListCommand creates the `devenv sdk list` command
func newSdkListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Lists all available SDK's",
		Run:     sdkListCommandfunc,
	}
}

// newSdkAddCommand creates the `devenv sdk add` command
func newSdkAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add [sdk]",
		Short: "Adds a local SDK",
		Args:  cobra.ExactArgs(1),
		Run:   sdkAddCommandfunc,
	}
}

func sdkAddCommandfunc(cmd *cobra.Command, args []string) {
	if len(args) > 1 {
		ExitWithError(1, fmt.Errorf("Too many arguments for command '%s'.", cmd.UsageTemplate()))
		return
	}
	
	addSDK := args[0]
	devEnvConfig := config.LoadViperConfig()
	for _, sdkConfig := range devEnvConfig.SDKConfig.SDKS {
		if sdkConfig.SDK == addSDK {
			ExitWithError(1, fmt.Errorf("SDK'%s' already configured.", addSDK))
			return
		}
	}
	
	devEnvConfig.SDKConfig.SDKS = append(devEnvConfig.SDKConfig.SDKS, config.DevEnvSDKConfig{
		SDK: addSDK,
	})
	
	err := config.UpdateDevEnvConfig(*devEnvConfig)
	if err != nil {
		ExitWithError(1, err)
	}
}



func sdkListCommandfunc(cmd *cobra.Command, args []string) {
	devenv := config.LoadViperConfig()
	var sdks []string
	for _, sdk := range devenv.SDKConfig.SDKS {
		sdks = append(sdks, sdk.SDK)
	}
	
	for _, sdk := range sdks {
		fmt.Println(sdk)
	}
}
