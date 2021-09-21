package golang

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

type Config struct {
	InstallPath string `yaml:"install_path"`
}

// CreateConfig creates the default plugin config
func CreateConfig(devctlPath string) map[string]string {
	goSDKInstallPath := devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
		return devctlPath
	})).SDK("go")

	return map[string]string{
		goSDKInstallPathKey: goSDKInstallPath,
	}
}

const goSDKInstallPathKey = "goSDKInstallPath"

func New(cfgMap map[string]string, args []string) (err error) {
	args = args[1:]

	cfg := &Config{InstallPath: cfgMap[goSDKInstallPathKey]}

	if len(args) == 0 {
		usage()
		return fmt.Errorf("must atleast provide one argument")
	}

	subcmd := args[0]
	args = args[1:]

	switch subcmd {
	case "list":
		if err = validateArgsForSubcommand(subcmd, args, 1); err != nil {
			return err
		}
		return handleList(cfg)
	case "current":
		if err = validateArgsForSubcommand(subcmd, args, 1); err != nil {
			return err
		}
		return handleCurrent(cfg)
	case "install":
		if err = validateArgsForSubcommand(subcmd, args, 2); err != nil {
			return err
		}
		return handleInstall(args[0], cfg)
	case "use":
		if err = validateArgsForSubcommand(subcmd, args, 2); err != nil {
			return err
		}
		return handleUse(args[0], cfg)
	default:
		return fmt.Errorf("unknown subcommand '%s'; expected on of 'list, current, install, use'", subcmd)
	}
}

func handleInstall(version string, config *Config) (err error) {
	installPath := path.Join(config.InstallPath, version)
	fmt.Printf("[sdk/go] installing go sdk version %s into %s", version, installPath)

	return err
}

func handleList(config *Config) (err error) {
	installPath := path.Join(config.InstallPath)
	ds, err := os.ReadDir(installPath)

	var installedVersions []string

	if err != nil {
		return err
	}

	for _, d := range ds {
		if d.IsDir() && strings.HasPrefix(d.Name(), "v") {
			installedVersions = append(installedVersions, d.Name())
		}
	}

	output := strings.Builder{}
	for _, version := range installedVersions {
		output.WriteString(version + "\n")
	}

	fmt.Printf("[sdk/go] installed versions:\n%s\n", output.String())
	return nil
}

func handleUse(version string, config *Config) (err error) {
	installPath := path.Join(config.InstallPath, version)
	fmt.Printf("[sdk/go] installing go sdk version %s into %s", version, installPath)

	return err
}

func handleCurrent(config *Config) (err error) {
	installPath := path.Join(config.InstallPath, "current")
	link, err := os.Readlink(installPath)
	if err != nil {
		return err
	}
	fmt.Println(link)
	return nil
}

func usage() {
	fmt.Printf("USAGE")
}

func validateArgsForSubcommand(subcmd string, args []string, expected int) error {
	if len(args) != expected {
		return fmt.Errorf("provided wrong number of argument for subcommand '%s'; expected=%d; provided=%d", subcmd, expected, len(args))
	}
	return nil
}
