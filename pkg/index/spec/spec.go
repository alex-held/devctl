package spec

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Plugin describes a plugin manifest file.
type Plugin struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata"`

	Spec PluginSpec `json:"spec"`
}

// PluginSpec is the plugin specification.
type PluginSpec struct {
	Version          string `json:"version,omitempty"`
	ShortDescription string `json:"shortDescription,omitempty"`
	Description      string `json:"description,omitempty"`
	Caveats          string `json:"caveats,omitempty"`
	Homepage         string `json:"homepage,omitempty"`

	Platforms []Platform `json:"platforms,omitempty"`
}

// Platform describes how to perform an installation on a specific platform
// and how to match the target platform (os, arch).
type Platform struct {
	URI    string `json:"uri,omitempty"`
	Sha256 string `json:"sha256,omitempty"`

	Selector *metav1.LabelSelector `json:"selector,omitempty"`
	Files    []FileOperation       `json:"files"`

	// Bin specifies the path to the plugin executable.
	// The path is relative to the root of the installation folder.
	// The binary will be linked after all FileOperations are executed.
	Bin string `json:"bin"`
}

// FileOperation specifies a file copying operation from plugin archive to the
// installation directory.
type FileOperation struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// Receipt describes a plugin receipt file.
type Receipt struct {
	Plugin `json:",inline" yaml:",inline"`

	Status ReceiptStatus `json:"status"`
}

// ReceiptStatus contains information about the installed plugin.
type ReceiptStatus struct {
	Source SourceIndex `json:"source"`
}

// SourceIndex contains information about the index a plugin was installed from.
type SourceIndex struct {
	// Name is the configured name of an index a plugin was installed from.
	Name string `json:"name"`
}
