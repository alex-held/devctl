package golang

import (
	"fmt"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/pkg/errors"

	"github.com/alex-held/devctl/pkg/system"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

func FormatGoArchiveArtifactName(ri system.RuntimeInfo, version string) string {
	return ri.Format("go%s.[os]-[arch].tar.gz", version)
}

func SymLink(pather devctlpath.Pather, fs vfs.VFS, version string) (err error) {
	sdkPath := pather.SDK("go", version)
	current := pather.SDK("go", "current")

	fs.Remove(current)
	fs.Symlink(sdkPath, current)

	readlink, err := fs.Readlink(current)
	if err == nil {
		err = fs.Remove(readlink)
		if err != nil {
			return errors.Wrapf(err, "failed to remove symlink %s", readlink)
		}
	}
	if err != nil && len(readlink) > 0 {
		return errors.Wrapf(err, "failed to remove symlink %s", readlink)
	}
	if readlink != "" {
		err = fs.Remove(readlink)
		if err != nil {
			return errors.Wrapf(err, "failed to remove symlink %s", readlink)
		}
	}
	fmt.Printf("there is no existing symlink that needs to be removed.")

	err = fs.Symlink(sdkPath, current)
	if err != nil {
		return errors.Wrapf(err, "failed to create symlink! \nsrc=%s\ndest=%s", sdkPath, current)
	}
	return err
}
