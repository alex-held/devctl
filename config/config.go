package config

import (
	"encoding/json"
	"github.com/spf13/afero"
	"sort"
)

var _ = afero.NewOsFs()

type Config struct {
	filepath string
	fs       afero.Fs
	Sdks     []SDK
}

type SDK struct {
	Name    string
	Version string
	Path    string
}

type SDKConfig interface {
	AddSDK(name string, version string, path string) *Config
}

type EnvVar struct {
	Key   string
	Value string
}

func NewConfig(fs afero.Fs, path string) *Config {
	config := &Config{
		filepath: path,
		fs:       fs,
		Sdks:     nil,
	}
	return config
}

func (config *Config) AddSDK(name string, version string, path string) *Config {
	config.Sdks = append(config.Sdks, SDK{
		Name:    name,
		Version: version,
		Path:    path,
	})
	config.Save()
	return config
}

func (config *Config) Save() error {
	return config.WriteToFile(config.fs, config.filepath)
}

func (config *Config) WriteToFile(fs afero.Fs, path string) error {
	file, marshalErr := json.MarshalIndent(config, "", " ")
	osErr := afero.WriteFile(fs, path, file, 0644)

	if marshalErr != nil {
		return marshalErr
	}
	if osErr != nil {
		return osErr
	}

	return nil
}

func ReadConfigFromFile(fs afero.Fs, path string) (*Config, error) {

	file, ioErr := afero.ReadFile(fs, path)

	if ioErr != nil {
		return nil, ioErr
	}

	config := NewConfig(fs, path)
	marshalErr := json.Unmarshal(file, &config)

	if marshalErr != nil {
		return nil, marshalErr
	}

	return config, nil
}

func (config *Config) ListSdks() []SDK {
	sort.Slice(config.Sdks, func(i, j int) bool {
		return config.Sdks[i].Name < config.Sdks[j].Name
	})
	return config.Sdks
}

func (config *Config) ListMatchingSdks(matcher func(sdk SDK) bool) []SDK {
	var result []SDK
	for _, sdk := range config.Sdks {
		if matcher(sdk) {
			result = append(result, sdk)
		}
	}
	return result
}
