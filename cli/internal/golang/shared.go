package golang

import (
	"fmt"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/pkg/errors"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

/*
type symlinkerFs struct {
	vfs      vfs.FileSystem
}

//LstatIfPossible
func (s *symlinkerFs) LstatIfPossible(name string) (stat os.FileInfo, lstatCalled bool, err error) {
	if symlinkDest, ok := s.symlinks[name]; ok {
		fmt.Printf("resolving symlink.\nsrc=%s;\ndest=%s\n", name, symlinkDest)
		lstatCalled = true
		stat, err = s.fs.Stat(symlinkDest)
		return stat, lstatCalled, err
	} else {
		lstatCalled = false
		stat, err = s.fs.Stat(name)
		return stat, false, err
	}
}

func (s *symlinkerFs) SymlinkIfPossible(oldname, newname string) (err error) {
	src, srcErr := s.fs.Stat(oldname)
	if src != nil {
		s.dumpFileInfo(src)
	}
	dest, destErr := s.fs.Stat(newname)
	if dest != nil {
		s.dumpFileInfo(dest)
	}

	if srcErr == nil {
		if destErr == nil {
			return errors.Wrapf(destErr, "SymlinkIfPossible failed! dest=%s already exists\n", newname)
		} else {
			if !src.IsDir() {
				var srcBytes, err = afero.ReadFile(s.fs, oldname)
				if err != nil {
					return errors.Wrapf(err, "SymlinkIfPossible failed! unable to read src file %s\n", oldname)
				}
				destFile, err := s.fs.Create(newname)
				if err != nil {
					return errors.Wrapf(err, "SymlinkIfPossible failed! unable to create dest file %s\n", newname)
				}
				n, err := destFile.Write(srcBytes)
				if err != nil {
					return errors.Wrapf(err, "SymlinkIfPossible failed! unable to write dest file %s; stopped at index %d \n", newname, n)
				}
				s.symlinks[oldname] = newname
				return nil
			} else {
				err = afero.Walk(s.fs, oldname, func(path string, info fs2.FileInfo, err error) error {
					if err != nil {
						return errors.Wrapf(err, "[%s] -> error while walking path %s, %v\n", oldname, path, info)
					}
					fmt.Printf("[%s] -> walking path %s, %v\n", oldname, path, info)
					return nil
				})
				if err != nil {
					return errors.Wrapf(err, "failed to walk dir %s", oldname)
				}
				return nil
			}
		}
	} else {
		return errors.Wrapf(srcErr, "SymlinkIfPossible failed! src=%s does not exists\n", oldname)
	}
}

func (s *symlinkerFs) dumpFileInfo(src os.FileInfo) {
	perm := src.Mode().Perm()
	mode := src.Mode()
	fmt.Printf("[SRC] perm: %s; filemode: %s", perm.String(), mode.String())
}

func (s *symlinkerFs) ReadlinkIfPossible(name string) (string, error) {
	panic("implement me")
}

func NewSymlinkerFs(fs afero.Fs) symlinkerFs {
	return symlinkerFs{
		symlinks: map[string]string{},
		fs:       fs,
	}
}
*/

func SymLink(pather devctlpath.Pather, fs vfs.VFS, version string) (err error) {
	sdkPath := pather.SDK("go", version)
	current := pather.SDK("go", "current")
	println(sdkPath)
	println(current)

	evalSymlinks, err := fs.EvalSymlinks(current)
	if err != nil {
		println("evalSymlinks failed: \n" + err.Error())
		return err
	}
	println(evalSymlinks)

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

/*

func Link(pather devctlpath.Pather, fs afero.Fs, version string) (err error) {
	sdkPath := pather.SDK("go", version)
	current := pather.SDK("go", "current")

	info, err := fs.Stat(current)
	if err != nil {
		fmt.Printf("no file system entry exists for path %s", current)
	} else if info.Mode().Type() == os.ModeSymlink {
		fmt.Printf("there is a symlink at path %s", current)
		_ = fs.Remove(current)
	}

	if ok, _ := afero.DirExists(fs, current); ok {
		_ = fs.Remove(current)
	}
	if ok, _ := afero.Exists(fs, current); ok {
		_ = fs.Remove(current)
	}

	symlink := exec.Command("ln", "-s", sdkPath, current)
	err = symlink.Run()
	if err != nil {
		return errors.Wrapf(err, "failed linking go sdk %s; src=%s; dest=%s\n", version, sdkPath, current)
	}
	return nil
}
*/
