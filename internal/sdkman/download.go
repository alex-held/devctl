package sdkman

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	
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


// CreateDownloadSDK creates the URI to download SDK archives from SDKMAN API
// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwinx64
func (u *uRLFactory) CreateDownloadSDK(sdk, version, arch string) URI {
	uri := u.createBaseURI()
	uri = uri.Append("broker", "download", sdk, version, arch)
	return uri
}



// DownloadSDK downloads the sdk from the sdkman broker.
// SDK specifies the sdk
// Version specifies the version
// aarch specifies the aarch [darwinx64,darwin]
// https://api.sdkman.io/2/broker/download/scala/scala-2.13.4/darwinx64
// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwinx64
func (s *DownloadService) DownloadSDK(ctx context.Context, filepath, sdk, version string, arch aarch.Arch) (*SDKDownload, *http.Response, error) {
	switch arch {
	case aarch.LINUX_64, aarch.MAC_OSX, aarch.LINUX_ARM32:
		return s.downloadSDK(ctx, filepath, sdk, version, string(arch))
	default:
		return nil, nil, errors.New(fmt.Sprintf("right now the provided aarc.Arch '%s' is not supported", arch))
	}
}



// downloadSDK
func (s *DownloadService) downloadSDK(ctx context.Context, filepath, sdk, version, arch string) (downloadDl *SDKDownload, resp *http.Response, err error) {
	uri := s.client.urlFactory.CreateDownloadSDK(sdk, version, arch)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), http.NoBody)
	
	resp, err = s.client.httpClient.Do(req)
	if err != nil {
		return nil, resp, err
	}
	
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, resp, err
	}
	
	err = afero.WriteFile(s.client.fs, filepath, dump, fileutil.PrivateDirMode)
	if err != nil {
		return nil, resp, err
	}
	
	download := &SDKDownload{
		Path:   filepath,
		Reader: bytes.NewBuffer(dump),
	}
	return download, resp, err
	
}

