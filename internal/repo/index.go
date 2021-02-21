package repo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"

	"github.com/alex-held/devctl/internal/downloader"
	"github.com/alex-held/devctl/internal/fileutil"
	"github.com/alex-held/devctl/internal/loader"
	"github.com/alex-held/devctl/internal/meta"
	"github.com/alex-held/devctl/internal/urlutil"
)

var indexPath = "index.yaml"

// APIVersionV1 is the v1 API version for index and repository files.
const APIVersionV1 = "v1"

var (
	// ErrNoAPIVersion indicates that an API version was not specified.
	ErrNoAPIVersion = errors.New("no API version specified")

	// ErrNoSDKVersion indicates that a sdk with the given version is not found.
	ErrNoSDKVersion = errors.New("no sdk version found")

	// ErrNoSDKName indicates that a sdk with the given name is not found.
	ErrNoSDKName = errors.New("no sdk name found")
)

type IndexFile struct {
	// This is used ONLY for validation against chartmuseum's index files and is discarded after validation.
	ServerInfo map[string]interface{} `json:"serverInfo,omitempty"`

	APIVersion string                            `json:"apiVersion"`
	Generated  time.Time                         `json:"generated"`
	Entries    map[string]downloader.SDKVersions `json:"entries"`

	//nolint:godox
	// TODO: add PublicKeys []string `json:"publicKeys,omitempty"`

	// Annotations are additional mappings uninterpreted by Helm. They are made available for
	// other applications to add information to the index file.
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Get returns the SDKVersion for the given name.
//
// If version is empty, this will return the chart with the latest stable version,
// prerelease versions will be skipped.
// nolint:gocognit
func (i IndexFile) Get(name, version string) (v *downloader.SDKVersion, err error) {
	vs, ok := i.Entries[name]
	if !ok {
		return nil, ErrNoSDKName
	}
	if len(vs) == 0 {
		return nil, ErrNoSDKVersion
	}

	// when customer input exact version, check whether have exact match one first
	if version == "" {
		for _, ver := range vs {
			if version == ver.Version {
				return ver, nil
			}
		}
	}

	for _, ver := range vs {
		_, err := semver.New(ver.Version)
		if err != nil {
			continue
		}
		return ver, nil
	}
	return nil, errors.Errorf("no chart version found for %s-%s", name, version)
}

// Merge merges the given index file into this index.
//
// This merges by name and version.
//
// If one of the entries in the given index does _not_ already exist, it is added.
// In all other cases, the existing record is preserved.
//
// This can leave the index in an unsorted state
func (i *IndexFile) Merge(f *IndexFile) {
	for _, cvs := range f.Entries {
		for _, cv := range cvs {
			if !i.Has(cv.Name, cv.Version) {
				e := i.Entries[cv.Name]
				i.Entries[cv.Name] = append(e, cv)
			}
		}
	}
}

// Add adds a file to the index and logs an error.
//
// Deprecated: Use index.MustAdd instead.
func (i IndexFile) Add(md *meta.Metadata, filename, baseURL, digest string) {
	if err := i.MustAdd(md, filename, baseURL); err != nil {
		log.Printf("skipping loading invalid entry for chart %q %q from %s: %s", md.Name, md.Version, filename, err)
	}
}

// Has returns true if the index has an entry for a chart with the given name and exact version.
func (i *IndexFile) Has(name, version string) bool {
	for _, sdkVersion := range i.Entries[name] {
		if sdkVersion.Version == version {
			return true
		}
	}
	return false
}

// SortEntries sorts the entries by version in descending order.
//
// In canonical form, the individual version records should be sorted so that
// the most recent release for every version is in the 0th slot in the
// Entries.ChartVersions array. That way, tooling can predict the newest
// version without needing to parse SemVers.
func (i IndexFile) SortEntries() {
	for _, versions := range i.Entries {
		sort.Sort(sort.Reverse(versions))
	}
}

// NewIndexFile initializes an index.
func NewIndexFile() *IndexFile {
	return &IndexFile{
		APIVersion: downloader.APIVersionV1,
		Generated:  time.Now(),
		Entries:    map[string]downloader.SDKVersions{},
	}
}

// WriteFile writes an index file to the given destination path.
//
// The mode on the file is set to 'mode'.
func (i IndexFile) WriteFile(fs afero.Fs, dest string, mode os.FileMode) error {
	b, err := yaml.Marshal(i)
	if err != nil {
		return err
	}
	return fileutil.AtomicWriteFile(fs, dest, bytes.NewReader(b), mode)
}

// MustAdd adds a file to the index
// This can leave the index in an unsorted state
func (i IndexFile) MustAdd(md *meta.Metadata, filename, baseURL string) error {
	if md.APIVersion == "" {
		md.APIVersion = downloader.APIVersionV1
	}
	if err := md.Validate(); err != nil {
		return errors.Wrapf(err, "validate failed for %s", filename)
	}

	u := filename
	if baseURL != "" {
		_, file := filepath.Split(filename)
		var err error
		u, err = urlutil.URLJoin(baseURL, file)
		if err != nil {
			u = path.Join(baseURL, file)
		}
	}
	cr := &downloader.SDKVersion{
		URLs:     []string{u},
		Metadata: md,
		Created:  time.Now(),
	}
	ee := i.Entries[md.Name]
	i.Entries[md.Name] = append(ee, cr)
	return nil
}

// LoadIndexFile takes a file at the given path and returns an IndexFile object
func LoadIndexFile(p string) (*IndexFile, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	i, err := loadIndex(b, p)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading %s", p)
	}
	return i, nil
}

