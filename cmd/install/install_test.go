package install

import (
	"github.com/alex-held/dev-env/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup(t *testing.T, config config.Config) (*assert.Assertions, afero.Fs) {
	fs := afero.NewMemMapFs()
	_ = config.WriteToFile(fs, "config.json")
	return assert.New(t), fs
}

func TestInstallAddsSDKToConfig(t *testing.T) {
	a, fs := setup(t, config.Config{Sdks: []config.SDK{
		{
			Name:    "java",
			Version: "1.8",
			Path:    "java-1.8",
		},
	}})

	executeInstall(fs, []string{"dotnet", "3.1.100"})
	config, _ := config.ReadConfigFromFile(fs, "config.json")
	a.Len(config.Sdks, 2)
	dotnetSDK := config.Sdks[1]
	a.Equal(dotnetSDK.Name, "dotnet")
	a.Equal(dotnetSDK.Version, "3.1.100")
	a.Equal(dotnetSDK.Name, "dotnet-3.1.00")
}
