package registry

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type api interface {
	getUrl(path string, segements ...string) *url.URL
	getContentUrl(paths ...string) *url.URL
	getGitUrl(paths ...string) *url.URL
	getFiles(uri *url.URL) (files []GitHubFile, err error)
}
type VersionAPI interface {
	API
	GetSDKVersionFiles(sdk string) (content []GitHubFile, err error)
	GetSDKVersions(sdk string) (versions []string, err error)
}

func (client GithubRegistryApiClient) GetSDKVersionFiles(sdk string) (files []GitHubFile, err error) {
	var response *http.Response
	var bytes []byte
	uri := client.getContentUrl("sdk", sdk)

	response, err = http.Get(uri.String())
	if err != nil {
		return nil, err
	}
	bytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (client GithubRegistryApiClient) GetSDKVersions(sdk string) (versions []string, err error) {
	var files []GitHubFile

	files, err = client.GetSDKVersionFiles(sdk)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		version := strings.TrimRight(file.Name, ".yaml")
		versions = append(versions, version)
	}
	return versions, nil
}
