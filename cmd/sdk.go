package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
	
	"github.com/alex-held/dev-env/config"
)

// NewSdkCommand creates the `devenv sdk` commands
func NewSdkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "sdk",
		Short:     "Configure SDK's",
		ValidArgs: []string{"list", "add", "remove"},
	}
	
	cmd.AddCommand(newSdkListCommand())
	cmd.AddCommand(newSdkAddCommand())
	cmd.AddCommand(newSdkRemoveCommand())
	cmd.AddCommand(newSdkVersionsCommand())
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

func newSdkVersionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "version",
		Short:     "Configures sdk versions",
		ValidArgs: []string{"list"},
	}
	
	cmd.AddCommand(
		newSdkVersionsListCommand(),
	)
	
	return cmd
}

// newSdkVersionsListCommand creates the `devenv sdk list` command
func newSdkVersionsListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list [sdk]",
		Aliases: []string{"ls"},
		Short:   "Lists all available versions for [sdk]",
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return listSdks(), cobra.ShellCompDirectiveDefault
		},
		Run: sdkVersionsListCommandfunc,
	}
}

func sdkVersionsListCommandfunc(cmd *cobra.Command, args []string) {
	sdkArg := args[0]
	
	versions := listSdkVersions(sdkArg)
	for _, version := range versions {
		fmt.Println(version)
	}
}

func listSdkVersions(sdk string) (versions []string) {
	cfg := config.LoadViperConfig()
	
	for _, sdkConfig := range cfg.SDKConfig.SDKS {
		if sdkConfig.SDK == sdk {
			for _, sdkInstallation := range sdkConfig.Installations {
				versions = append(versions, sdkInstallation.Version)
			}
		}
	}
	
	return versions
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

// newSdkRemoveCommand creates the `devenv sdk remove` command
func newSdkRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "remove [sdk]",
		Aliases: []string{"rm"},
		Short:   "Adds a local SDK",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return listSdks(), 0
			
		},
		Args: cobra.ExactArgs(1),
		Run:  sdkRemoveCommandfunc,
	}
}

func sdkRemoveCommandfunc(cmd *cobra.Command, args []string) {
	if len(args) > 1 {
		ExitWithError(1, fmt.Errorf("Too many arguments for command '%s'.", cmd.UsageTemplate()))
		return
	}
	
	removeSDK := args[0]
	devEnvConfig := config.LoadViperConfig()
	
	filteredSdks := devEnvConfig.SDKConfig.SDKS[:0]
	for _, sdkConfig := range devEnvConfig.SDKConfig.SDKS {
		if sdkConfig.SDK != removeSDK {
			filteredSdks = append(filteredSdks, sdkConfig)
		}
	}
	devEnvConfig.SDKConfig.SDKS = filteredSdks
	
	err := config.UpdateDevEnvConfig(*devEnvConfig)
	if err != nil {
		ExitWithError(1, err)
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
	
	sdks := listSdks()
	for _, sdk := range sdks {
		fmt.Println(sdk)
	}
}

func listSdks() (sdks []string) {
	devenv := config.LoadViperConfig()
	for _, sdk := range devenv.SDKConfig.SDKS {
		sdks = append(sdks, sdk.SDK)
	}
	return sdks
}
