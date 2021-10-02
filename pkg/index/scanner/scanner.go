// Copyright 2019 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// COPIED AND MODIFIED FROM: https://github.com/kubernetes-sigs/krew

package scanner

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alex-held/devctl-kit/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	errors2 "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/yaml"

	"github.com/alex-held/devctl/internal/git"
	"github.com/alex-held/devctl/pkg/constants"
	"github.com/alex-held/devctl/pkg/env"
	"github.com/alex-held/devctl/pkg/index/spec"
	"github.com/alex-held/devctl/pkg/index/validate"
)

type Plugins []spec.Plugin
type PluginMap map[string]spec.Plugin

func (plugins *PluginMap) Filter(predicate func(key string, p spec.Plugin) bool) (res PluginMap) {
	res = PluginMap{}
	for k, v := range *plugins {
		if predicate(k, v) {
			res[k] = v
		}
	}
	return res
}

func (plugins *Plugins) Filter(predicate func(p spec.Plugin) bool) (res Plugins) {
	for _, v := range *plugins {
		if predicate(v) {
			res = append(res, v)
		}
	}
	return res
}

func (plugins *Plugins) Map(fn func(p spec.Plugin) interface{}) (res []interface{}) {
	for _, v := range *plugins {
		res = append(res, fn(v))
	}
	return res
}

func (plugins *PluginMap) Values() (res []spec.Plugin) {
	for _, v := range *plugins {
		res = append(res, v)
	}
	return res
}

func (plugins *PluginMap) Keys() (res []string) {
	for k := range *plugins {
		res = append(res, k)
	}
	return res
}

func LoadPluginsFromFS(f env.Factory, indexName string) (plugins []spec.Plugin, errors errors2.Aggregate) {
	paths := f.Paths()
	indexDir := paths.IndexPluginsPath(indexName)
	files, err := findPluginManifestFiles(indexDir)
	if err != nil {
		return nil, errors2.NewAggregate([]error{err})
	}
	log.Debugf("found %d plugins in dir %s", len(files), indexDir)

	var errs []error
	for _, file := range files {
		pluginName := strings.TrimSuffix(file, filepath.Ext(file))
		p, err := LoadPluginByName(f, indexDir, pluginName)
		if err != nil {
			errs = append(errs, err)
			log.Errorf("failed to read or parse plugin manifest %q: %v", pluginName, err)
			continue
		}
		plugins = append(plugins, p)
	}

	return plugins, errors2.NewAggregate(errs)
}

func LoadPluginByName(f env.Factory, pluginsDir string, pluginName string) (p spec.Plugin, err error) {
	pluginFile := filepath.Join(pluginsDir, pluginName+constants.ManifestExtension)
	return ReadPluginFromFile(f.Fs(), pluginFile)
}

func ReadPluginFromFile(fs afero.Fs, path string) (p spec.Plugin, err error) {
	p = spec.Plugin{}
	err = readFromFile(fs, path, &p)
	return p, err
}

func ReadPlugin(f io.ReadCloser) (spec.Plugin, error) {
	var plugin spec.Plugin
	err := decodeFile(f, &plugin)
	if err != nil {
		return plugin, errors.Wrap(err, "failed to decode plugin manifest")
	}
	return plugin, errors.Wrap(validate.ValidatePlugin(plugin.Name, plugin), "plugin manifest validation error")
}

// ReadReceiptFromFile loads a file from the FS. When receipt file not found, it
// returns an error that can be checked with os.IsNotExist.
func ReadReceiptFromFile(fs afero.Fs, path string) (spec.Receipt, error) {
	var receipt spec.Receipt
	err := readFromFile(fs, path, &receipt)
	if receipt.Status.Source.Name == "" {
		receipt.Status.Source.Name = constants.DefaultIndexName
	}
	return receipt, err
}

func readFromFile(fs afero.Fs, path string, as interface{}) error {
	file, err := fs.Open(path)
	if err != nil {
		return err
	}
	err = decodeFile(file, &as)
	return errors.Wrapf(err, "failed to parse yaml file %q", path)
}

// decodeFile tries to decode a plugin/receipt
func decodeFile(r io.ReadCloser, as interface{}) error {
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, &as)
}

func findPluginManifestFiles(indexDir string) ([]string, error) {
	var out []string
	files, err := ioutil.ReadDir(indexDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open index dir")
	}
	for _, file := range files {
		if file.Mode().IsRegular() && filepath.Ext(file.Name()) == constants.ManifestExtension {
			out = append(out, file.Name())
		}
	}
	return out, nil
}

func DefaultIndex() string {
	if uri := os.Getenv(constants.DEVCTL_DEFAULT_INDEX_URI_KEY); uri != "" {
		return uri
	}
	return constants.DefaultIndexURI
}

var validNamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// Index describes the name and URL of a configured index.
type Index struct {
	Name string
	URL  string
}

// ListIndexes returns a slice of Index objects. The path argument is used as
// the base path of the index.
func ListIndexes(paths env.Paths) ([]Index, error) {
	entries, err := ioutil.ReadDir(paths.IndexBase())
	if err != nil {
		return nil, errors.Wrap(err, "failed to list directory")
	}

	indexes := []Index{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		indexName := e.Name()
		remote, err := git.GetRemoteURL(paths.IndexPath(indexName))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list the remote URL for index %s", indexName)
		}

		indexes = append(indexes, Index{
			Name: indexName,
			URL:  remote,
		})
	}
	return indexes, nil
}

// AddIndex initializes a new index to install plugins from.
func AddIndex(paths env.Paths, name, url string) error {
	dir := paths.IndexPath(name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return git.EnsureCloned(url, dir)
	} else if err != nil {
		return err
	}
	return errors.New("index already exists")
}

// DeleteIndex removes specified index name. If index does not exist, returns an error that can be tested by os.IsNotExist.
func DeleteIndex(paths env.Paths, name string) error {
	dir := paths.IndexPath(name)
	if _, err := os.Stat(dir); err != nil {
		return err
	}

	return os.RemoveAll(dir)
}

// IsValidIndexName validates if an index name contains invalid characters
func IsValidIndexName(name string) bool {
	return validNamePattern.MatchString(name)
}
