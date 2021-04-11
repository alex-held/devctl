package testutils

import (
	"fmt"
	"os"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

//BeASymlink succeeds if a file exists and is a directory.
//Actual must be a string representing the abs path to the file being checked.
func BeASymlink(fs vfs.VFS) types.GomegaMatcher {
	return &BeASymlinkMatcher{
		Fs: fs,
	}
}

type notASymlinkError struct {
	VFS      vfs.VFS
	FileInfo os.FileInfo
	Err      error
}

func (t notASymlinkError) Error() string {
	fi := t.FileInfo
	mode := fi.Mode()
	return fmt.Sprintf("file mode is: %v", mode)
}

type BeASymlinkMatcher struct {
	Fs  vfs.VFS
	err error
}

func (matcher *BeASymlinkMatcher) Match(actual interface{}) (success bool, err error) {
	fileName, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("BeASymlinkMatcher matcher expects a file path")
	}

	fi, err := matcher.Fs.Lstat(fileName)
	if err != nil {
		matcher.err = err
		return false, nil
	}

	switch mode := fi.Mode(); {
	case mode&os.ModeSymlink != 0:
		return true, nil
	case mode.IsRegular():
	case mode.IsDir():
	default:
		matcher.err = notASymlinkError{
			VFS:      matcher.Fs,
			FileInfo: fi,
			Err:      fmt.Errorf("the file has a wrong os.FileMode. %v", mode),
		}
		return false, nil
	}
	return true, nil
}

func (matcher *BeASymlinkMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to be a symlink: %+v\n", matcher.err))
}

func (matcher *BeASymlinkMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("not be a symlink: %v\n", matcher.err))
}
