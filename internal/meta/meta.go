package meta

// APIVersionV1 is the API version number for version 1.
const APIVersionV1 = "v1"

// APIVersionV2 is the API version number for version 2.
const APIVersionV2 = "v2"

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
func (m *Meta) Name() string {
	if m.Metadata == nil {
		return ""
	}
	return m.Metadata.Name
}

// MetaPath returns the full path to this chart in dot notation.
func (m *Meta) MetaPath() string {
	if !m.IsRoot() {
		return m.Parent().MetaPath() + "." + m.Name()
	}
	return m.Name()
}

// MetaFullPath returns the full path to this chart.
func (m *Meta) MetaFullPath() string {
	if !m.IsRoot() {
		return m.Parent().MetaFullPath() + "/charts/" + m.Name()
	}
	return m.Name()
}

// Validate validates the metadata.
func (m *Meta) Validate() error {
	return m.Metadata.Validate()
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
