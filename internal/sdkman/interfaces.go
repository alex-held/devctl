package sdkman

import (
	"net/http"
)

// Client provides the SDKMAN Api
type Client interface {
	ListCandidates() (candidates []string, err error)
}

// HTTPClient Sends http.Request returning http.Response or an error.Error
type HTTPClient interface {
	Do(req *http.Request) (response *http.Response, err error)
}

type HTTPDoFunc func(req *http.Request) (*http.Response, error)
