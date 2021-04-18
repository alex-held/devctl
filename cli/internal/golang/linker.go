package golang

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/alex-held/devctl/pkg/ui/taskrunner"
	"github.com/gobuffalo/plugins"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/pkg/errors"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

type GoLinkerCmd struct {
	plugins plugins.Plugins
	path    devctlpath.Pather
	fs      vfs.VFS
}

type Taskable interface {
	AsTasker(version string) taskrunner.Tasker
}

func (cmd *GoLinkerCmd) RemoveBrokenSymlinks(directory string) (err error) {
	fileInfos, err := cmd.fs.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, fi := range fileInfos {
		filepath := path.Join(directory, fi.Name())
		_ = tryRemoveSymlink(cmd.fs, filepath)
	}
	return nil
}

func tryRemoveSymlink(vfs vfs.VFS, name string) error {
	fi, err := vfs.Lstat(name)
	if err != nil && os.IsNotExist(err) {
		return nil
	}

	if fi.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("%s: not a symlink. Inspect and remove it manually to continue", name)
	}

	if err = vfs.Remove(name); err != nil {
		return err
	}

	return nil
}

func (cmd *GoLinkerCmd) AsTasker(version string) taskrunner.Tasker {
	return &taskrunner.SimpleTask{
		Desc: fmt.Sprintf("Linking go sdk %s to @current", version),
		Action: func(ctx context.Context) error {
			sdkPath := cmd.path.SDK("go", version)
			currentPath := cmd.path.SDK("go", "current")
			err := cmd.RemoveBrokenSymlinks(cmd.path.SDK("go"))
			if err != nil {
				return err
			}
		//		err = cmd.fs.Remove(currentPath)
			_ = err // ingore the error.. rm -f
			return cmd.fs.Symlink(sdkPath, currentPath)
		},
	}
}

func (cmd *GoLinkerCmd) WithPlugins(feeder plugins.Feeder) { cmd.plugins = feeder() }

func (cmd *GoLinkerCmd) CreateTaskRunner(version string) (runner taskrunner.Runner) {
	/*sdkPath := cmd.path.SDK("go", version)
	current := cmd.path.SDK("go", "current")*/
	/*
		removeCurrentIfExitstsTask := taskrunner.NewConditionalTask(
			"Clean up existing @current",
			func(ctx context.Context) error {
				return cmd.fs.RemoveAll(current)
			},
			func() bool {
				fi, err := cmd.fs.Stat(current)
				fi, err = cmd.fs.Lstat(current)
				_ = fi
				return err != nil
			},
		)

		linkVersionToCurrent := taskrunner.NewConditionalTask(
			fmt.Sprintf("Link %s to @current", version),
			func(ctx context.Context) error {
				err := cmd.fs.Symlink(sdkPath, current)
				return err
			},
			func() bool {
				e, _ := cmd.fs.DirExists(sdkPath)
				return e
			},
		)*/

	tasker := cmd.AsTasker(version)
	runner = taskrunner.NewTaskRunner(
		taskrunner.WithTitle(fmt.Sprintf("link go sdk %s to @current", version)),
		taskrunner.WithTasks(
			tasker,
			/*	removeCurrentIfExitstsTask,
				linkVersionToCurrent),*/
		),
	)

	return runner
}

func (cmd *GoLinkerCmd) PluginName() string {
	return GoLinkCmdName
}

func (cmd *GoLinkerCmd) CmdName() string {
	return "link"
}

func (cmd *GoLinkerCmd) ExecuteCommand(ctx context.Context, root string, args []string) (err error) {
	version := args[1]
	/*	runner := cmd.CreateTaskRunner(version)
		err = runner.Run(ctx)
		if err != nil {
			return errors.Wrapf(err, "GoUse-TaskRunner execution failed.. ERROR=%v, GoSDKVersion=%s, Tasks=%s", err, version, runner.Describe())
		}
	*/

	sdkPath := cmd.path.SDK("go", version)
	current := cmd.path.SDK("go", "current")

	removeCurrentIfExitstsTask := taskrunner.NewConditionalTask(
		"Clean up existing @current",
		func(ctx context.Context) error {
			return cmd.fs.RemoveAll(current)
		},
		func() bool {
			fi, err := cmd.fs.Stat(current)
			fi, err = cmd.fs.Lstat(current)
			_ = fi
			return err != nil
		},
	)

	linkVersionToCurrent := taskrunner.NewConditionalTask(
		fmt.Sprintf("Link %s to @current", version),
		func(ctx context.Context) error {
			err := cmd.fs.Symlink(sdkPath, current)
			return err
		},
		func() bool {
			e, _ := cmd.fs.DirExists(sdkPath)
			return e
		},
	)

	err = removeCurrentIfExitstsTask.Task(ctx)
	if err != nil {
		return errors.Wrapf(err, "Unable to remove @current before linking.. ABORT\n")
	}

	err = linkVersionToCurrent.Task(ctx)
	if err != nil {
		return errors.Wrapf(err, "Unable to link %s to @current", version)
	}

	return nil
}
