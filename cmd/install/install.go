package install

import (
	"fmt"
	"github.com/alex-held/dev-env/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeInstall(afero.NewOsFs(), args)
	},
}

type RemoteSDKProvider interface {
	GetLatestVersion(sdk string) string
	Install(sdk string, version string, path string) error
}

type SDKProvider struct {
	Config *config.Config
}

func (provider *SDKProvider) GetLatestVersion(sdk string) string {
	switch sdk {
	case "java":
		return "1.8"
	case "dotnet":
		return "3.1.100"
	default:
		return "1.0"
	}
}

func (provider *SDKProvider) Install(sdk string, version string) error {
	installPath := fmt.Sprintf("%s-%s", sdk, version)
	provider.Config.AddSDK(sdk, version, installPath)
	return nil
}

func executeInstall(fs afero.Fs, args []string) {
	config, err := config.ReadConfigFromFile(fs, "config.json")
	provider := SDKProvider{Config: config}

	sdk := args[0]
	version := provider.GetLatestVersion(sdk)
	err = provider.Install(sdk, version)

	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}
}
