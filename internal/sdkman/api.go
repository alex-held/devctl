package sdkman

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type sdkmanClient struct {
	baseURL    string
	version    string
	httpClient HTTPClient
}

// NewSdkManClient creates the default sdkman.Client
func NewSdkManClient() Client {
	return &sdkmanClient{
		baseURL:    "https://api.sdkman.io",
		version:    "2",
		httpClient: http.DefaultClient,
	}
}

// ListCandidates lists all available sdk candidates
func (c *sdkmanClient) ListCandidates() (candidates []string, err error) {
	url := c.listCandidatesURL()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url(), http.NoBody)
	if err != nil {
		return nil, err
	}
	result, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	defer result.Body.Close()
	responseBodyBytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	candidatesList := string(responseBodyBytes)
	println(candidatesList)

	candidates = strings.Split(candidatesList, ",")
	return candidates, nil
}

func (c *sdkmanClient) listCandidatesURL() func() string {
	return func() string {
		return c.createURL("candidates", "all")
	}
}

func (c *sdkmanClient) createURL(segments ...string) (uri string) {
	uri = fmt.Sprintf("%s/%s/", c.baseURL, c.version)
	for _, seg := range segments {
		uri = fmt.Sprintf("%s/%s", uri, seg)
	}
	return uri
}
