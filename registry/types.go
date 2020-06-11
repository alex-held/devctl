package registry

import (
    "net/url"
)

type GithubRegistryApiClient struct {
    baseUrl     string
}

type GitHubFile struct {
    Name string    `json:"name"`
    Path string    `json:"path"`
    HtmlUrl string    `json:"html_url"`
    GitUrl string    `json:"git_url"`
    DownloadUrl string    `json:"download_url"`
    Type string    `json:"type"`
}


type RegistryAPI interface {
    VersionAPI
}


func NewRegistryAPI() RegistryAPI  {
    return GithubRegistryApiClient{
        baseUrl: "https://api.github.com/repos/alex-held/dev-env-registry",
    }
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
