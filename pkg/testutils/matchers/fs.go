package matchers

import (
	"fmt"
	"os"

	"github.com/mandelsoft/vfs/pkg/vfs"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/spf13/afero"
)

func BeAnExistingFileFs(fs interface{}) OmegaMatcher {
	switch fs.(type) {
	case vfs.VFS:
		return &FsBeAnExistingFileMatcher{
			FS: fs,
		}
	case afero.Fs:
		return &FsBeAnExistingFileMatcher{
			FS: fs,
		}
	default:
		panic(fmt.Sprintf("[FsBeAnExistingFile] fs has type '%T', but it's not supported. \t%v", fs, fs))
	}

	return &FsBeAnExistingFileMatcher{
		FS: fs,
	}
}

type FsBeAnExistingFileMatcher struct {
	FS interface{}
}

func (matcher *FsBeAnExistingFileMatcher) Match(actual interface{}) (success bool, err error) {
	actualFilename, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("FsBeAnExistingFileMatcher matcher expects a file path")
	}

	var existFn func() error
	switch t := matcher.FS.(type) {
	case afero.Fs:
		existFn = func() error {
			_, err = t.Stat(actualFilename)
			return err
		}
	case vfs.VFS:
		existFn = func() error {
			_, err = t.Stat(actualFilename)
			return err
		}
	default:
		return false, fmt.Errorf("[FsBeAnExistingFile] fs has type '%T', but it's not supported. \t%v\n", t, t)
	}

	if err = existFn(); err != nil {
		switch {
		case os.IsNotExist(err):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (matcher *FsBeAnExistingFileMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to exist"))
}

func (matcher *FsBeAnExistingFileMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("not to exist"))
}