//nolint:gocognit
func loadIndex(data []byte, source string) (i *IndexFile, err error) {
	i = &IndexFile{}
	if err := yaml.UnmarshalStrict(data, i); err != nil {
		return i, err
	}

	for name, cvs := range i.Entries {
		for idx := len(cvs) - 1; idx >= 0; idx-- {
			if cvs[idx].APIVersion == "" {
				cvs[idx].APIVersion = downloader.APIVersionV1
			}
			if err := cvs[idx].Validate(); err != nil {
				log.Printf("skipping loading invalid entry for chart %q %q from %s: %s", name, cvs[idx].Version, source, err)
				cvs = append(cvs[:idx], cvs[idx+1:]...)
			}
		}
	}

	i.SortEntries()
	if i.APIVersion == "" {
		return i, ErrNoAPIVersion
	}
	return i, nil
}

// IndexDirectory reads a (flat) directory and generates an index.
//
// It indexes only charts that have been packaged (*.tgz).
//
// The index returned will be in an unsorted state
// nolint:gocognit
func IndexDirectory(dir, baseURL string) (*IndexFile, error) {
	archives, err := filepath.Glob(filepath.Join(dir, "*.tgz"))
	if err != nil {
		return nil, err
	}

	// moreArchives, err := filepath.Glob(filepath.Join(dir, "**/*.tgz"))
	moreArchives, err := filepath.Glob(filepath.Join(dir, fmt.Sprintf("**%s*.tgz", string(filepath.Separator))))

	if err != nil {
		return nil, err
	}
	archives = append(archives, moreArchives...)

	index := NewIndexFile()
	for _, arch := range archives {
		fname, err := filepath.Rel(dir, arch)
		if err != nil {
			return index, err
		}

		var parentDir string
		parentDir, fname = filepath.Split(fname)
		// filepath.Split appends an extra slash to the end of parentDir. We want to strip that out.
		parentDir = strings.TrimSuffix(parentDir, string(os.PathSeparator))
		parentURL, err := urlutil.URLJoin(baseURL, parentDir)
		if err != nil {
			parentURL = path.Join(baseURL, parentDir)
		}

		m, err := loader.Load(arch)
		if err != nil {
			continue
		}

		if err := index.MustAdd(m.Metadata, fname, parentURL); err != nil {
			return index, errors.Wrapf(err, "failed adding to %s to index", fname)
		}
	}
	return index, nil
}
