package download

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

var _ plugcmd.Namer = &GoDownloadCmd{}
var _ plugins.Plugin = &GoDownloadCmd{}

type GoDownloadCmd struct {
	pluginsFn plugins.Feeder
}

func (l *GoDownloadCmd) CmdName() string {
	return "download"
}

func (l *GoDownloadCmd) PluginName() string {
	return "sdk/go/download"
}

func (l *GoDownloadCmd) Download(ctx context.Context, root string, args []string) error {
	version := args[1]
	fmt.Println(version)
	httpClient := http.Client{}
	baseURI := "https://golang.org"

	fs := afero.NewOsFs()
	dlPath := devctlpath.DownloadPath("go", version)
	err := fs.MkdirAll(dlPath, fileutil.PrivateDirMode)
	if err != nil {
		return errors.Wrapf(err, "failed creating go sdk download path; version=%s; err=%v", version, err)
	}
	filename := fmt.Sprintf("go%s.darwin-amd64.tar.gz", version)
	dlURI := fmt.Sprintf("%s/dl/%s", baseURI, filename)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dlURI, http.NoBody)
	if err != nil {
		return errors.Wrapf(err, "failed downloading go sdk archive; version=%s; err=%v", version, err)
	}
	response, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed downloading go sdk archive; version=%s; err=%v", version, err)
	}
	defer response.Body.Close()
	filePath := path.Join(dlPath, filename)
	err = afero.WriteReader(fs, filePath, response.Body)
	if err != nil {
		return errors.Wrapf(err, "failed writing go sdk archive; version=%s; err=%v", version, err)
	}
	return nil
}
