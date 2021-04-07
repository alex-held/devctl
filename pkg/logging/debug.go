package logging

import (
	"os"
	"strings"
)

type DebugFunc func(format string, args ...interface{})

var _ Log = NewLogger(WithName(""), WithLevel(LogLevelWarn))

var (
	enabledDebugs           = getEnabledDebugs()
	nopDebugf     DebugFunc = func(string, ...interface{}) {}
)

func getEnabledDebugs() map[string]bool {
	ret := map[string]bool{}
	debugVar := os.Getenv("GL_DEBUG")
	if debugVar == "" {
		return ret
	}

	for _, tag := range strings.Split(debugVar, ",") {
		ret[tag] = true
	}

	return ret
}

func HaveDebugTag(tag string) bool {
	return enabledDebugs[tag]
}

func SetupVerboseLog(log Log, isVerbose bool) {
	if isVerbose {
		log.SetLevel(LogLevelInfo)
	}
}

func Debug(tag string) DebugFunc {
	if !enabledDebugs[tag] {
		return nopDebugf
	}

	l := NewLogger(WithName(tag), WithLevel(LogLevelDebug))
	return func(format string, args ...interface{}) {
		l.Debugf(format, args...)
	}
}
