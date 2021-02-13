package sdkman

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
)

type ListAllSDKService service

// CreateListAllAvailableSDKURI gets all available SDK and returns them as an array of strings
// https://api.sdkman.io/2/candidates/all
func (s *ListAllSDKService) ListAllSDK(ctx context.Context) (candidates []string, resp *http.Response, err error) {
	// CreateListAllAvailableSDKURI creates the URI to list all available SDK
	// https://api.sdkman.io/2/candidates/all
	req, err := s.client.NewRequest(ctx, "GET", "candidates/all", nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err = s.client.client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	candidatesList := string(responseBodyBytes)
	candidates = strings.Split(candidatesList, ",")
	return candidates, resp, nil
}
