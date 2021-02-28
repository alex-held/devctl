package action

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Use action

var (
	ErrSdkVersionNotInstalled error = errors.New("SDK with version is not installed")
)

func (u *Use) Use(sdk, version string) (path string, err error) {
	path = u.Pather.SDK(sdk, version)
	exists, err := afero.Exists(u.Fs, path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to check whether sdk dir/file exists; path=%s", path)
	}
	if !exists {
		return "", errors.Wrapf(ErrSdkVersionNotInstalled, "sdk=%s; version=%s path=%s", sdk, version, path)
	}

	err = u.Actions.Config.SetCurrentSdk(sdk, version, path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to set current version for sdk; sdk=%s; version=%s", sdk, version)
	}
	return path, nil
}
