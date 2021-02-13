package sdkman

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/afero"
)

// Client provides access to the sdkman api
type Client struct {
	context    context.Context
	baseUrl    *url.URL
	client     *http.Client
	httpClient HTTPClient

	// allocate a single struct instead of one for each service
	common service

	// Services used for talking to different parts of the SDKMAN API.
	Download *DownloadService
	ListSdks *ListAllSDKService
	fs       afero.Fs
}

// ClientConfig contains configurable values for the creation of the sdkman.Client
type ClientConfig struct {
	httpClient       *http.Client
	context          context.Context
	fs               afero.Fs
	baseUrl, version string
}

// ClientOption is a function which configures ClientConfig
type ClientOption func(config *ClientConfig) *ClientConfig

// HttpClientOption configures the internal http.Client for the sdkman.Client
func HttpClientOption(client *http.Client) ClientOption {
	return func(c *ClientConfig) *ClientConfig {
		c.httpClient = client
		return c
	}
}

// FileSystemOption configures the afero.Fs used in the sdkman.Client
func FileSystemOption(fs afero.Fs) ClientOption {
	return func(c *ClientConfig) *ClientConfig {
		c.fs = fs
		return c
	}
}

// DefaultSdkManOptions configures the sdkman.Client using defaults
func DefaultSdkManOptions() []ClientOption {
	return []ClientOption{
		HttpClientOption(&http.Client{}),
		FileSystemOption(afero.NewOsFs()),
		UrlOptions(BaseUrl),
	}
}

// SdkManUrlOptions configures the api baseurl
func UrlOptions(baseUrl string) ClientOption {
	return func(c *ClientConfig) *ClientConfig {
		c.baseUrl = baseUrl
		return c
	}
}

// BaseUrl BaseUrl of the remote sdkman api
const BaseUrl = "https://api.sdkman.io"

// NewSdkManClient creates the default *Client using defaults and then the provided options
func NewSdkManClient(options ...ClientOption) *Client {
	config := &ClientConfig{}
	for _, defaultOption := range DefaultSdkManOptions() {
		config = defaultOption(config)
	}
	for _, option := range options {
		config = option(config)
	}

	baseUrl, _ := url.Parse(fmt.Sprintf("%s/%s", config.baseUrl, config.version))

	c := &Client{
		baseUrl:    baseUrl,
		context:    config.context,
		client:     config.httpClient,
		httpClient: http.DefaultClient,
		fs:         config.fs,
	}

	c.common.client = c
	c.Download = (*DownloadService)(&c.common)
	c.ListSdks = (*ListAllSDKService)(&c.common)

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the ClientIn.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, ctx context.Context, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.baseUrl.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.baseUrl)
	}
	u, err := c.baseUrl.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	return req, nil
}
