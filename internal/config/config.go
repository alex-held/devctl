package config

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"

	"github.com/alex-held/devctl/internal/devctlpath"
)

var (
	ErrSDKNotInConfig = errors.New("The config.SDKConfig for the given SDK-key is not available inside the config.Config")
)

const (
	VersionV1 = "v1"
)

type SdkConfig struct {
	SDK           string            `yaml:"sdk"`
	Current       string            `yaml:"current"`
	Installations map[string]string `yaml:"installations"`
}

func NewSdkConfig(sdk string) SdkConfig {
	return SdkConfig{SDK: sdk}
}

func (c *SdkConfig) GetInstallation(version string) (path string, err error) {
	if val, ok := c.Installations[version]; ok {
		return val, nil
	}
	return "", fmt.Errorf("installation")
}

type SdksConfig map[string]SdkConfig

type Config struct {
	Version string     `yaml:"version"`
	Sdks    SdksConfig `yaml:"sdks"`
}

func NewBlankConfig() (c *Config) {
	c = &Config{
		Version: VersionV1,
		Sdks:    SdksConfig{},
	}
	return c
}

func (c *Config) Apply(updateFn func(*Config)) *Config {
	updateFn(c)
	return c
}

func Save(fs afero.Fs, pather devctlpath.Pather, cfg *Config) (err error) {
	cfgPath := pather.ConfigFilePath()
	configBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal config.Config; config=%+v", *cfg)
	}

	err = afero.WriteFile(fs, cfgPath, configBytes, 0777)
	if err != nil {
		return errors.Wrapf(err, "failed to write config file to fs; path=%s; yaml=%s", cfgPath, string(configBytes))
	}

	return nil
}

func LoadOrCreate(fs afero.Fs, pather devctlpath.Pather) (cfg *Config, err error) {
	cfgPath := pather.ConfigFilePath()
	exists, err := afero.Exists(fs, cfgPath)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to check whether the config file exists on fs; path=%s", cfgPath)
	} else if exists {
		return Load(fs, pather)
	} else {
		_, err = fs.Create(cfgPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create empty file on fs; path=%s", cfgPath)
		}
		return Load(fs, pather)
	}
}

func Load(fs afero.Fs, pather devctlpath.Pather) (cfg *Config, err error) {
	cfgPath := pather.ConfigFilePath()
	configBytes, err := afero.ReadFile(fs, cfgPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read config file from fs; path=%s", cfgPath)
	}
	cfg = NewBlankConfig()
	err = yaml.Unmarshal(configBytes, cfg)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to unmarshal config.Config from config file; path=%s; yaml=%s", cfgPath, string(configBytes))
	}
	return cfg, nil
}
