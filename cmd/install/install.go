package install

import (
	"fmt"
	"github.com/alex-held/dev-env/config"
	. "github.com/alex-held/dev-env/execution"
	. "github.com/alex-held/dev-env/manifest"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
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

//Install
func Install(manifest Manifest) error {
	executor := NewCommandExecutor(&manifest, func(str string) {
		fmt.Print(str)
	})
	err := executor.Execute()
	if err != nil {
		return err
	}
	fmt.Println("Successfully installed " + manifest.SDK)
	return nil
}

func executeInstall(fs afero.Fs, args []string) {
	cfg, err := config.ReadConfigFromFile(fs, "config.json")
	_ = config.GetManifests()
	if cfg == nil {
		os.Exit(1)
		return
	}

	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}

	cfg.AddSDK(args[0], args[1])
}
