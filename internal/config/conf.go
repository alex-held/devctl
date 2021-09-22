package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	yaml2 "sigs.k8s.io/yaml"

	"github.com/alex-held/devctl-kit/pkg/constants"
)

// DefaultConfig creates a blank *DevCtlConfig
func DefaultConfig() *DevCtlConfig {
	return &DevCtlConfig{
		GlobalConfig: DevEnvGlobalConfig{Version: constants.VersionV1},
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

type DevCtlConfig struct {
	GlobalConfig DevEnvGlobalConfig         `yaml:"global,omitempty" json:"global,omitempty" mapstructure:"global,omitempty"`
	Sdks         map[string]DevEnvSDKConfig `yaml:"sdks,omitempty" json:"sdks,omitempty" mapstructure:"sdks,omitempty"`
}

func (d *DevCtlConfig) GoString() string {
	b, e := yaml.Marshal(d)
	if e != nil {
		return e.Error()
	}
	return fmt.Sprintf("%+v", string(b))
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

func ParseConfigFile(filename string) (configuration *DevCtlConfig, err error) {
	return parseConfigFile(filename)
}

func parseConfigFile(filename string) (configuration *DevCtlConfig, err error) {
	data, _ := ReadConfigFile(filename)
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

func WriteDevEnvConfig(filepath string, cfg DevCtlConfig) error {
	data, err := yaml2.Marshal(&cfg)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal the DevCtlConfig. %+v\n", cfg)
	}

	err = ioutil.WriteFile(filepath, data, fileutil.PrivateFileMode)
	if err != nil {
		return errors.Wrapf(err, "failed to write config file to disk")
	}
	return nil
}
