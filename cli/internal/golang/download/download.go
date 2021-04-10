package download

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/plugins/downloader"
	plugins2 "github.com/alex-held/devctl/pkg/plugins"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

var _ plugcmd.Namer = &GoDownloadCmd{}
var _ plugins.Plugin = &GoDownloadCmd{}

type GoDownloadCmd struct {
	*goSDKCore
	*dlOptions
	pluginsFn plugins.Feeder
}

func (cmd *GoDownloadCmd) CmdName() string {
	return "download"
}

func (cmd *GoDownloadCmd) PluginName() string {
	return "sdk/go/download"
}

type dlOptions struct {
	version string
	baseURI string
}

type goSDKCore struct {
	sdk               string
	runtimeInfoGetter plugins2.OSRuntimeInfoGetter
	fs                afero.Fs
	pather            devctlpath.Pather
	out               *output
}

func (cmd *GoDownloadCmd) DownloadDir() string {
	return devctlpath.DownloadPath(cmd.sdk, cmd.version)
}

func (cmd *GoDownloadCmd) DownloadUri() (uri string) {
	artifact := cmd.DownloadArtifactName()
	uri = cmd.runtimeInfoGetter.Format("%s/dl/%s", cmd.baseURI, artifact)
	return uri
}

func (cmd *GoDownloadCmd) DownloadArtifactPath() string {
	return path.Join(cmd.DownloadDir(), cmd.DownloadArtifactName())
}

func (cmd *GoDownloadCmd) DownloadProgressDesc() string {
	return fmt.Sprintf("downloading sdk: %s %s", cmd.sdk, cmd.version)
}

func (cmd *GoDownloadCmd) DownloadArtifactName() string {
	artifactName := cmd.runtimeInfoGetter.Format("go%s.[os]-[arch].tar.gz", cmd.dlOptions.version)
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

func NewConsoleOutput() (out *output) {
	out = NewOutput(os.Stdout, os.Stderr, os.Stdin)
	return out
}

func (cmd *GoDownloadCmd) artifactCached() (exists bool) {
	downloadArtifactPath := cmd.DownloadArtifactPath()
	_, err := cmd.fs.Stat(downloadArtifactPath)
	return err == nil
}

func (cmd *GoDownloadCmd) ExecuteCommand(ctx context.Context, root string, args []string) error {
	err := cmd.Init(ctx, root, args[0:])
	if err != nil {
		return errors.Wrapf(err, "failed to initialize %T", *cmd)
	}

	if cmd.artifactCached() {
		_, err := cmd.Out().Write([]byte(fmt.Sprintf("go sdk %s already downloaded\n", cmd.version)))
		return err
	}

	err = cmd.fs.MkdirAll(cmd.DownloadDir(), fileutil.PrivateDirMode)
	if err != nil {
		return errors.Wrapf(err, "failed creating go sdk download path; version=%s", cmd.dlOptions.version)
	}
	artifactFile, err := cmd.fs.Create(cmd.DownloadArtifactPath())
	if err != nil {
		return errors.Wrapf(err, "failed creating / opening file handle for the download")
	}

	dl := downloader.NewDownloader(cmd.DownloadUri(), cmd.DownloadProgressDesc(), artifactFile, cmd.Out())
	err = dl.Download(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed downloading go sdk %v from the remote server %s", cmd.version, cmd.baseURI)
	}
	return nil
}

func (cmd *GoDownloadCmd) Init(_ context.Context, _ string, args []string) error {
	cmd.goSDKCore = DefaultSDKCore()
	cmd.dlOptions = &dlOptions{
		version: args[1],
		baseURI: "https://golang.org",
	}
	return nil
}

func (cmd *GoDownloadCmd) Out() io.Writer {
	return cmd.out.out
}

func DefaultSDKCore() *goSDKCore {
	return &goSDKCore{
		runtimeInfoGetter: plugins2.OSRuntimeInfoGetter{},
		fs:                afero.NewOsFs(),
		pather:            devctlpath.DefaultPather(),
		out:               NewConsoleOutput(),
		sdk:               "go",
	}
}
