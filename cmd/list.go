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

	"github.com/rs/zerolog/log"

	"github.com/alex-held/dev-env/api"
	. "github.com/alex-held/dev-env/config"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
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
		client := api.NewGithubAPI(nil)
		sdks, err := client.GetPackages("sdk")
		if err != nil {
			panic(err)
		}
		log.Info().Strs("sdks", sdks).Send()
	},
}

func readOrCreateConfig() *Config {
	fs := afero.NewOsFs()
	filepath := "/Users/dev/.dev-env/config.json"
	printError := func(e error) {
		_ = fmt.Errorf("Could not read config file from %s\n%s", filepath, e.Error())
	}
	ensureNoError := func(config *Config, e error) *Config {
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
		config, err := ReadConfigFromFile(fs, filepath)
		return ensureNoError(config, err)
	}

	config := NewConfig(fs, filepath)
	err = config.Save()
	return ensureNoError(config, err)
}
