package install

import (
	. "github.com/alex-held/dev-env/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup(t *testing.T, config Config) (*assert.Assertions, afero.Fs) {
	fs := afero.NewMemMapFs()
	_ = config.WriteToFile(fs, "config.json")
	return assert.New(t), fs
}

func TestInstallAddsSDKToConfig(t *testing.T) {
	a, fs := setup(t, Config{Sdks: []SDK{
		{
			Name:    "java",
			Version: "1.8",
			Path:    "java-1.8",
		},
	}})

	executeInstall(fs, []string{"dotnet", "3.1.100"})
	config, _ := ReadConfigFromFile(fs, "config.json")
	a.Len(config.Sdks, 2)
	dotnetSDK := config.Sdks[1]
	a.Equal("dotnet", dotnetSDK.Name)
	a.Equal("3.1.100", dotnetSDK.Version)
	a.Equal("dotnet-3.1.100", dotnetSDK.Path)
}
