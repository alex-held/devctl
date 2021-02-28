package sdkman

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/alex-held/devctl/internal/system"
)

type SDKDownload struct {
	bytes.Buffer
}


// DownloadService downloads SDKs to the filesystem
type DownloadService service

// DownloadSDK downloads the sdk from the sdkman broker.
// SDK specifies the sdk
// Version specifies the apiVersion
// system specifies the system [darwinx64,darwin]
// https://api.sdkman.io/2/broker/download/scala/scala-2.13.4/darwinx64
// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwinx64
// nolint: lll,gocognit
func (s *DownloadService) DownloadSDK(ctx context.Context, sdk, version string, arch system.Arch) (dl *SDKDownload, err error) {
	switch {
	case arch.IsLinux() || arch.IsDarwin():
		req, err := s.client.NewRequest(ctx, "GET", fmt.Sprintf("broker/download/%s/%s/%s", sdk, version, string(arch)), nil)
		if err != nil {
			return nil, err
		}
		resp, err := s.client.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		dl := &SDKDownload{}
		_, err = io.Copy(dl, resp.Body)
		if err != nil {
			return nil, err
		}

		return dl, nil

	default:
		return nil, fmt.Errorf("right now the provided aarc.Arch '%s' is not supported", arch)
	}
}
