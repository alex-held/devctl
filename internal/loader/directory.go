package loader

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/alex-held/devctl/internal/meta"

	"github.com/pkg/errors"
)

var utf8bom = []byte{0xEF, 0xBB, 0xBF}

// DirLoader loads a chart from a directory
type DirLoader string

// Load loads the chart
func (l DirLoader) Load() (*meta.Meta, error) {
	return LoadDir(string(l))
}

// LoadDir loads from a directory.
//
// This loads charts only from directories.
func LoadDir(dir string) (*meta.Meta, error) {
	topdir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	// Just used for errors.
	c := &meta.Meta{}

	files := []*BufferedFile{}
	topdir += string(filepath.Separator)

	walk := func(name string, fi os.FileInfo, err error) error {
		n := strings.TrimPrefix(name, topdir)
		if n == "" {
			// No need to process top level. Avoid bug with helmignore .* matching
			// empty names. See issue 1779.
			return nil
		}

		// Normalize to / since it will also work on Windows
		n = filepath.ToSlash(n)

		if err != nil {
			return err
		}
		if fi.IsDir() {
			// Directory-based ignore rules should involve skipping the entire
			// contents of that directory.
			return nil
		}

		// Irregular files include devices, sockets, and other uses of files that
		// are not regular files. In Go they have a file mode type bit set.
		// See https://golang.org/pkg/os/#FileMode for examples.
		if !fi.Mode().IsRegular() {
			return fmt.Errorf("cannot load irregular file %s as it has file mode type bits set", name)
		}

		data, err := ioutil.ReadFile(name)
		if err != nil {
			return errors.Wrapf(err, "error reading %s", n)
		}

		data = bytes.TrimPrefix(data, utf8bom)

		files = append(files, &BufferedFile{Name: n, Data: data})
		return nil
	}
	if err = Walk(topdir, walk); err != nil {
		return c, err
	}

	return LoadFiles(files)
}

// Walk walks the file tree rooted at root, calling walkFn for each file or directory
// in the tree, including root. All errors that arise visiting files and directories
// are filtered by walkFn. The files are walked in lexical order, which makes the
// output deterministic but means that for very large directories Walk can be
// inefficient. Walk follows symbolic links.
func Walk(root string, walkFn filepath.WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = symwalk(root, info, walkFn)
	}
	if err == filepath.SkipDir {
		return nil
	}
	return err
}

// readDirNames reads the directory named by dirname and returns
// a sorted list of directory entries.
func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

// symwalk recursively descends path, calling walkFn.
func symwalk(path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	// Recursively walk symlinked directories.
	if IsSymlink(info) {
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			return errors.Wrapf(err, "error evaluating symlink %s", path)
		}
		log.Printf("found symbolic link in path: %s resolves to %s", path, resolved)
		if info, err = os.Lstat(resolved); err != nil {
			return err
		}
		if err := symwalk(path, info, walkFn); err != nil && err != filepath.SkipDir {
			return err
		}
		return nil
	}

	if err := walkFn(path, info, nil); err != nil {
		return err
	}

	if !info.IsDir() {
		return nil
	}

	names, err := readDirNames(path)
	if err != nil {
		return walkFn(path, info, err)
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := os.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = symwalk(filename, fileInfo, walkFn)
			if err != nil {
				if (!fileInfo.IsDir() && !IsSymlink(fileInfo)) || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

// IsSymlink is used to determine if the fileinfo is a symbolic link.
func IsSymlink(fi os.FileInfo) bool {
	return fi.Mode()&os.ModeSymlink != 0
}
