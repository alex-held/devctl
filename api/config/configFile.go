package api

import (
	"gopkg.in/yaml.v2"
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
