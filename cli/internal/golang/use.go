package golang

import (
	"context"
	"os"
	"time"

	"github.com/alex-held/devctl/cli/cmds/sdk"
	taskrunner2 "github.com/alex-held/devctl/pkg/ui/taskrunner"
	"github.com/gobuffalo/plugins"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/pkg/errors"

	// "github.com/alex-held/devctl/pkg/plugins"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

type GoUseCmd struct {
	path devctlpath.Pather
	fs   vfs.VFS

	plugins   plugins.Plugins
	executors []sdk.Sdker
}

func (cmd *GoUseCmd) WithPlugins(feeder plugins.Feeder) {
	plugs := feeder()
	cmd.plugins = plugs

	for _, plug := range plugs {
		if executor, ok := plug.(sdk.Sdker); ok {
			cmd.executors = append(cmd.executors, executor)
		}
	}
}

func (cmd *GoUseCmd) CreateTaskRunner(pluginNames ...string) (runner taskrunner2.Runner) {
	var tasks taskrunner2.Tasks

	// I am a little slower looping like this, but we don't
	// care to order the plugins afterwards
	for _, pluginName := range pluginNames {
		for _, execPlug := range cmd.executors {
			if execPlug.PluginName() == pluginName {
				tasks = append(tasks, taskrunner2.Task{
					Plugin:      execPlug,
					Description: execPlug.PluginName(),
				})
			}
		}
	}

	runner = taskrunner2.NewTaskRunner(
		taskrunner2.WithTasks(tasks...),
		taskrunner2.WithTitle("use go sdk <version>"),
		taskrunner2.WithTimeout(50*time.Millisecond),
	)
	return runner
}

func (cmd *GoUseCmd) Init() {
	if cmd.path == nil {
		cmd.path = devctlpath.DefaultPather()
	}
	if cmd.fs == nil {
		cmd.fs = vfs.New(osfs.New())
	}
}

func (cmd *GoUseCmd) PluginName() string {
	return "sdk/go/use"
}

func (cmd *GoUseCmd) CmdName() string {
	return "use"
}

func (cmd *GoUseCmd) ExecuteCommand(ctx context.Context, root string, args []string) (err error) {
	cmd.Init()
	version := args[1]
	sdkPath, _ := cmd.fs.EvalSymlinks(cmd.path.SDK("go", version))
	current, _ := cmd.fs.EvalSymlinks(cmd.path.SDK("go", "current"))

	// 1. Clean up existing @current
	fi, err := cmd.fs.Stat(current)
	if err == nil {
		_ = cmd.fs.Remove(current)
		fi.Name()
	}

	// 2. Make sure directories exists
	_ = cmd.fs.MkdirAll(cmd.path.SDK("go"), os.ModePerm)

	// 4. Is the go sdk version installed?
	if exists, _ := cmd.fs.DirExists(sdkPath); !exists {

		// 4 -> Start different plugin do install
		// todo: search plugin and start it
		// todo: ask for user input, if the sdk should be installed
		if err = cmd.fs.MkdirAll(sdkPath, os.ModePerm); err != nil {
			runner := cmd.CreateTaskRunner(GoLinkCmdName, GoInstallCmdName, GoDownloadCmdName)
			err = runner.Run(ctx)
			return errors.Wrapf(err, "TaskRunner execution failed.. ERROR=%v, Tasks=%s", err, runner.Describe())
		}
	}

	// 5. Link the go sdk version to @current
	// ln -s -v -F  /root/sdks/go/1.16.3  /root/sdks/go/current
	// err = cmd.fs.Symlink(sdkPath, current)
	return err
}
