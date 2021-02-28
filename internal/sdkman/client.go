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

// BaseUrl BaseUrl of the remote sdkman api
const BaseURL = "https://api.sdkman.io/2"

// Client provides access to the sdkman api
type Client struct {
	baseURL *url.URL
	client  *http.Client

	// allocate a single struct instead of one for each service
	common service

	// Services used for talking to different parts of the SDKMAN API.
	Download *DownloadService
	Registry *RegistryService
	Version  *VersionService
	ListSdks *ListAllSDKService
	fs       afero.Fs
}

// NewSdkManClient creates the default *Client using defaults and then the provided options
func NewSdkManClient(options ...ClientOption) *Client {
	config := &ClientConfig{}
	for _, defaultOption := range DefaultClientOptions() {
		config = defaultOption(config)
	}
	for _, option := range options {
		config = option(config)
	}
	sanitizedBaseURL := strings.TrimSuffix(config.baseURL, "/") + "/"
	baseURL, _ := url.Parse(sanitizedBaseURL)

	c := &Client{
		baseURL: baseURL,
		client:  config.httpClient,
		fs:      config.fs,
	}

	c.common.client = c
	c.Download = (*DownloadService)(&c.common)
	c.ListSdks = (*ListAllSDKService)(&c.common)
	c.Version = (*VersionService)(&c.common)
	c.Registry = (*RegistryService)(&c.common)

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the ClientIn.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.baseURL)
	}
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
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
