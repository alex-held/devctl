package index

import (
	"io"
	"io/ioutil"
	log2 "log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alex-held/devctl-kit/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	errors2 "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/yaml"

	"github.com/alex-held/devctl/pkg/constants"
	"github.com/alex-held/devctl/pkg/env"
)

func FindPluginManifests(f env.Factory) (plugins []Plugin, errors errors2.Aggregate) {
	paths := f.Paths()
	indexDir := paths.IndexBase()
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

func LoadPluginByName(f env.Factory, pluginsDir string, pluginName string) (p Plugin, err error) {
	pluginFile := filepath.Join(pluginsDir, pluginName+constants.ManifestExtension)
	log2.Printf("\n---\npluginsDir: %s\npluginName: %s\npluginFile: %s\n---\n", pluginsDir, pluginName, pluginFile)
	return ReadPluginFromFile(f.Fs(), pluginFile)
}

func ReadPluginFromFile(fs afero.Fs, path string) (p Plugin, err error) {
	p = Plugin{}
	err = readFromFile(fs, path, &p)
	return p, err
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
