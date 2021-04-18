package golang

import (
	"context"
	"fmt"
	"time"

	"github.com/alex-held/devctl/pkg/ui/taskrunner"
	"github.com/gobuffalo/plugins"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/pkg/errors"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

type GoUseCmd struct {
	plugins plugins.Plugins
	path    devctlpath.Pather
	fs      vfs.VFS
}

func (cmd *GoUseCmd) WithPlugins(feeder plugins.Feeder) { cmd.plugins = feeder() }
func (cmd *GoUseCmd) CreateTaskRunner(version string) (runner taskrunner.Runner) {
	var tasks taskrunner.Tasks

	for _, plug := range cmd.plugins {
		switch t := plug.(type) {
		case *GoDownloadCmd:
			tasks = append(tasks, taskrunner.NewConditionalTask(
				fmt.Sprintf("Download go sdk %s", version),
				func(ctx context.Context) error {
					return t.ExecuteCommand(ctx, "use", []string{version})
				}, func() bool {
					_, err := cmd.fs.Stat(cmd.path.Download(version))
					return err == nil
				},
			))
		case *GoInstallCmd:
			tasks = append(tasks, &taskrunner.ConditionalTask{
				Description: fmt.Sprintf("Install go sdk %s", version),
				Action: func(ctx context.Context) error {
					return t.ExecuteCommand(ctx, "use", []string{version})
				},
				ShouldExecute: func() bool {
					exists, _ := cmd.fs.Exists(cmd.path.SDK("go", version))
					return !exists
				},
			})
		case *GoLinkerCmd:
			tasks = append(tasks, t.AsTasker(version))
		default: // no-op
		}
	}

	runner = taskrunner.NewTaskRunner(
		taskrunner.WithTasks(tasks...),
		taskrunner.WithTitle(fmt.Sprintf("use go sdk %s", version)),
		taskrunner.WithTimeout(50*time.Millisecond),
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
	return GoUseCmdName
}

func (cmd *GoUseCmd) CmdName() string {
	return "use"
}

func (cmd *GoUseCmd) ExecuteCommand(ctx context.Context, _ string, args []string) (err error) {
	version := args[1]
	runner := cmd.CreateTaskRunner(version)
	err = runner.Run(ctx)
	if err != nil {
		return errors.Wrapf(err, "GoUse-TaskRunner execution failed.. ERROR=%v, GoSDKVersion=%s, Tasks=%s", err, version, runner.Describe())
	}

	return nil

	/*
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
				runner := cmd.CreateTaskRunner(version)
				err = runner.Run(ctx)
				return errors.Wrapf(err, "GoUse-TaskRunner execution failed.. ERROR=%v, GoSDKVersion=%s, Tasks=%s", err, version, runner.Describe())
			}
		}

		// 5. Link the go sdk version to @current
		// ln -s -v -F  /root/sdks/go/1.16.3  /root/sdks/go/current
		// err = cmd.fs.Symlink(sdkPath, current)
		return err*/
}
