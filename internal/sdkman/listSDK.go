package sdkman

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
)
type ListAllSDKService service

// CreateListAllAvailableSDKURI creates the URI to list all available SDK
// https://api.sdkman.io/2/candidates/all
func (u *uRLFactory) CreateListAllAvailableSDKURI() URI {
	return u.createBaseURI().Append("candidates", "all")
}

// CreateListAllAvailableSDKURI gets all available SDK and returns them as an array of strings
// https://api.sdkman.io/2/candidates/all
func (s *ListAllSDKService) ListAllSDK(ctx context.Context) (candidates []string, resp *http.Response, err error) {
	uri := s.client.urlFactory.CreateListAllAvailableSDKURI()
	urlString := uri.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlString, http.NoBody)
	
	resp, err = s.client.httpClient.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	
	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	
	candidatesList := string(responseBodyBytes)
	println(candidatesList)
	
	candidates = strings.Split(candidatesList, ",")
	return candidates, resp, nil
	
}
