package exec

import (
	"context"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/gobuffalo/plugins/plugio"

	"github.com/alex-held/devctl/pkg/cli"
)


func Run(ctx context.Context, root string, args []string) error {
	main := filepath.Join(root, "cmd", "devctl")
	if _, err := os.Stat(filepath.Dir(main)); err != nil {
		buff, err := cli.New()
		if err != nil {
			return err
		}
		return buff.Main(ctx, root, args)
	}
	bargs := []string{"run", "-v", "./cmd/devctl"}
	bargs = append(bargs, args...)

	cmd := exec.CommandContext(ctx, "go", bargs...)
	cmd.Stdin = plugio.Stdin()
	cmd.Stdout = plugio.Stdout()
	cmd.Stderr = plugio.Stderr()
	return cmd.Run()
}

// CreatePluginCmd provides a lightweight version to delegate the execution to the plugins
func CreatePluginCmd(ctx context.Context, args []string) (cmd *exec.Cmd, err error) {
	var pluginRoot = os.Getenv("DEVCTL_PLUGIN_ROOT")
	if pluginRoot == "" {
		devctlRoot, err := ResolveDevctlRoot()
		if err != nil {
			return nil, err
		}
		pluginRoot = path.Join(devctlRoot, "plugins")
	}

	pluginPath := pluginRoot
	for _, arg := range args {
		pluginPath = path.Join(pluginPath, arg)
	}
	pluginExecArgs := []string{"run", "-v", pluginPath}
	cmd = exec.CommandContext(ctx, "go", pluginExecArgs...)
	cmd.Stdin = plugio.Stdin()
	cmd.Stdout = plugio.Stdout()
	cmd.Stderr = plugio.Stderr()
	return cmd, nil
}

func ResolveDevctlRoot() (string, error) {
	const root = "devctl"
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for path.Base(wd) != root {
		wd = path.Dir(wd)
	}

	return wd, nil
}
