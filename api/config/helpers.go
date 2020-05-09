package api

import (
	"fmt"
	"github.com/blang/semver"
	yaml2 "gopkg.in/yaml.v2"
	"strings"
)

func isConfigEmpty(config *Config) bool {
	return len(config.Versions) == 0 && len(config.Contexts) == 0 && len(config.Sdks) == 0
}

func parseYamlFromConfigFile(yaml string) (*Config, error) {
	configFile := NewConfigFile()
	err := yaml2.Unmarshal([]byte(yaml), &configFile)

	for key, context := range configFile.Contexts {
		context.SdkId = key
	}

	config := NewConfig()
	config.Contexts = configFile.Contexts
	config.Sdks = configFile.Sdks

	for _, version := range configFile.Versions {
		var preVersions []semver.PRVersion

		if len(version.Version.Pre) > 0 {
			for _, s := range version.Version.Pre {
				pr, er := semver.NewPRVersion(*s)
				if er != nil {
					preVersions = nil
					break
				}
				preVersions = append(preVersions, pr)
			}
		} else {
			preVersions = nil
		}

		var build []string
		if len(version.Version.Build) > 0 {

			for _, s := range version.Version.Build {
				if s == nil {
					build = nil
					break
				}
				build = append(build, *s)
			}
		} else {
			build = nil
		}

		config.Versions = append(config.Versions, &Version{
			Id: version.ID,
			Version: semver.Version{
				Major: version.Version.Major,
				Minor: version.Version.Minor,
				Patch: version.Version.Patch,
				Pre:   preVersions,
				Build: build,
			},
			Vendor: version.Vendor,
			Path:   version.Path,
		})
	}

	return config, err
}

func (config *Config) toYaml() string {
	bytes, err := yaml2.Marshal(config)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}
	output := string(bytes)
	output = strings.TrimSpace(output)
	return output
}
