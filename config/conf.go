package config

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type DevEnvGlobalConfig struct {
	Version string `yaml:"version,omitempty" mapstructure:"version,omitempty"`
}

type DevEnvSDKConfig struct {
	SDK           string                        `yaml:"sdk" json:"sdk" mapstructure:"sdk"`
	Current       string                        `yaml:"current,omitempty" mapstructure:"current,omitempty"`
	Installations []DevEnvSDKInstallationConfig `yaml:"installations,omitempty" mapstructure:"installations,omitempty"`
}

type DevEnvSDKInstallationConfig struct {
	Path    string `yaml:"path" mapstructure:"path"`
	Version string `yaml:"version" mapstructure:"version"`
}

type DevEnvSDKSConfig struct {
	SDKS []DevEnvSDKConfig `yaml:"sdks,omitempty" mapstructure:"sdks,omitempty"`
}

type DevEnvConfig struct {
	GlobalConfig DevEnvGlobalConfig `yaml:"global" mapstructure:"global"`
	SDKConfig    DevEnvSDKSConfig   `yaml:"sdk,omitempty" mapstructure:"sdk,omitempty"`
}

var (
	DefaultDevEnvConfig = DevEnvConfig{
		GlobalConfig: DefaultDevEnvGlobalConfig,
	}

	DefaultDevEnvGlobalConfig = DevEnvGlobalConfig{
		Version: "v1",
	}
)

// ResolveDevEnvConfigPath Resolves the current dev-env config file.
func ResolveDevEnvConfigPath() (err error, dir string) {
	if err, dir = ResolveDevEnvConfigDir(); err == nil {
		return nil, filepath.Join(dir, "config", "yaml")
	}
	return err, ""
}

// ResolveDevEnvConfigDir Resolves the current dev-env root directory.
func ResolveDevEnvConfigDir() (err error, dir string) {
	home, err := homedir.Dir()
	return err, path.Join(home, ".devctl")
}

func InitViper(filename string) {
	dir := path.Dir(filename)
	config := path.Base(filename)

	fmt.Printf("Config Directory: '%s'\n", dir)
	fmt.Printf("Config File: '%s'\n", config)
	viper.AddConfigPath(dir)
	viper.SetConfigName(config)
	viper.SetConfigType("yaml")
}

func LoadViperConfig() *DevEnvConfig {

	configuration := &DevEnvConfig{}

	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("Error reading config file, %s\n", err)
	}

	err := viper.Unmarshal(configuration)
	if err != nil {
		_ = fmt.Errorf("unable to decode into struct, %v\n", err)
	}
	return configuration
}

func UpdateDevEnvConfig(cfg DevEnvConfig) error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	devEnvConfigMap := &map[string]interface{}{}
	err = mapstructure.Decode(&cfg, devEnvConfigMap)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(devEnvConfigMap)
	err = viper.MergeConfig(bytes.NewReader(b))
	if err != nil {
		return err
	}

	return viper.WriteConfig()
}
