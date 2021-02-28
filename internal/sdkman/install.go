package sdkman

import (
	"context"
	"fmt"

	"github.com/alex-held/devctl/internal/devctlpath"
	"github.com/alex-held/devctl/internal/system"
)

type InstallService service

func (c *InstallService) Install(ctx context.Context, sdk, version string) (dlPath string, err error) {

	dlPath = devctlpath.DownloadPath(sdk, version, fmt.Sprintf("%s-%s.zip", sdk, version))

	dl, err := c.client.Download.DownloadSDK(ctx, sdk, version, system.DarwinX64)

	fmt.Println(dl.String())

	return dlPath, err
}
