package golang

/*
import (
	"context"
	"fmt"

	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/pkg/errors"

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

func (g *GoUseCmd) ExecuteCommand(_ context.Context, _ string, args []string) (err error) {
	g.Init()
	version := args[1]
	sdkPath := g.path.SDK("go", version)
	current := g.path.SDK("go", "current")

	evalSym, err := g.fs.EvalSymlinks(current)
	fmt.Printf("[current] eval symlinks: %s\nERROR: %+v\n", evalSym, err)
	evalSym, err = g.fs.EvalSymlinks(sdkPath)
	fmt.Printf("[sdkPath] eval symlinks: %s\nERROR: %+v\n", evalSym, err)

	fInfo, err := g.fs.Lstat(current)
	if err = g.fs.RemoveAll(current); err != nil {
		err = errors.Wrapf(err, "failed to unlink:\nfile info: %+v\n", fInfo)
	} else {
		_, _ = fmt.Println("Success!!!!\nSuccess!!!!\nSuccess!!!!\n ")
	}

	if err == g.fs.Symlink(sdkPath, current) {
		fmt.Printf("SYMLINK from src=%s -> dest=%s\nERROR: %+v\n", sdkPath, current, err)
		dest, err := g.fs.Readlink(sdkPath)
		if err != nil {
			fmt.Printf("Read link from dest=%s -> src=%s\nERROR: %+v\n", sdkPath, dest, err)
			err = errors.Wrapf(err, "failed -> unlink:\nfile info: %+v\n", fInfo)
		}
		fmt.Printf("Read link from dest=%s -> src=%s", sdkPath, current)
		if err != g.fs.Symlink(dest, current) {
			return errors.Wrapf(err, "failed to reverse the src=%s and dest=%s", sdkPath, current)
		}
	}

	if err = g.fs.Symlink(sdkPath, current); err != nil {
		err = errors.Wrapf(err, "dfailed to create symlink!!!\nSRC=%s -> DEST=%s\n", sdkPath, current)
		fmt.Printf("Err: %+v\nfileInfo: %+v\n", err, fInfo)
		return err
	}
	return nil
}
*/
