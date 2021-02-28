package action

import (
	"github.com/pkg/errors"

	config2 "github.com/alex-held/devctl/internal/config"
)

type Config action

func (c *Config) Save(cfg *config2.Config) (err error) {
	err = config2.Save(c.Fs, c.Pather, cfg)
	if err != nil {
		return errors.Wrapf(err, "failed to save config file; config=%v\n", *cfg)
	}
	return nil
}

func (c *Config) Load() (cfg *config2.Config, err error) {
	cfg, err = config2.LoadOrCreate(c.Fs, c.Pather)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load config file")
	}
	return cfg, nil
}

func (c *Config) SetCurrentSdk(sdk, version, path string) (err error) {
	cfg, err := c.Load()
	if err != nil {
		return errors.Wrapf(err, "failed to load config.Config")
	}

	cfg.Apply(func(c *config2.Config) {
		if sdksVal, ok := c.Sdks[sdk]; ok {
			sdksVal.Current = version
			sdksVal.Installations[version] = path
			c.Sdks[sdk] = sdksVal
			return
		}

		c.Sdks[sdk] = config2.SdkConfig{
			SDK:     sdk,
			Current: version,
			Installations: map[string]string{
				version: path,
			},
		}
	})

	err = c.Save(cfg)
	if err != nil {
		return errors.Wrapf(err, "failed to save config file; config=%+v", *cfg)
	}
	return nil
}
