package sdkman

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/alex-held/devctl/internal/system"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/spf13/afero"
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
// system specifies the system [darwinx64,darwin]
// https://api.sdkman.io/2/broker/download/scala/scala-2.13.4/darwinx64
// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwinx64
// nolint: lll,gocognit
func (s *DownloadService) DownloadSDK(ctx context.Context, dlPath, sdk, version string, arch system.Arch) (*SDKDownload, error) {
	switch arch {
	case system.Darwin, system.MacOsx64, system.Linux, system.Linux64, system.LinuxArm32:
		req, err := s.client.NewRequest(ctx, "GET", fmt.Sprintf("broker/download/%s/%s/%s", sdk, version, string(arch)), http.NoBody)
		if err != nil {
			return nil, err
		}
		resp, err := s.client.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		var body bytes.Buffer
		_, err = io.Copy(&body, resp.Body)
		if err != nil {
			return nil, err
		}

		dump := body.Bytes()
		err = afero.WriteFile(s.client.fs, dlPath, dump, fileutil.PrivateFileMode)
		if err != nil {
			return nil, err
		}

		downloadFile, err := s.client.fs.Open(dlPath)
		if err != nil {
			return nil, err
		}
		download := &SDKDownload{
			Path:   dlPath,
			Reader: downloadFile,
		}
		return download, nil

	default:
		return nil, fmt.Errorf("right now the provided aarc.Arch '%s' is not supported", arch)
	}
}
