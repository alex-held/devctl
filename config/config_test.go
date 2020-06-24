package config

import (
	"encoding/json"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var (
	configFilePath = "config.json"
	jdk            = SDK{
		Name:    "java",
		Version: "1.8",
		Path:    "/some/path",
	}
	defaultConfig = Config{Sdks: []SDK{jdk}}
)

func setup(t *testing.T, json string) (*assert.Assertions, afero.Fs) {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, configFilePath, []byte(json), 0644)
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
	config, _ := ReadConfigFromFile(fs, configFilePath)

	sdk := config.Sdks[0]

	a.Equal(sdk.Name, "java")
	a.Equal(sdk.Version, "1.8")
	a.Equal(sdk.Path, "/some/path")
}

func TestWriteFile(t *testing.T) {
	a, fs := setup(t, "")
	config := defaultConfig
	err := config.WriteToFile(fs, configFilePath)
	if err != nil {
		t.Error(err)
	}

	result, _ := afero.ReadFile(fs, configFilePath)
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
	config := NewConfig(fs, configFilePath)
	a.Empty(config.Sdks)
}

func TestAddSDKAddsASDK(t *testing.T) {
	a, fs := setup(t, "")
	config := NewConfig(fs, configFilePath)
	config.AddSDK("java", "1.8")

	file, err := afero.ReadFile(fs, configFilePath)
	if err != nil {
		t.Error(err)
	}

	newConfig := NewConfig(fs, configFilePath)
	json.Unmarshal(file, &newConfig)
	a.Equal("java", newConfig.Sdks[0].Name)
	a.Equal("1.8", newConfig.Sdks[0].Version)
	a.Equal("java-1.8", newConfig.Sdks[0].Path)
}

func TestJConfig_ListMatchingSdks(t *testing.T) {
	a, fs := setup(t, "{}")
	config := NewConfig(fs, "config.json")
	config.Sdks = append(config.Sdks, jdk, SDK{
		Name:    "dotnet",
		Version: "3.1.100",
		Path:    "dotnet-3.1.100",
	}, SDK{
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
	config.Sdks = append(config.Sdks, jdk, SDK{
		Name:    "dotnet",
		Version: "3.1.100",
		Path:    "dotnet-3.1.100",
	}, SDK{
		Name:    "dotnet",
		Version: "2.1.0",
		Path:    "dotnet-2.1.0",
	})

	result := config.ListSdks()
	length := len(result)
	a.Equal(length, 3)
}
