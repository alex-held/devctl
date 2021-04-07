package sdk

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

type Actions struct {
	Pather devctlpath.Pather
	FS     afero.Fs
}

func (a *Actions) List() (sdks []string, err error) {
	sdkDir := a.Pather.SDK()
	return GetSubdirNames(a.FS, sdkDir, NoOpExcluder)
}

var NoOpExcluder Excluder = func(s string) bool {
	return false
}

type Excluder func(s string) bool

func (a *Actions) ListVersions(sdk string) (sdks []string, err error) {
	sdkDir := a.Pather.SDK(sdk)
	if exists, _ := afero.DirExists(a.FS, sdkDir); !exists {
		return sdks, nil
	}
	return GetSubdirNames(a.FS, sdkDir, func(s string) bool {
		return s == "current"
	})
}

func GetSubdirNames(fs afero.Fs, dir string, excluder Excluder) (names []string, err error) {
	fis, err := afero.ReadDir(fs, dir)
	if err != nil {
		return names, errors.Wrapf(err, "failed to collect subdirs for %s", dir)
	}

	for _, fi := range fis {
		dirName := fi.Name()
		if !excluder(dirName) {
			names = append(names, dirName)
		}
	}

	return names, nil
}
