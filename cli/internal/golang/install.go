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

	"github.com/alex-held/devctl/pkg/devctlpath"
	"github.com/alex-held/devctl/pkg/plugins"
	"github.com/alex-held/devctl/pkg/system"
	"github.com/alex-held/devctl/pkg/ui/taskrunner"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/pkg/errors"
)

var _ plugins.Executor = &GoInstallCmd{}

type Renamer func(p string) string

type GoInstallCmd struct {
	path    devctlpath.Pather
	runtime system.RuntimeInfoGetter
	Fs      vfs.VFS
}

func (cmd *GoInstallCmd) AsTasker(version string) taskrunner.Tasker {
	archivePath := cmd.path.Download("go", version)
	installPath := cmd.path.SDK("go", version)

	return &taskrunner.ConditionalTask{
		Description: "installing go sdk %s into the go sdk directory",
		Action: func(ctx context.Context) error {
			archive, err := cmd.Fs.OpenFile(archivePath, os.O_WRONLY, os.ModePerm)
			if err != nil {
				return errors.Wrapf(err, "failed to open go sdk archive=%s\n", archivePath)
			}
			err = cmd.Fs.MkdirAll(installPath, os.ModePerm)
			if err != nil {
				return errors.Wrapf(err, "failed to Extract go sdk %s; dest=%s; archive=%s\n", version, installPath, archivePath)
			}
			err = UnTarGzip(archive, installPath, GoSDKUnarchiveRenamer(), cmd.Fs)
			if err != nil {
				return errors.Wrapf(err, "failed to Extract go sdk %s; dest=%s; archive=%s\n", version, installPath, archivePath)
			}
			return nil
		},
		ShouldExecute: func() bool {
			// dont run installer if the version is already installed
			exists, _ := cmd.Fs.Exists(installPath)
			return !exists
		},
	}
}

/*


tasks = append(tasks,
taskrunner.NewConditionalTask(
"Make sure the go-sdk directory exists",
func(ctx context.Context) error {
	return cmd.fs.MkdirAll(cmd.path.SDK("go"), os.ModePerm)
},
func() bool {return true},
),
)

*//*


tasks = append(tasks,
taskrunner.NewConditionalTask(
"Make sure the go-sdk directory exists",
func(ctx context.Context) error {
	return cmd.fs.MkdirAll(cmd.path.SDK("go"), os.ModePerm)
},
func() bool {return true},
),
)

*/

func (cmd *GoInstallCmd) PluginName() string {
	return GoInstallCmdName
}

func (cmd *GoInstallCmd) CmdName() string {
	return "install"
}

func (cmd *GoInstallCmd) Init() {
	if cmd.path == nil {
		cmd.path = devctlpath.DefaultPather()
	}
	if cmd.Fs == nil {
		cmd.Fs = vfs.New(osfs.New())
	}
	if cmd.runtime == nil {
		cmd.runtime = system.OSRuntimeInfoGetter{}
	}
}

func (cmd *GoInstallCmd) ExecuteCommand(ctx context.Context, root string, args []string) (err error) {

	return cmd.Fs.MkdirAll(cmd.path.SDK("go", "1.16.3"), os.ModePerm)

	//	cmd.Init()

	version := args[1]
	fmt.Printf("executing: install; version=%s\n", version)
	dlDir := cmd.path.Download("go", version)
	filename := cmd.runtime.Get().Format("go%s.[os]-[arch].tar.gz", version)
	archivePath := path.Join(dlDir, filename)
	println(archivePath)

	sdkDir := cmd.path.SDK("go", version)
	println(sdkDir)

	_ = cmd.Fs.MkdirAll(dlDir, os.ModePerm)
	if exists, err := cmd.Fs.Exists(archivePath); !exists || err == nil { //nolint:govet
		/*	downloadCmd := &GoDownloadCmd{
				Fs:      cmd.Fs,
				Pather:  cmd.path,
				Runtime: cmd.runtime,
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
			}*/
	}

	_ = cmd.Fs.MkdirAll(cmd.path.Download("go", version), os.ModePerm)
	archive, err := cmd.Fs.OpenFile(archivePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to open go sdk archive=%s\n", archivePath)
	}
	err = cmd.Fs.MkdirAll(sdkDir, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to Extract go sdk %s; dest=%s; archive=%s\n", version, sdkDir, archivePath)
	}
	err = UnTarGzip(archive, sdkDir, GoSDKUnarchiveRenamer(), cmd.Fs)
	if err != nil {
		return errors.Wrapf(err, "failed to Extract go sdk %s; dest=%s; archive=%s\n", version, sdkDir, archivePath)
	}

	return nil
}

func (cmd *GoInstallCmd) Link(version string) (err error) {
	return SymLink(cmd.path, cmd.Fs, version)
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
