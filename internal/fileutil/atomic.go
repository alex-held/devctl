package fileutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func AtomicWriteFile(fs afero.Fs, filename string, r *bytes.Reader, mode os.FileMode) error {
	dir, file := filepath.Split(filename)

	tempFile, err := afero.TempFile(fs, dir, file)
	if err != nil {
		return err
	}
	tempName := tempFile.Name()

	if _, err := io.Copy(tempFile, r); err != nil {
		tempFile.Close() // return value is ignored as we are already on error path
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	if err := fs.Chmod(tempName, mode); err != nil {
		return err
	}

	return RenameWithFallback(fs, tempName, filename)
}

// RenameWithFallback attempts to rename a file or directory, but falls back to
// copying in the event of a cross-device link error. If the fallback copy
// succeeds, src is still removed, emulating normal rename behavior.
func RenameWithFallback(fs afero.Fs, src, dst string) error {
	_, err := fs.Stat(src)
	if err != nil {
		return errors.Wrapf(err, "cannot stat %s", src)
	}

	err = fs.Rename(src, dst)
	if err == nil {
		return nil
	}

	return renameFallback(fs, err, src, dst)
}

// renameFallback attempts to determine the appropriate fallback to failed rename
// operation depending on the resulting error.
func renameFallback(fs afero.Fs, err error, src, dst string) error {
	// Rename may fail if src and dst are on different devices; fall back to
	// copy if we detect that case. syscall.EXDEV is the common name for the
	// cross device link error which has varying output text across different
	// operating systems.
	terr, ok := err.(*os.LinkError)
	if !ok {
		return err
	} else if terr.Err != syscall.EXDEV {
		return errors.Wrapf(terr, "link error: cannot rename %s to %s", src, dst)
	}

	return renameByCopy(fs, src, dst)
}

func CopyDir(fs afero.Fs, src, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = fs.Stat(src); err != nil {
		return err
	}

	if err = fs.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = afero.ReadDir(fs, src); err != nil {
		return err
	}

	for _, fd := range fds {
		srcfd := path.Join(src, fd.Name())
		dstfd := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(fs, srcfd, dstfd); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(fs, srcfd, dstfd); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func CopyFile(fs afero.Fs, src string, dst string) error {
	var err error
	var srcfd afero.File
	var dstfd afero.File
	var srcinfo os.FileInfo

	if srcfd, err = fs.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = fs.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = fs.Stat(src); err != nil {
		return err
	}
	return fs.Chmod(dst, srcinfo.Mode())
}

// renameByCopy attempts to rename a file or directory by copying it to the
// destination and then removing the src thus emulating the rename behavior.
func renameByCopy(fs afero.Fs, src, dst string) error {
	var cerr error
	if dir, _ := afero.IsDir(fs, src); dir {
		cerr = CopyDir(fs, src, dst)
		if cerr != nil {
			cerr = errors.Wrap(cerr, "copying directory failed")
		}
	} else {
		cerr = CopyFile(fs, src, dst)
		if cerr != nil {
			cerr = errors.Wrap(cerr, "copying file failed")
		}
	}

	if cerr != nil {
		return errors.Wrapf(cerr, "rename fallback failed: cannot rename %s to %s", src, dst)
	}

	return errors.Wrapf(fs.RemoveAll(src), "cannot delete %s", src)
}
