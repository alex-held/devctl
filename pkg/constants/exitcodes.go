package constants

import (
	"fmt"
	"os"
)

type ExitCode int

type ExitError struct {
	ExitCode ExitCode
	Message  string
}

// Code returns the ExitCode as int of the ExitError
func (e ExitError) Code() int { return int(e.ExitCode) }
func (e ExitError) Error() string {
	return e.Message
}

// ExitWithError  prints an error message and exits the application with ErrorCode: code
func ExitWithError(code int, err error) {
	if err == nil {
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(code)
}

const (
	OK ExitCode = 0

	IssuesFound = 1

	WarningInTest = 2

	Failure = 3

	Timeout = 4

	NoGoFiles = 5

	NoConfigFileDetected = 6

	ErrorWasLogged = 7
)

var (
	// ErrNoGoFiles is the pre-defined ExitError NoGoFiles
	ErrNoGoFiles = &ExitError{
		Message:  "no go files to analyze",
		ExitCode: NoGoFiles,
	}
	// ErrFailure is the pre-defined ExitError Failure
	ErrFailure = &ExitError{
		Message:  "failed to analyze",
		ExitCode: Failure,
	}
)
