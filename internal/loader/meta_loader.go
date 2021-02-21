package loader

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/alex-held/devctl/internal/meta"
)

// MetaLoader loads a chart.
type MetaLoader interface {
	Load() (*meta.Meta, error)
}

// Loader returns a new MetaLoader appropriate for the given chart name
func Loader(name string) (MetaLoader, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return DirLoader(name), nil
	}
	return FileLoader(name), nil
}

// Load takes a string name, tries to resolve it to a file or directory, and then loads it.
//
// This is the preferred way to load a chart. It will discover the chart encoding
// and hand off to the appropriate chart reader.
//
// If a .helmignore file is present, the directory loader will skip loading any files
// matching it. But .helmignore is not evaluated when reading out of an archive.
func Load(name string) (*meta.Meta, error) {
	l, err := Loader(name)
	if err != nil {
		return nil, err
	}
	return l.Load()
}

// BufferedFile represents an archive file buffered for later processing.
type BufferedFile struct {
	Name string
	Data []byte
}

// LoadFiles loads from in-memory files.
// nolint:gocognit
func LoadFiles(files []*BufferedFile) (*meta.Meta, error) {
	c := new(meta.Meta)
	subcharts := make(map[string][]*BufferedFile)

	// do not rely on assumed ordering of files in the chart and crash
	// if Meta.yaml was not coming early enough to initialize metadata
	for _, f := range files {
		c.Raw = append(c.Raw, &meta.File{Name: f.Name, Data: f.Data})
		if f.Name == "Meta.yaml" {
			if c.Metadata == nil {
				c.Metadata = new(meta.Metadata)
			}
			if err := yaml.Unmarshal(f.Data, c.Metadata); err != nil {
				return c, errors.Wrap(err, "cannot load Meta.yaml")
			}
			if c.Metadata.APIVersion == "" {
				c.Metadata.APIVersion = meta.APIVersionV1
			}
		}
	}
	for _, f := range files {
		switch {
		case f.Name == "Meta.yaml":
			// already processed
			continue
		case f.Name == "values.yaml":
			c.Values = make(map[string]interface{})
			if err := yaml.Unmarshal(f.Data, &c.Values); err != nil {
				return c, errors.Wrap(err, "cannot load values.yaml")
			}
		case f.Name == "values.schema.json":
			c.Schema = f.Data

		case strings.HasPrefix(f.Name, "templates/"):
			c.Templates = append(c.Templates, &meta.File{Name: f.Name, Data: f.Data})

		case strings.HasPrefix(f.Name, "charts/"):
			if filepath.Ext(f.Name) == ".prov" {
				c.Files = append(c.Files, &meta.File{Name: f.Name, Data: f.Data})
				continue
			}

			fname := strings.TrimPrefix(f.Name, "charts/")
			cname := strings.SplitN(fname, "/", 2)[0]
			subcharts[cname] = append(subcharts[cname], &BufferedFile{Name: fname, Data: f.Data})
		default:
			c.Files = append(c.Files, &meta.File{Name: f.Name, Data: f.Data})
		}
	}

	if c.Metadata == nil {
		return c, errors.New("Meta.yaml file is missing")
	}

	if err := c.Validate(); err != nil {
		return c, err
	}

	return c, nil
}
