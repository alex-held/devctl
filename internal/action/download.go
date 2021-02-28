package action

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/system"
)

type Download action

func (i *Download) Download(ctx context.Context, sdk, version string) (archive afero.File, err error) {
	dl, err := i.Client.Download.DownloadSDK(ctx, sdk, version, system.GetCurrent())
	if err != nil {
		return nil, errors.Wrap(err, "error downloading sdk from api.sdkman.io")
	}
	archivePath := i.Pather.Download(sdk, version, fmt.Sprintf("%s-%s.zip", sdk, version))
	archive, err = saveArchive(i.Fs, dl.Buffer, archivePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to save http content to zip file; path=%s\n", archivePath)
	}

	return archive, nil
}
