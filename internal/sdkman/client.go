package sdkman

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"strings"
	
	"github.com/spf13/afero"
	
	"github.com/alex-held/devctl/pkg/aarch"
)

type URI fmt.Stringer

type uri struct {
	scheme      string
	host        string
	segments    []string
	queryString []string
}

func (u *uri) String() string {
	return u.Stringer()
}

func (u *uri) Append(segments ...string) (uri *uri) {
	for _, s := range segments {
		u.segments = append(u.segments, s)
	}
	return u
}

func (u uri) Stringer() string {
	path := func() string {
		if len(u.segments) >= 1 {
			return "/" + strings.Join(u.segments, "/")
		}
		return ""
	}()
	
	query := func() string {
		if len(u.queryString) >= 1 {
			return "?" + strings.Join(u.segments, "&")
		}
		return ""
	}()
	
	unsafeString := fmt.Sprintf("%s://%s%s%s", u.scheme, u.host, path, query)
	escapedString := html.EscapeString(unsafeString)
	return escapedString
}

type APIURLFactoryFunc func() (uri string, err error)

type APIURLFactory interface {
	DownloadSDK() APIURLFactoryFunc
	ListAllSDKURI() APIURLFactoryFunc
}

func (u *uRLFactory) createBaseURI() *uri {
	host := fmt.Sprintf("%s/%s", u.hostname, u.version)
	return &uri{
		scheme:      "https",
		host:        host,
		segments:    []string{},
		queryString: []string{},
	}
}

type uRLFactory struct {
	hostname string
	version  string
}

type sdkmanClient struct {
	context    context.Context
	urlFactory uRLFactory
	httpClient HTTPClient
	
	// allocate a single struct instead of one for each service
	common            service
	
	// Services used for talking to different parts of the SDKMAN API.
	download   *DownloadService
	sdkService *ListAllSDKService
	fs         afero.Fs
}

func (s *sdkmanClient) ListCandidates() (candidates []string, resp *http.Response, err error) {
	return s.sdkService.ListAllSDK(s.context)
}

func (s *sdkmanClient) DownloadSDK(filepath, sdk, version string, arch aarch.Arch) (download *SDKDownload, resp *http.Response, err error) {
	return s.download.DownloadSDK(s.context, filepath, sdk, version, arch)
}

// NewSdkManClient creates the default sdkman.Client
func NewSdkManClient() Client {
	return &sdkmanClient{
		urlFactory: uRLFactory{
			hostname: "https://api.sdkman.io",
			version:  "2",
		},
		httpClient: http.DefaultClient,
		download:   nil,
		sdkService: nil,
		fs:         afero.NewOsFs(),
	}
}
