package registry

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)


func NewRegistryAPI() RegistryAPI  {
    return GithubRegistryApiClient{
        baseUrl: "https://api.github.com/repos/alex-held/dev-env-registry",
    }
}

type GithubRegistryApiClient struct {
    baseUrl     string
}

func (client GithubRegistryApiClient) getUrl(path string, segements ...string) *url.URL {
    for _, p := range segements {
        path += "/" + p
    }
    full := client.baseUrl + path
    uri,err := url.Parse(full)
    if err != nil {
        panic(err)
    }
    return uri
}

func (client GithubRegistryApiClient) getContentUrl(paths ...string) *url.URL {
    return client.getUrl("/content", paths...)
}
func (client GithubRegistryApiClient) getGitUrl(paths ...string) *url.URL {
    return client.getUrl("/git", paths...)
}

func (client GithubRegistryApiClient) GetSDKVersionFiles(sdk string) (files []GitHubFile, err error) {
    var response *http.Response
    var bytes[]byte
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

type GitHubFile struct {
    Name string    `json:"name"`
    Path string    `json:"path"`
    HtmlUrl string    `json:"html_url"`
    GitUrl string    `json:"git_url"`
    DownloadUrl string    `json:"download_url"`
    Type string    `json:"type"`
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

type VersionAPI interface {
    GetSDKVersionFiles(sdk string) (content []GitHubFile, err error)
    GetSDKVersions(sdk string) (versions []string, err error)
}

type RegistryAPI interface {
    VersionAPI
}
