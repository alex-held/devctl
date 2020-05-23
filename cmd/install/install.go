package install

import (
	"fmt"
	"github.com/alex-held/dev-env/config"
	. "github.com/alex-held/dev-env/manifest"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
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

func prepareLinkingCommands(m Manifest) []Instructing {
	var linkingCommands []Instructing

	for _, link := range m.ResolveLinks() {
		command := DevEnvCommand{
			Command: "ln",
			Args:    []string{"-s", link.Source, link.Target},
		}
		linkingCommands = append(linkingCommands, command)
	}
	return linkingCommands
}

type CommandExecutor struct {
	installationCommands []Instructing
	linkingCommands      []Instructing
}

func NewCommandExecutor(installationCommands []Instructing, linkingCommands []Instructing) *CommandExecutor {
	return &CommandExecutor{
		installationCommands,
		linkingCommands,
	}
}

func (executor *CommandExecutor) Execute() error {
	commands := executor.GetCommands()

	for _, command := range commands {
		switch c := command.(type) {
		case DevEnvCommand:

			fmt.Println("Executing command: " + c.Format())

			cmd := exec.Command(c.Command, c.Args...)
			output, err := cmd.Output()
			if err != nil {
				fmt.Println("Error while running command. Error:'" + err.Error())
				os.Exit(1)
			}

			println(output)
		}

		/*
		   if err := cmd.Run(); err != nil {
		       fmt.Println("Error while running command. Error:'" + err.Error())
		       return err
		   }*/
	}
	os.Exit(0)
	return nil
}

func (executor *CommandExecutor) GetCommands() []Instructing {
	linkingCommands := executor.linkingCommands
	installationCommands := executor.installationCommands
	commands := append(installationCommands, linkingCommands...)
	return commands
}

//Install
func Install(manifest Manifest) error {
	executor := NewCommandExecutor(manifest.ResolveInstructions(), prepareLinkingCommands(manifest))
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
