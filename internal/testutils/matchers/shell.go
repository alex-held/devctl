package matchers

import (
	"fmt"
	"strings"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

type ContainExportMatcher struct {
	Key   string
	Value string
	Args  []interface{}
}

func (m *ContainExportMatcher) stringToMatch() string {
	stringToMatch := fmt.Sprintf("export %s=%s\n", m.Key, m.Value)
	if len(m.Args) > 0 {
		stringToMatch = fmt.Sprintf(stringToMatch, m.Args...)
	}
	return stringToMatch
}

func (m ContainExportMatcher) Match(actual interface{}) (success bool, err error) {
	var actualString string
	switch a := actual.(type) {
	case string:
		actualString = a
	case fmt.Stringer:
		actualString = a.String()
	case fmt.GoStringer:
		actualString = a.GoString()
	default:
		return false, fmt.Errorf("ContainExport m requires a string or stringer.  Got:\n%s", format.Object(actual, 1))
	}
	return strings.Contains(actualString, m.stringToMatch()), nil
}

func (m ContainExportMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to contain substring", m.stringToMatch())
}

func (m ContainExportMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to contain substring", m.stringToMatch())
}

// ContainExport succeeds if actual is a string or stringer that contains an export declaration
// the passed-in key and value.
// Optional arguments can be provided to construct the substring
// via fmt.Sprintf().
func ContainExport(key, value string, args ...interface{}) types.GomegaMatcher {
	return &ContainExportMatcher{
		Key:   key,
		Value: value,
		Args:  args,
	}
}

// HaveSectionHeader succeeds if actual is a string or stringer that contains a section header with
// the passed-in key title
// Optional arguments can be provided to construct the substring
// via fmt.Sprintf().
func HaveSectionHeader(title string, args ...interface{}) types.GomegaMatcher {
	header := fmt.Sprintf(
		`# --------------------------------------------------
# %s
# --------------------------------------------------`, title)

	return &matchers.ContainSubstringMatcher{
		Substr: header,
		Args:   args,
	}
}
