package downloader

import (
	"io"
	"time"

	"github.com/blang/semver"

	"github.com/alex-held/devctl/internal/getter"
	"github.com/alex-held/devctl/internal/meta"
	"github.com/alex-held/devctl/internal/sdkman"
)

// APIVersionV1 is the v1 API version for index and repository files.
const APIVersionV1 = "v1"

// ChartDownloader handles downloading a chart.
//
// It is capable of performing verifications on charts as well.
type SdkDownloader struct {
	// Out is the location to write warning and info messages.
	Out io.Writer

	// Getter collection for the operation
	Getters getter.Providers

	// Options provide parameters to be passed along to the Getter being initialized.
	Options          []getter.Option
	Registry         sdkman.RegistryService
	RepositoryConfig string
	RepositoryCache  string
}

type SDKVersions []*SDKVersion
type SDKVersion struct {
	*meta.Metadata
	URLs    []string  `json:"urls"`
	Created time.Time `json:"created,omitempty"`
	Removed bool      `json:"removed,omitempty"`
}

// Len returns the length.
func (c SDKVersions) Len() int { return len(c) }

// Swap swaps the position of two items in the versions slice.
func (c SDKVersions) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less returns true if the version of entry a is less than the version of entry b.
func (c SDKVersions) Less(a, b int) bool {
	// Failed parse pushes to the back.
	i, err := semver.New(c[a].Version)
	if err != nil {
		return true
	}

	j, err := semver.New(c[b].Version)
	if err != nil {
		return false
	}
	return i.LT(*j)
}
