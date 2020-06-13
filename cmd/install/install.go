package install

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alex-held/dev-env/config"
)

// installCmd represents the install command
//noinspection GoUnusedGlobalVariable
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Installed called: %v", args)
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
