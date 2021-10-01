package errors

import (
	"errors"
)

var ErrExit = errors.New("devctl has exited unsuccessfully")
var ErrWindowsNotSupported = errors.New("error: windows is not supported")
