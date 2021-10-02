package installation

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"

	"github.com/alex-held/devctl/pkg/constants"
	"github.com/alex-held/devctl/pkg/env"
	"github.com/alex-held/devctl/pkg/index/scanner"
	"github.com/alex-held/devctl/pkg/index/spec"
)

// Store saves the given receipt at the destination.
// The caller has to ensure that the destination directory exists.
func Store(fs afero.Fs, receipt spec.Receipt, dest string) error {
	yamlBytes, err := yaml.Marshal(receipt)
	if err != nil {
		return errors.Wrapf(err, "convert to yaml")
	}

	err = afero.WriteFile(fs, dest, yamlBytes, 0644)
	return errors.Wrapf(err, "write plugin receipt %q", dest)
}

// Load reads the plugin receipt at the specified destination.
// If not found, it returns os.IsNotExist error.
func Load(fs afero.Fs, path string) (spec.Receipt, error) {
	return scanner.ReadReceiptFromFile(fs, path)
}

// New returns a new receipt with the given plugin and index name.
func New(plugin spec.Plugin, indexName string, timestamp metav1.Time) spec.Receipt {
	plugin.CreationTimestamp = timestamp
	return spec.Receipt{
		Plugin: plugin,
		Status: spec.ReceiptStatus{
			Source: spec.SourceIndex{
				Name: indexName,
			},
		},
	}
}

// InstalledPluginsFromIndex returns a list of all install plugins from a particular spec.
func InstalledPluginsFromIndex(f env.Factory, indexName string) ([]spec.Receipt, error) {
	var out []spec.Receipt
	receipts, err := GetInstalledPluginReceipts(f)
	if err != nil {
		return nil, err
	}
	for _, r := range receipts {
		if r.Status.Source.Name == indexName {
			out = append(out, r)
		}
	}
	return out, nil
}

// GetInstalledPluginReceipts returns a list of receipts.
func GetInstalledPluginReceipts(f env.Factory) ([]spec.Receipt, error) {
	receiptsDir := f.Paths().InstallReceiptsPath()
	files, err := filepath.Glob(filepath.Join(receiptsDir, "*"+constants.ManifestExtension))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to glob receipts directory (%s) for manifests", receiptsDir)
	}
	out := make([]spec.Receipt, 0, len(files))
	for _, file := range files {
		r, err := Load(f.Fs(), file)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse plugin install receipt %s", f)
		}
		out = append(out, r)
		klog.V(4).Infof("parsed receipt for %s: version=%s", r.GetObjectMeta().GetName(), r.Spec.Version)

	}
	return out, nil
}
