package golang

import (
	"context"

	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/alex-held/devctl/pkg/plugins"

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

	if exsist, _ := g.fs.DirExists(g.path.SDK("go", version)); !exsist {
		install := &GoInstallCmd{
			path:    g.path,
			runtime: plugins.OSRuntimeInfoGetter{},
			Fs:      g.fs,
		}
		return install.ExecuteCommand(ctx, root, []string{"install", version})
	}

	sdkPath, _ := g.fs.EvalSymlinks(g.path.SDK("go", version))
	current, _ := g.fs.EvalSymlinks(g.path.SDK("go", "current"))

	_ = g.fs.RemoveAll(current)
	_ = g.fs.Remove(current)
	err = g.fs.Symlink(sdkPath, current)
	return err
}
