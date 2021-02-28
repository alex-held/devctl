package action

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/alex-held/devctl/internal/config"
)

type Use action

func (u Use) UseSDK(sdk, version string) (err error) {
	cfg, err := config.Load(u.Options.Fs, u.Options.Pather)
	if err != nil {
		return errors.Wrap(err, "failed to load config from file")
	}

	if err = validateUseSDK(cfg, sdk, version); err != nil {
		return errors.Wrapf(err, "failed to validate UseSDK inputs; sdk=%s; version=%s; config=%+v", sdk, version, *cfg)
	}

	cfg = cfg.Apply(func(c *config.Config) {
		sdkCfg := c.Sdks[sdk]
		c.Sdks[sdk] = config.SdkConfig{
			SDK:           sdkCfg.SDK,
			Current:       version,
			Installations: sdkCfg.Installations,
		}
	})
	err = u.Config.Save(cfg)
	if err != nil {
		return errors.Wrapf(err, "failed to save config to fs")
	}

	_, err = u.Symlink.LinkCurrentSDK(sdk, version)
	if err != nil {
		return errors.Wrapf(err, "failed to link the extracted sdk to current; sdk=%s; version=%s", sdk, version)
	}

	return nil
}

func validateUseSDK(cfg *config.Config, sdk, version string) (err error) {
	if valSdks, ok := cfg.Sdks[sdk]; ok {
		if _, ok = valSdks.Installations[version]; ok {
			return nil
		}
		return fmt.Errorf("invalid version; reason=version not installed")
	}
	return fmt.Errorf("invalid sdk; reason=sdk not installed")
}
