package api

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type (
	ConfigFileSemVersion struct {
		Major uint64    `yaml:"major"`
		Minor uint64    `yaml:"minor"`
		Patch uint64    `yaml:"patch"`
		Pre   []*string `yaml:"pre,omitempty"`
		Build []*string `yaml:"build,omitempty"`
	}

	ConfigFileVersion struct {
		ID      string               `yaml:"id"`
		Version ConfigFileSemVersion `yaml:"version"`
		Vendor  string               `yaml:"vendor,omitempty"`
		Path    string               `yaml:"path"`
	}

	ConfigFile struct {
		Sdks     []*Sdk               `yaml:"sdks"`
		Contexts map[string]*Context  `yaml:"contexts"`
		Versions []*ConfigFileVersion `yaml:"versions"`
	}
)

func (config Config) toConfigFileYaml() string {
	configFile := NewConfigFile()
	configFile.Sdks = config.Sdks
	configFile.Contexts = config.Contexts

	for _, version := range config.Versions {

		configFileSemVer := ConfigFileSemVersion{
			Major: version.Version.Major,
			Minor: version.Version.Minor,
			Patch: version.Version.Patch,
			Pre:   nil,
			Build: nil,
		}

		configFileVersion := ConfigFileVersion{
			ID:      version.Id,
			Version: configFileSemVer,
			Vendor:  version.Vendor,
			Path:    version.Path,
		}

		configFile.Versions = append(configFile.Versions, &configFileVersion)
	}

	yamlBytes, err := yaml.Marshal(configFile)

	if err != nil {
		return "<yaml|nil>"
	}

	yamlString := string(yamlBytes)

	return yamlString
}

func toYaml(config ConfigFile) string {
	yamlBytes, err := yaml.Marshal(&config)

	if err != nil {
		return "<nil>"
	}

	yaml := string(yamlBytes)
	return yaml
}

func NewConfigFile() *ConfigFile {
	return &ConfigFile{
		Sdks:     []*Sdk{},
		Contexts: map[string]*Context{},
		Versions: []*ConfigFileVersion{},
	}
}

func readYaml(yml string) ConfigFile {
	config := ConfigFile{}

	err := yaml.Unmarshal([]byte(yml), config)

	if err != nil {
		return *NewConfigFile()
	}

	return config
}

func LoadConfig() (*Config, error) {
	configPath := "/Users/dev/.dev-env/config.yaml"
	fmt.Printf("Loading Configuration from %v", configPath)

	yamlBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Errorf("Error reading Configuration File '%v'!\n\n%v", configPath, err.Error())
		return nil, err
	}

	yaml := string(yamlBytes)
	config, parseError := parseYamlFromConfigFile(yaml)

	if parseError != nil {
		fmt.Errorf("Error parsing yaml into Config!\n\n%v\n\nError:\n%v", yaml, err.Error())
		return nil, err
	}
	return config, nil
}
