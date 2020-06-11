package registry

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

const basePath = "repos/alex-held/dev-env-registry"
const githubHost = "https://api.github.com"

type GithubRegistryApiClient struct {
    rootPath string
    baseUrl  string
}

type GitHubFile struct {
    Name        string `json:"name"`
    Path        string `json:"path"`
    HtmlUrl     string `json:"html_url"`
    GitUrl      string `json:"git_url"`
    DownloadUrl string `json:"download_url"`
    Type        string `json:"type"`
}

type RegistryAPI interface {
    VersionAPI
    SDKApi
}

func NewRegistryAPI() RegistryAPI {
    return GithubRegistryApiClient{
        rootPath: basePath,
        baseUrl:  githubHost,
    }
}


func (client GithubRegistryApiClient) getUrl(path string, segements ...string) *url.URL {
    for _, p := range segements {
        path += "/" + p
    }
    path = strings.Trim(path, "/")
    full := fmt.Sprintf("%s/%s/%s", client.baseUrl, client.rootPath, path)
    fmt.Println(full)
    uri, err := url.Parse(full)
    if err != nil {
        panic(err)
    }
    return uri
}

func (client GithubRegistryApiClient) getContentUrl(paths ...string) *url.URL {
    return client.getUrl("/contents", paths...)
}
func (client GithubRegistryApiClient) getGitUrl(paths ...string) *url.URL {
    return client.getUrl("/git", paths...)
}

func (client GithubRegistryApiClient) getFiles(uri *url.URL) (files []GitHubFile, err error) {
    response, err := http.Get(uri.String())
    if err != nil {
        return files, err
    }
    bytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return files, err
    }
    err = json.Unmarshal(bytes, &files)
    if err != nil {
        return files, err
    }
    return files, nil
}
