package testutils

import (
	"fmt"
	"os"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

//BeADirectoryFs succeeds if a file exists and is a directory.
//Actual must be a string representing the abs path to the file being checked.
func BeADirectoryFs(fs vfs.VFS) types.GomegaMatcher {
	return &BeADirMatcher{
		Fs: fs,
	}
}

type notADirError struct {
	VFS      vfs.VFS
	FileInfo os.FileInfo
	Err      error
}

func (t notADirError) Error() string {
	fi := t.FileInfo
	mode := fi.Mode()
	return fmt.Sprintf("file mode is: %v", mode)
}

type BeADirMatcher struct {
	Fs       vfs.VFS
	expected interface{}
	err      error
}

func (matcher *BeADirMatcher) Match(actual interface{}) (success bool, err error) {
	fileName, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("BeADirMatcher matcher expects a file path")
	}

	fi, err := matcher.Fs.Lstat(fileName)
	if err != nil {
		matcher.err = err
		return false, nil
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
	case mode&os.ModeDir != 0:
		return true, nil
	case mode.IsRegular():
	default:
		matcher.err = notADirError{
			VFS:      matcher.Fs,
			FileInfo: fi,
			Err:      fmt.Errorf("the file has a wrong os.FileMode. %v", mode),
		}
		return false, nil
	}
	return true, nil
}

func (matcher *BeADirMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to be a directory: %s", matcher.err))
}

func (matcher *BeADirMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("not be a directory"))
}
