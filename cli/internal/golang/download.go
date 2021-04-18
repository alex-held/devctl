package golang

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	plugins2 "github.com/alex-held/devctl/pkg/plugins"
	downloader2 "github.com/alex-held/devctl/pkg/plugins/downloader"
	"github.com/alex-held/devctl/pkg/system"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/pkg/errors"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

var _ plugcmd.Namer = &GoDownloadCmd{}
var _ plugins.Plugin = &GoDownloadCmd{}
var _ plugins2.Executor = &GoDownloadCmd{}

type GoDownloadCmd struct {
	Fs      vfs.VFS
	Pather  devctlpath.Pather
	Runtime system.RuntimeInfoGetter
	*dlOptions
	Output *output
}

func (cmd *GoDownloadCmd) CmdName() string {
	return "download"
}

func (cmd *GoDownloadCmd) PluginName() string {
	return GoDownloadCmdName
}

type dlOptions struct {
	version string
	baseURI string
}

func (cmd *GoDownloadCmd) DownloadDir() string {
	return devctlpath.DownloadPath("go", cmd.version)
}

func (cmd *GoDownloadCmd) downloadURI() (uri string) {
	artifact := cmd.DownloadArtifactName()
	uri = cmd.Runtime.Get().Format("%s/dl/%s", cmd.baseURI, artifact)
	return uri
}

func (cmd *GoDownloadCmd) DownloadArtifactPath() string {
	return path.Join(cmd.DownloadDir(), cmd.DownloadArtifactName())
}

func (cmd *GoDownloadCmd) DownloadProgressDesc() string {
	return fmt.Sprintf("downloading sdk: %s %s", "go", cmd.version)
}

func (cmd *GoDownloadCmd) DownloadArtifactName() string {
	artifactName := cmd.Runtime.Get().Format("go%s.[os]-[arch].tar.gz", cmd.dlOptions.version)
	return artifactName
}

type output struct {
	out io.Writer
	err io.Writer
	in  io.Reader
}

type nopReader struct{}

func (r *nopReader) Read(_ []byte) (n int, err error) {
	return 0, nil
}

func NewOutput(out, err io.Writer, in io.Reader) (o *output) {
	if out == nil {
		out = ioutil.Discard
	}
	if err == nil {
		err = ioutil.Discard
	}
	if in == nil {
		in = &nopReader{}
	}
	o = &output{
		out: out,
		err: err,
		in:  in,
	}
	return o
}

func NewConsoleOutput() *output {
	return NewOutput(os.Stdout, os.Stderr, os.Stdin)
}

func (cmd *GoDownloadCmd) artifactCached() (exists bool) {
	downloadArtifactPath := cmd.DownloadArtifactPath()
	_, err := cmd.Fs.Stat(downloadArtifactPath)
	return err == nil
}

func (cmd *GoDownloadCmd) ExecuteCommand(ctx context.Context, root string, args []string) error {
	cmd.Init()
	version := args[1]
	cmd.version = version

	if cmd.artifactCached() {
		cmd.Out().Write([]byte(fmt.Sprintf("go sdk %s already downloaded\n", version)))
		return nil
	}

	err := cmd.Fs.MkdirAll(cmd.DownloadDir(), fileutil.PrivateDirMode)
	if err != nil {
		return errors.Wrapf(err, "failed creating go sdk download path; version=%s", version)
	}
	artifactFile, err := cmd.Fs.Create(cmd.DownloadArtifactPath())
	if err != nil {
		return errors.Wrapf(err, "failed creating / opening file handle for the download")
	}

	dl := downloader2.NewDownloader(cmd.downloadURI(), cmd.DownloadProgressDesc(), artifactFile, cmd.Out())
	err = dl.Download(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed downloading go sdk %v from the remote server %s", version, cmd.baseURI)
	}
	return nil
}

func (cmd *GoDownloadCmd) Init() {
	if cmd == nil {
		cmd = &GoDownloadCmd{}
	}
	cmd.Fs = vfs.New(osfs.New())
	cmd.Runtime = system.OSRuntimeInfoGetter{}
	cmd.Pather = devctlpath.DefaultPather()
	cmd.Output = NewConsoleOutput()
	cmd.dlOptions = &dlOptions{
		baseURI: "https://golang.org",
	}
}

func (cmd *GoDownloadCmd) Out() io.Writer {
	return cmd.Output.out
}
