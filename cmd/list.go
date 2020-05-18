/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/alex-held/dev-env/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"strings"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := readOrCreateConfig()
		executeList(*config, args)
	},
}

func readOrCreateConfig() *config.Config {
	fs := afero.NewOsFs()
	filepath := "/Users/dev/.dev-env/config.json"
	printError := func(e error) {
		_ = fmt.Errorf("Could not read config file from %s\n%s", filepath, e.Error())
	}
	ensureNoError := func(config *config.Config, e error) *config.Config {
		if e != nil {
			printError(e)
			return nil
		}
		return config
	}

	exists, err := afero.Exists(fs, filepath)
	if err != nil {
		printError(err)
		return nil
	}
	if exists {
		config, err := config.ReadConfigFromFile(fs, filepath)
		return ensureNoError(config, err)
	}

	config := config.NewConfig(fs, filepath)
	err = config.Save()
	return ensureNoError(config, err)
}

func executeList(config config.Config, args []string) []config.SDK {
	if len(args) == 0 {
		result := config.ListSdks()
		prettyPrintSdkTable(result)
		return result
	}

	result := []config.SDK{}

	for _, arg := range args {

		for _, sdk := range config.ListMatchingSdks(func(sdk config.SDK) bool {
			return sdk.Name == arg
		}) {
			result = append(result, sdk)
		}
	}

	prettyPrintSdkTable(result)
	return result

}

func prettyPrintSdkTable(sdks []config.SDK) {
	println(formatSdkTable(sdks))
}

func formatSdkTable(sdks []config.SDK) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("\n|%10s|%10s|%20s|\n", "sdk", "version", "path"))
	for _, sdk := range sdks {
		sb.WriteString(fmt.Sprintf("|%10s|%10s|%20s|\n", sdk.Name, sdk.Version, sdk.Path))
	}
	return sb.String()
}
