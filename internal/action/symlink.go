package action

import (
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Symlink action

func (s *Symlink) LinkCurrentSDK(sdk, version string) (current string, err error) {
	currentPath := s.Options.Pather.SDK(sdk, "current")
	versionPath := s.Options.Pather.SDK(sdk, version)
	exists, err := afero.Exists(s.Options.Fs, currentPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed check whether current symlink already exists; sdk=%s; path=%s", sdk, currentPath)
	} else if exists {
		err = s.Options.Fs.Remove(currentPath)
		s.Options.Logger.Debugf("removing symlink %s", currentPath)
		if err != nil {
			return "", errors.Wrapf(err, "failed to remove current symlink; path=%s", currentPath)
		}
	}

	cmd := exec.Command("ln", "-Ffsn", versionPath, currentPath)
	s.Options.Logger.Debugf("Command=%s\n", cmd.String())

	err = cmd.Run()
	if err != nil {
		return "", errors.Wrapf(err, "failed to exectute symlink command; command=%s; sdk=%s", cmd.String(), sdk)
	}

	return currentPath, nil
}
