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
// arch specifies the arch [darwinx64,darwin]
// https://api.sdkman.io/2/broker/download/scala/scala-2.13.4/darwinx64
// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwinx64
// nolint: lll
func (s *DownloadService) DownloadSDK(ctx context.Context, filepath, sdk, version string, arch aarch.Arch) (*SDKDownload, *http.Response, error) {
	switch arch {
	case aarch.Linux64, aarch.MacOsx, aarch.LinuxArm32:
		// CreateDownloadSDK creates the URI to Download SDK archives from SDKMAN API
		// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwinx64
		// "broker", "Download", sdk, version, arch
		req, err := s.client.NewRequest(ctx, "GET", fmt.Sprintf("broker/download/%s/%s/%s", sdk, version, arch), nil)

		if err != nil {
			return nil, nil, err
		}

		resp, err := s.client.httpClient.Do(req)
		if err != nil {
			return nil, resp, err
		}

		var buf bytes.Buffer
		if _, err = buf.ReadFrom(resp.Body); err != nil {
			return nil, nil, err
		}

		defer resp.Body.Close()

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
	default:
		return nil, nil, fmt.Errorf("right now the provided aarc.Arch '%s' is not supported", arch)
	}
}
