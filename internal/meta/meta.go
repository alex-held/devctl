package meta

import (
	"path/filepath"
	"regexp"
	"strings"
)

// APIVersionV1 is the API version number for version 1.
const APIVersionV1 = "v1"

// APIVersionV2 is the API version number for version 2.
const APIVersionV2 = "v2"

// aliasNameFormat defines the characters that are legal in an alias name.
var aliasNameFormat = regexp.MustCompile("^[a-zA-Z0-9_-]+$")

// Chart is a helm package that contains metadata, a default config, zero or more
// optionally parameterizable templates, and zero or more charts (dependencies).
type Meta struct {

	// Raw contains the raw contents of the files originally contained in the chart archive.
	//
	// This should not be used except in special cases like `helm show values`,
	// where we want to display the raw values, comments and all.
	Raw []*File `json:"-"`

	// Metadata is the contents of the Chartfile.
	Metadata *Metadata `json:"metadata"`

	// Templates for this chart.
	Templates []*File `json:"templates"`

	// Values are default config for this chart.
	Values map[string]interface{} `json:"values"`

	// Schema is an optional JSON schema for imposing structure on Values
	Schema []byte `json:"schema"`

	// Files are miscellaneous files in a chart archive,
	// e.g. README, LICENSE, etc.
	Files []*File `json:"files"`

	parent *Meta
	// dependencies []*Meta
}

// Name returns the name of the chart.
func (ch *Meta) Name() string {
	if ch.Metadata == nil {
		return ""
	}
	return ch.Metadata.Name
}

// MetaPath returns the full path to this chart in dot notation.
func (ch *Meta) MetaPath() string {
	if !ch.IsRoot() {
		return ch.Parent().MetaPath() + "." + ch.Name()
	}
	return ch.Name()
}

// MetaFullPath returns the full path to this chart.
func (ch *Meta) MetaFullPath() string {
	if !ch.IsRoot() {
		return ch.Parent().MetaFullPath() + "/charts/" + ch.Name()
	}
	return ch.Name()
}

// Validate validates the metadata.
func (ch *Meta) Validate() error {
	return ch.Metadata.Validate()
}

func (m *Meta) IsRoot() bool  { return m.parent == nil }
func (m *Meta) Parent() *Meta { return m.parent }

// Root finds the root chart.
func (m *Meta) Root() *Meta {
	if m.IsRoot() {
		return m
	}
	return m.Parent().Root()
}

func hasManifestExtension(fname string) bool {
	ext := filepath.Ext(fname)
	return strings.EqualFold(ext, ".yaml") || strings.EqualFold(ext, ".yml") || strings.EqualFold(ext, ".json")
}
