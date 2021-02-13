package sdkman

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/pkg/aarch"
)

type SDKDownload struct {
	Path   string
	Reader io.Reader
}

// DownloadService downloads SDKs to the filesystem
type DownloadService service

// DownloadSDK downloads the sdk from the sdkman broker.
// SDK specifies the sdk
// Version specifies the apiVersion
// aarch specifies the aarch [darwinx64,darwin]
// https://api.sdkman.io/2/broker/download/scala/scala-2.13.4/darwinx64
// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwinx64
// nolint: lll
func (s *DownloadService) DownloadSDK(ctx context.Context, filepath, sdk, version string, arch aarch.Arch) (*SDKDownload, *http.Response, error) {
	switch arch {
	case aarch.Linux64, aarch.MacOsx, aarch.LinuxArm32:
		return s.downloadSDK(ctx, filepath, sdk, version, string(arch))
	default:
		return nil, nil, fmt.Errorf("right now the provided aarc.Arch '%s' is not supported", arch)
	}
}

// downloadSDK
func (s *DownloadService) downloadSDK(ctx context.Context, filepath, sdk, version, arch string) (
	downloadDl *SDKDownload,
	resp *http.Response,
	err error,
) {
	// CreateDownloadSDK creates the URI to Download SDK archives from SDKMAN API
	// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwinx64
	// "broker", "Download", sdk, version, arch
	req, err := s.client.NewRequest("GET", fmt.Sprintf("broker/download/%s/%s/%s", sdk, version, arch), ctx, nil)

	if err != nil {
		return nil, nil, err
	}

	resp, err = s.client.httpClient.Do(req)
	if err != nil {
		return nil, resp, err
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return nil, nil, err
	}
	if err = resp.Body.Close(); err != nil {
		return nil, nil, err
	}
	body := buf.Bytes()

	err = afero.WriteFile(s.client.fs, filepath, body, fileutil.PrivateDirMode)
	if err != nil {
		return nil, resp, err
	}

	download := &SDKDownload{
		Path:   filepath,
		Reader: bytes.NewBuffer(body),
	}
	return download, resp, err
}