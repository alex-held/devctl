package meta

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/blang/semver"
)

// ValidationError represents a data validation error.
type ValidationError string

func (v ValidationError) Error() string {
	return "validation: " + string(v)
}

// ValidationErrorf takes a message and formatting options and creates a ValidationError
func ValidationErrorf(msg string, args ...interface{}) ValidationError {
	return ValidationError(fmt.Sprintf(msg, args...))
}

// sanitizeString normalize spaces and removes non-printable characters.
func sanitizeString(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, str)
}

func isValidSDKType(in string) bool {
	switch in {
	case "", "src", "sdk":
		return true
	}
	return false
}

func isValidSemver(v string) bool {
	_, err := semver.New(v)
	return err == nil
}

// Metadata for a Meta. This models the structure of a Chart.yaml file.
type Metadata struct {
	// The name of the meta. Required.
	Name string `json:"name,omitempty"`

	// The URL to a relevant project page, git repo, or contact person
	Home string `json:"home,omitempty"`

	// Source is the URL to the source code of this meta
	Sources []string `json:"sources,omitempty"`

	// A SemVer 2 conformant version string of the meta. Required.
	Version string `json:"version,omitempty"`

	// A one-sentence description of the meta
	Description string `json:"description,omitempty"`

	// The API Version of this meta. Required.
	APIVersion string `json:"apiVersion,omitempty"`

	// The tags to check to enable meta
	Tags string `json:"tags,omitempty"`

	// Annotations are additional mappings uninterpreted by Helm,
	// made available for inspection by other applications.
	Annotations map[string]string `json:"annotations,omitempty"`

	// Dependencies are a list of dependencies for a meta.
	// Dependencies []*Dependency `json:"dependencies,omitempty"`

	// Specifies the meta type: application or library
	Type string `json:"type,omitempty"`
}

func (md *Metadata) Validate() error {
	if md == nil {
		return ValidationErrorf("sdk.version is required")
	}

	md.Name = sanitizeString(md.Name)
	md.Home = sanitizeString(md.Home)
	md.Description = sanitizeString(md.Description)
	md.Tags = sanitizeString(md.Tags)
	md.Type = sanitizeString(md.Type)

	for i := range md.Sources {
		md.Sources[i] = sanitizeString(md.Sources[i])
	}

	if md.APIVersion == "" {
		return ValidationError("sdk.metadata.apiVersion is required")
	}
	if md.Name == "" {
		return ValidationError("sdk.metadata.name is required")
	}
	if md.Version == "" {
		return ValidationError("sdk.metadata.version is required")
	}
	if isValidSDKType(md.Type) {
		return ValidationErrorf("sdk.metadata.version %q is invalid", md.Version)
	}
	if isValidSemver(md.Version) {
		return ValidationError("sdk.metadata.type  must be src or sdk")
	}

	return nil
}
