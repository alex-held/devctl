package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	yaml2 "sigs.k8s.io/yaml"

	"github.com/alex-held/devctl/internal/devctlpath"
)

// DefaultConfig creates a blank *DevEnvConfig
func DefaultConfig() *DevEnvConfig {
	return &DevEnvConfig{
		GlobalConfig: DevEnvGlobalConfig{Version: VersionV1},
		Sdks:         map[string]DevEnvSDKConfig{},
	}
}

type DevEnvGlobalConfig struct {
	Version string `yaml:"version,omitempty" json:"version,omitempty" mapstructure:"version,omitempty"`
}

type DevEnvSDKConfig struct {
	Current    string         `yaml:"current,omitempty" json:"current,omitempty" mapstructure:"current,omitempty"`
	Candidates []SDKCandidate `yaml:"candidates,omitempty"  json:"candidates,omitempty" mapstructure:"candidates,omitempty"`
}

type SDKCandidate struct {
	Path    string `yaml:"path,omitempty" json:"path,omitempty" mapstructure:"path,omitempty"`
	Version string `yaml:"version,omitempty" json:"version,omitempty" mapstructure:"version,omitempty"`
}

type DevEnvConfig struct {
	GlobalConfig DevEnvGlobalConfig         `yaml:"global,omitempty" json:"global,omitempty" mapstructure:"global,omitempty"`
	Sdks         map[string]DevEnvSDKConfig `yaml:"sdks,omitempty" json:"sdks,omitempty" mapstructure:"sdks,omitempty"`
}

func (d *DevEnvConfig) GoString() string {
	b, e := yaml.Marshal(d)
	if e != nil {
		return e.Error()
	}
	return fmt.Sprintf("%+v", string(b))
}

var initialized bool = false

func InitViper(filename string) {
	dir := path.Dir(filename)
	config := path.Base(filename)
	viper.AddConfigPath(dir)
	viper.SetConfigName(config)
	viper.SetConfigType("yaml")
	initialized = true
}

var ReadConfigFile = func(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open '%s'", filename)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseConfigFile(filename string) (configuration *DevEnvConfig, err error) {
	data, err := ReadConfigFile(filename)
	configuration = DefaultConfig()
	err = yaml.Unmarshal(data, configuration)

	if err != nil {
		return configuration, errors.Wrapf(err, "failed to read config file from disk")
	}

	c := DefaultConfig()
	err = yaml2.Unmarshal(data, &c)

	if err != nil {
		return configuration, errors.Wrapf(err, "failed to unmarshal yaml-file; yaml=%s", string(data))
	}
	return configuration, nil
}

func LoadViperConfig() *DevEnvConfig {
	if !initialized {
		cfgPath := devctlpath.DevCtlConfigFilePath()
		InitViper(cfgPath)
	}

	configuration := DefaultConfig()
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("error reading config file, %s\n ", err)
	}
	err := viper.Unmarshal(configuration)
	devEnvConfigMap := &map[string]interface{}{}
	err = mapstructure.Decode(&configuration, devEnvConfigMap)
	if err != nil {
		_ = fmt.Errorf("unable to decode into struct, %v\n ", err)
	}
	return configuration
}

func WriteDevEnvConfig(filepath string, cfg DevEnvConfig) error {
	data, err := yaml2.Marshal(&cfg)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal the DevEnvConfig. %+v\n", cfg)
	}

	err = ioutil.WriteFile(filepath, data, fileutil.PrivateFileMode)
	if err != nil {
		return errors.Wrapf(err, "failed to write config file to disk")
	}
	return nil
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

	b, _ := yaml.Marshal(devEnvConfigMap)
	err = viper.MergeConfig(bytes.NewReader(b))
	if err != nil {
		return err
	}

	return viper.WriteConfig()
}
