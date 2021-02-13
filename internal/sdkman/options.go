package sdkman

import (
	"net/http"

	"github.com/spf13/afero"
)

// ClientOption is a function which configures ClientConfig
type ClientOption func(config *ClientConfig) *ClientConfig

// HttpClientOption configures the internal http.Client for the sdkman.Client
func HTTPClientOption(client *http.Client) ClientOption {
	return func(c *ClientConfig) *ClientConfig {
		if client == nil {
			client = http.DefaultClient
		}
		c.httpClient = client
		return c
	}
}

// FileSystemOption configures the afero.Fs used in the sdkman.Client
func FileSystemOption(fs afero.Fs) ClientOption {
	return func(c *ClientConfig) *ClientConfig {
		if fs == nil {
			fs = afero.NewOsFs()
		}
		c.fs = fs
		return c
	}
}

// DefaultClientOptions configures the sdkman.Client using defaults
func DefaultClientOptions() []ClientOption {
	return []ClientOption{
		HTTPClientOption(&http.Client{}),
		FileSystemOption(afero.NewOsFs()),
		URLOptions(BaseURL),
	}
}

// SdkManUrlOptions configures the api baseurl
func URLOptions(baseURL string) ClientOption {
	return func(c *ClientConfig) *ClientConfig {
		c.baseURL = baseURL
		return c
	}
}
