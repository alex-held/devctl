package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bndr/gotabulate"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	config2 "github.com/alex-held/devctl/internal/config"
	"github.com/alex-held/devctl/pkg/constants"
	"github.com/alex-held/devctl/pkg/devctlpath"

	"github.com/alex-held/devctl/internal/cli"
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
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config2.LoadViperConfig()

			var tVals [][]interface{}
			for _, sdk := range cfg.Sdks {
				currentVal := sdk.Current
				if currentVal == "" {
					currentVal = "<not installed>"
				}
				tRowHeader := []interface{}{sdk, currentVal, "", ""}
				tVals = append(tVals, tRowHeader)
				for _, sdkInstallationConfig := range sdk.Candidates {
					tRow := []interface{}{
						"",
						"",
						sdkInstallationConfig.Version,
						sdkInstallationConfig.Path,
					}
					tVals = append(tVals, tRow)
				}
			}
			t := gotabulate.Create(tVals)
			t.SetHeaders([]string{"sdk", "current", "version", "path"})
			t.SetEmptyString(" ")

			fmt.Println(t.Render("simple"))
		},
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
	versions := listSdkVersions()
	for _, version := range versions {
		fmt.Println(version)
	}
}

func listSdkVersions() (versions []string) {
	cfg := config2.LoadViperConfig()

	for _, sdkConfig := range cfg.Sdks {
		for _, sdkInstallation := range sdkConfig.Candidates {
			versions = append(versions, sdkInstallation.Version)
		}
	}
	return versions
}

var sdkCurrentPath *string
var sdkName *string
var sdkCandidatePath *string
var sdkVersion *string
var sdkCandidates []string

// newSdkAddCommand creates the `devenv sdk add` command
func newSdkAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [sdk]",
		Short: "Adds a local SDK",
		Args:  cobra.ExactArgs(1),
		Run:   sdkAddCommandfunc,
	}

//	cmd.Flags().StringVarP(sdkCurrentPath, "current", "c", "", "--current '~/.devctl/sdks/java/adopt-jdk/16.1'")
//	cmd.Flags().StringVarP(sdkName, "name", "n", "", "--name 'java'")
//	cmd.Flags().StringVarP(sdkVersion, "version", "v", "", "--version '1.0.0'")
//	cmd.Flags().StringVarP(sdkCandidatePath, "path", "p", "", "--path '~/.devctl/sdks/java/adopt-jdk/16.1'")
//	cmd.Flags().StringArrayVarP(&sdkCandidates, "candidates", "c", []string{}, "--candidates ~/.devctl/sdks/java/adopt-jdk/16.1'")

	return cmd
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
		cli.ExitWithError(1, fmt.Errorf("too many arguments for command '%s'. ", cmd.UsageTemplate()))
		return
	}

	removeSDK := args[0]
	devEnvConfig := config2.LoadViperConfig()

	filteredSdks := devEnvConfig.Sdks

	for sdk, config := range devEnvConfig.Sdks {
		if sdk != removeSDK {
			filteredSdks[sdk] = config
		}
	}
	devEnvConfig.Sdks = filteredSdks
	err := config2.UpdateDevEnvConfig(*devEnvConfig)
	if err != nil {
		cli.ExitWithError(1, err)
	}
}

// nolint: gocognit
func sdkAddCommandfunc(cmd *cobra.Command, args []string) {
	if len(args) > 1 {
		cli.ExitWithError(1, fmt.Errorf("too many arguments for command '%s'. ", cmd.UsageTemplate()))
		return
	}

	addSDK := args[0]

	cfg := config2.LoadViperConfig()

	// Determine whether or not the sdk is already already tracked
	if sdk, ok := cfg.Sdks[addSDK]; ok { //nolint:nestif
		// Add sdk-candidate path to sdk
		sdk.Candidates = append(sdk.Candidates, config2.SDKCandidate{
			Path:    *sdkCandidatePath,
			Version: *sdkVersion,
		})
		cfg.Sdks[addSDK] = sdk

		err := config2.WriteDevEnvConfig("$HOME/.devctl2/ccc33.yaml", *cfg)
		if err != nil {
			err = errors.Wrapf(err, "failed to save config. sdk not added.\n")
			cli.ExitWithError(constants.Failure, err)
		}

		// Quit
		os.Exit(0)
	} else {
		// Add new sdk
		newSDKConfig := config2.DevEnvSDKConfig{
			Current:    "latest",
			Candidates: []config2.SDKCandidate{},
		}

		*sdkName = devctlpath.SDKsPath(addSDK)
		sdkCandidateDirs, err := ioutil.ReadDir(*sdkName)
		if err != nil {
			panic(err)
		}

		for _, candidateDir := range sdkCandidateDirs {
			if candidateDir.IsDir() {
				name := candidateDir.Name()
				newSDKConfig.Candidates = append(newSDKConfig.Candidates, config2.SDKCandidate{
					Path:    name,
					Version: filepath.Dir(name),
				})
			}
			fmt.Printf("Candidate '%s' is not a directory.", candidateDir.Name())
		}

		cfg.Sdks[addSDK] = newSDKConfig
		err = config2.WriteDevEnvConfig("$HOME/.devctl2/ccc.yaml", *cfg)
		if err != nil {
			err = errors.Wrapf(err, "failed to save config. sdk not added.\n")
			cli.ExitWithError(constants.Failure, err)
		}
	}
}

func sdkListCommandfunc(cmd *cobra.Command, args []string) {
	sdks := listSdks()
	for _, sdk := range sdks {
		fmt.Println(sdk)
	}
}

func listSdks() (sdks []string) {
	devenv := config2.LoadViperConfig()
	for _, sdk := range devenv.Sdks {
		sdks = append(sdks, sdk.Current)
	}
	return sdks
}
