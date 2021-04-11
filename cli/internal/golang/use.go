package golang

import (
	"context"
	"os"

	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/pkg/errors"

	// "github.com/alex-held/devctl/pkg/plugins"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

type GoUseCmd struct {
	path devctlpath.Pather
	fs   vfs.VFS
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
			return errors.Wrapf(err, "DUMMY sdk version installtion dir creation failed. path=%s\n", sdkPath)
		}
	}

	// 5. Link the go sdk version to @current
	// ln -s -v -F  /root/sdks/go/1.16.3  /root/sdks/go/current
	err = cmd.fs.Symlink(sdkPath, current)
	return err
}
