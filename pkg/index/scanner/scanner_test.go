package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alex-held/gold"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/alex-held/devctl/pkg/env"
)

func TestLoadIndexFileFromFS(t *testing.T) {
	type args struct {
		pluginsDir string
		pluginName string
	}
	tests := []struct {
		name              string
		args              args
		wantErr           bool
		wantIsNotExistErr bool
	}{
		{
			name: "load single index file",
			args: args{
				pluginsDir: filepath.FromSlash("./testdata/testindex/plugins"),
				pluginName: "foo",
			},

			wantErr:           false,
			wantIsNotExistErr: false,
		},
		{
			name: "plugin file not found",
			args: args{
				pluginsDir: filepath.FromSlash("./testdata/plugins"),
				pluginName: "not",
			},
			wantErr:           true,
			wantIsNotExistErr: true,
		},
		{
			name: "plugin file bad name",
			args: args{
				pluginsDir: filepath.FromSlash("./testdata/plugins"),
				pluginName: "wrongname",
			},
			wantErr:           true,
			wantIsNotExistErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := env.NewFactory()

			got, err := LoadPluginByName(f, tt.args.pluginsDir, tt.args.pluginName)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadIndexFileFromFS() got = %T,error = %v, wantErr %v", got, err, tt.wantErr)
				return
			}
			if os.IsNotExist(err) != tt.wantIsNotExistErr {
				t.Errorf("LoadIndexFileFromFS() got = %##v,error = %v, wantIsNotExistErr %v", got, err, tt.wantErr)
				return
			}
		})
	}
}

func touch(t *testing.T, filename string) {
	filename = strings.TrimPrefix(filename, "./testdata/")
	g := G(t)
	g.Assert(t, filename, byteArr(""))
}

func byteArr(s string) []byte {
	return []byte(s)
}

func G(t *testing.T) *gold.Gold {
	return gold.New(t, goldie.WithNameSuffix(""), goldie.WithSubTestNameForDir(false), goldie.WithTestNameForDir(false))
}

func TestFindPluginManifests(t *testing.T) {

	type args struct {
		indexFilePath string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		matchFirst labels.Set
	}{
		{
			name: "read index file",
			args: args{
				indexFilePath: filepath.Join("index", "plugins", "foo.yaml"),
			},
			wantErr: false,
			matchFirst: labels.Set{
				"os": "macos",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}

}

func TestReadPluginFile(t *testing.T) {
	type args struct {
		indexFilePath string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		matchFirst labels.Set
	}{
		{
			name: "read index file",
			args: args{
				indexFilePath: filepath.Join("testdata", "testindex", "plugins", "foo.yaml"),
			},
			wantErr: false,
			matchFirst: labels.Set{
				"os": "macos",
			},
		},
	}
	neverMatch := labels.Set{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ReadPluginFromFile(afero.NewOsFs(), tt.args.indexFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPluginFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Name != "foo" && got.Kind != "Plugin" {
				t.Errorf("ReadPluginFromFile() has not parsed the metainformations %v", got)
				return
			}

			sel, err := metav1.LabelSelectorAsSelector(got.Spec.Platforms[0].Selector)
			if err != nil {
				t.Errorf("ReadPluginFromFile() error parsing label err: %v", err)
				return
			}
			if !sel.Matches(tt.matchFirst) || sel.Matches(neverMatch) {
				t.Errorf("ReadPluginFromFile() didn't parse label selector properly: %##v", sel)
				return
			}
		})
	}
}
