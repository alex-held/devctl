package sdkman

import (
	"context"
	"fmt"

	"github.com/alex-held/devctl/internal/devctlpath"
	"github.com/alex-held/devctl/internal/system"
)

type InstallService service

func (c *InstallService) Install(ctx context.Context, sdk, version string) (path string, err error) {
	archive := fmt.Sprintf("%s-%s.zip", sdk, version)
	dlPath := devctlpath.DownloadPath(sdk, version, archive)

	downloadSDK, err := c.client.Download.DownloadSDK(ctx, dlPath, sdk, version, system.DarwinX64)

	fmt.Println(downloadSDK.Path)

	return downloadSDK.Path, err
}
