package sdkman

import (
	"net/http"
	
	"github.com/alex-held/devctl/pkg/aarch"
)

type HTTPDoFunc func(req *http.Request) (*http.Response, error)

// Client provides the SDKMAN Api
type Client interface {
	ListCandidates() (candidates []string, resp *http.Response, err error)
	DownloadSDK(filepath, sdk, version string, arch aarch.Arch) (download *SDKDownload, resp *http.Response, err error)
}

// HTTPClient Sends http.Request returning http.Response or an error.Error
type HTTPClient interface {
	Do(req *http.Request) (response *http.Response, err error)
}
