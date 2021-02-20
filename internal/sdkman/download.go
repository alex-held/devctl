package sdkman

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/alex-held/devctl/internal/system"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/cli"
)

type SDKDownload struct {
	Path   string
	Reader io.Reader
}

// DownloadService downloads SDKs to the filesystem
type DownloadService service

func (s *DownloadService) Resolve() func(...string) string {
	c := cli.GetOrCreateCLI()
	rootDir := c.ConfigDir()

	return func(paths ...string) (fp string) {
		fp = filepath.Join(rootDir, string(cli.Downloads))

		for _, p := range paths {
			fp = filepath.Join(fp, p)
		}
		return fp
	}
}

// DownloadSDK downloads the sdk from the sdkman broker.
// SDK specifies the sdk
// Version specifies the apiVersion
// system specifies the system [darwinx64,darwin]
// https://api.sdkman.io/2/broker/download/scala/scala-2.13.4/darwinx64
// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwinx64
// nolint: lll,gocognit
func (s *DownloadService) DownloadSDK(ctx context.Context, dlPath, sdk, version string, arch system.Arch) (*SDKDownload, *http.Response, error) {
	switch arch {
	case system.Linux64, system.MacOsx, system.LinuxArm32:
		req, err := s.client.NewRequest(ctx, "GET", fmt.Sprintf("broker/download/%s/%s/%s", sdk, version, string(arch)), http.NoBody)
		if err != nil {
			return nil, nil, err
		}
		resp, err := s.client.client.Do(req)
		if err != nil {
			return nil, resp, err
		}

		var body bytes.Buffer
		_, err = io.Copy(&body, resp.Body)
		if err != nil {
			return nil, nil, err
		}

		dump := body.Bytes()
		err = afero.WriteFile(s.client.fs, dlPath+".zip", dump, fileutil.PrivateFileMode)
		if err != nil {
			return nil, nil, err
		}

		downloadFile, err := s.client.fs.Open(dlPath + ".zip")
		if err != nil {
			return nil, nil, err
		}
		download := &SDKDownload{
			Path:   dlPath + ".zip",
			Reader: downloadFile,
		}
		return download, resp, nil

	default:
		return nil, nil, fmt.Errorf("right now the provided aarc.Arch '%s' is not supported", arch)
	}
}

type FilePathResolver interface {
	Resolve() string
}
