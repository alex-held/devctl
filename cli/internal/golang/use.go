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

func (c *GoUseCmd) Init() {
	if c.path == nil {
		c.path = devctlpath.DefaultPather()
	}
	if c.fs == nil {
		c.fs = vfs.New(osfs.New())
	}
}

func (g *GoUseCmd) PluginName() string {
	return "sdk/go/use"
}

func (g *GoUseCmd) CmdName() string {
	return "use"
}

func (g *GoUseCmd) ExecuteCommand(ctx context.Context, root string, args []string) (err error) {
	g.Init()
	version := args[1]
	sdkPath, _ := g.fs.EvalSymlinks(g.path.SDK("go", version))
	current, _ := g.fs.EvalSymlinks(g.path.SDK("go", "current"))

	// 1. Clean up existing @current
	fi, err := g.fs.Stat(current)
	if err == nil {
		_ = g.fs.Remove(current)
		fi.Name()
	}

	// 2. Make sure directories exists
	_ = g.fs.MkdirAll(g.path.SDK("go"), os.ModePerm)

	// 4. Is the go sdk version installed?
	if exists, _ := g.fs.DirExists(sdkPath); !exists {
		// 4 -> Start different plugin do install
		// todo: search plugin and start it
		// todo: ask for user input, if the sdk should be installed
		if err = g.fs.MkdirAll(sdkPath, os.ModePerm); err != nil {
			return errors.Wrapf(err, "DUMMY sdk version installtion dir creation failed. path=%s\n", sdkPath)
		}
	}

	// 5. Link the go sdk version to @current
	// ln -s -v -F  /root/sdks/go/1.16.3  /root/sdks/go/current
	err = g.fs.Symlink(sdkPath, current)
	return err
}
