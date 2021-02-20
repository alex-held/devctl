package sdkman

import (
	"context"
	"fmt"

	"github.com/alex-held/devctl/internal/system"
)

type InstallService service

func (c *InstallService) Install(ctx context.Context, sdk, version string, arch system.Arch) (path string, err error) {
	dlPath := ""
	downloadSDK, err := c.client.Download.DownloadSDK(ctx, dlPath, sdk, version, arch)

	fmt.Println(downloadSDK.Path)

	return downloadSDK.Path, err
}