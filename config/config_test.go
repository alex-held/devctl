package config

import (
	"encoding/json"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	filepath = "config.json"
	jdk      = SDK{
		Name:    "java",
		Version: "1.8",
		Path:    "/some/path",
	}
	defaultConfig = Config{Sdks: []SDK{jdk}}
)

func setup(t *testing.T, json string) (*assert.Assertions, afero.Fs) {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, filepath, []byte(json), 0644)
	return assert.New(t), fs
}

func TestReadConfigFromFile(t *testing.T) {
	a, fs := setup(t, `{
 "Sdks": [
  {
   "Name": "java",
   "Version": "1.8",
   "Path": "/some/path"
  }
 ]
}`)
	config, _ := ReadConfigFromFile(fs, filepath)

	sdk := config.Sdks[0]

	a.Equal(sdk.Name, "java")
	a.Equal(sdk.Name, "1.8")
	a.Equal(sdk.Name, "/some/path")
}

func TestWriteFile(t *testing.T) {
	a, fs := setup(t, "")
	config := defaultConfig
	err := config.WriteToFile(fs, filepath)
	if err != nil {
		t.Error(err)
	}

	result, _ := afero.ReadFile(fs, filepath)
	a.JSONEq(`{
 "Sdks": [
  {
   "Name": "java",
   "Version": "1.8",
   "Path": "/some/path"
  }
 ]
}`, string(result))
}

func TestNewConfig(t *testing.T) {
	a, fs := setup(t, "")
	config := NewConfig(fs, filepath)
	a.Empty(t, config.Sdks)
}

func TestAddSDKAddsASDK(t *testing.T) {
	a, fs := setup(t, "")
	config := NewConfig(fs, filepath)
	config.AddSDK("java", "1.8", "/some/path")

	a.Contains(config.Sdks, jdk)
	file, err := afero.ReadFile(fs, filepath)
	if err != nil {
		t.Error(err)
	}

	newConfig := NewConfig(fs, filepath)
	json.Unmarshal(file, &newConfig)
	a.Contains(newConfig.Sdks, jdk)
}

func TestJConfig_ListMatchingSdks(t *testing.T) {
	a, fs := setup(t, "")
	config := NewConfig(fs, "config.json")
	config.Sdks = append(config.Sdks, jdk)
	config.Sdks = append(config.Sdks, SDK{
		Name:    "dotnet",
		Version: "3.1.100",
		Path:    "dotnet-3.1.100",
	})
	config.Sdks = append(config.Sdks, SDK{
		Name:    "dotnet",
		Version: "2.1.0",
		Path:    "dotnet-2.1.0",
	})

	result := config.ListMatchingSdks(func(sdk SDK) bool {
		return sdk.Name == "dotnet"
	})

	a.Len(result, 2)
}

func TestJConfig_ListSdks(t *testing.T) {
	a, fs := setup(t, "")
	config := NewConfig(fs, "config.json")
	config.Sdks = append(config.Sdks, jdk)
	config.Sdks = append(config.Sdks, SDK{
		Name:    "dotnet",
		Version: "3.1.100",
		Path:    "dotnet-3.1.100",
	})
	config.Sdks = append(config.Sdks, SDK{
		Name:    "dotnet",
		Version: "2.1.0",
		Path:    "dotnet-2.1.0",
	})

	result := config.ListSdks()

	a.Len(result, 1)
}
