package golang

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/pkg/errors"

	"github.com/alex-held/devctl/pkg/devctlpath"
	"github.com/alex-held/devctl/pkg/plugins"
)

type Renamer func(p string) string

type GoInstallCmd struct {
	path    devctlpath.Pather
	runtime plugins.RuntimeInfoGetter
	Fs      vfs.VFS
}

func (g *GoInstallCmd) PluginName() string {
	return "sdk/go/install"
}

func (g *GoInstallCmd) CmdName() string {
	return "install"
}

func (g *GoInstallCmd) Init() {
	if g.path == nil {
		g.path = devctlpath.DefaultPather()
	}
	if g.Fs == nil {
		g.Fs = vfs.New(osfs.New())
	}
	if g.runtime == nil {
		g.runtime = plugins.OSRuntimeInfoGetter{}
	}
}

func (g *GoInstallCmd) ExecuteCommand(ctx context.Context, root string, args []string) (err error) {
	g.Init()

	version := args[1]
	fmt.Printf("executing: install; version=%s\n", version)
	dlDir := g.path.Download("go", version)
	filename := g.runtime.Get().Format("go%s.[os]-[arch].tar.gz", version)
	archivePath := path.Join(dlDir, filename)
	println(archivePath)

	sdkDir := g.path.SDK("go", version)
	println(sdkDir)

	_ = g.Fs.MkdirAll(dlDir, os.ModePerm)
	if exists, err := g.Fs.Exists(archivePath); !exists || err == nil {
		downloadCmd := &GoDownloadCmd{
			Fs:      g.Fs,
			Pather:  g.path,
			Runtime: g.runtime,
			Output:  NewConsoleOutput(),
			dlOptions: &dlOptions{
				version: version,
				baseURI: "https://golang.org",
			},
		}

		err = downloadCmd.ExecuteCommand(ctx, root, []string{"download", version})
		if err != nil {
			err = errors.Wrapf(err, "failed downloading the gosdk %s\n", version)
			fmt.Printf("%+v", err)
			return err
		}
	}

	_ = g.Fs.MkdirAll(g.path.Download("go", version), os.ModePerm)
	archive, err := g.Fs.OpenFile(archivePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to open go sdk archive=%s\n", archivePath)
	}
	err = g.Fs.MkdirAll(sdkDir, fileutil.PrivateDirMode)
	if err != nil {
		return errors.Wrapf(err, "failed to Extract go sdk %s; dest=%s; archive=%s\n", version, sdkDir, archivePath)
	}
	err = UnTarGzip(archive, sdkDir, GoSDKUnarchiveRenamer(), g.Fs)
	if err != nil {
		return errors.Wrapf(err, "failed to Extract go sdk %s; dest=%s; archive=%s\n", version, sdkDir, archivePath)
	}

	return nil
}

func (p *GoInstallCmd) Link(version string) (err error) {
	return SymLink(p.path, p.Fs, version)
}

//nolint:gocognit
func UnTarGzip(file io.Reader, target string, renamer Renamer, v vfs.VFS) error {
	gr, _ := gzip.NewReader(file)
	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		filename := header.Name
		if renamer != nil {
			filename = renamer(filename)
		}

		p := filepath.Join(target, filename)
		fi := header.FileInfo()

		if fi.IsDir() {
			if e := v.MkdirAll(p, fi.Mode()); e != nil {
				return e
			}
			continue
		}
		file, err := v.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(file, tr)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func GoSDKUnarchiveRenamer() Renamer {
	return func(p string) string {
		parts := strings.Split(p, string(filepath.Separator))
		parts = parts[1:]
		newPath := strings.Join(parts, string(filepath.Separator))
		return newPath
	}
}
