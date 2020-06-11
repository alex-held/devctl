package registry

import (
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func NewTestServerRegistryAPI(handler func(w http.ResponseWriter, r *http.Request)) (api RegistryAPI, server *httptest.Server)  {
    server = httptest.NewServer(http.HandlerFunc(handler))
    return GithubRegistryApiClient{
        baseUrl: server.URL,
    }, server
}


type TestRequestResponse struct {
    path, query       string    // request
    body string                 // response
}


func TestGetVersionFiles(t *testing.T) {
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _, _ = io.WriteString(w, "[{\"name\":\"3.1.202.yaml\",\"path\":\"sdk/dotnet/3.1.202.yaml\",\"sha\":\"ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"size\":1184,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet/3.1.202.yaml?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/blob/master/sdk/dotnet/3.1.202.yaml\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/blobs/ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"download_url\":\"https://raw.githubusercontent.com/alex-held/dev-env-registry/master/sdk/dotnet/3.1.202.yaml\",\"type\":\"file\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet/3.1.202.yaml?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/blobs/ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"html\":\"https://github.com/alex-held/dev-env-registry/blob/master/sdk/dotnet/3.1.202.yaml\"}}]")
    }
    api, ts := NewTestServerRegistryAPI(handler)
    defer ts.Close()

    versions, err := api.GetSDKVersionFiles("dotnet")
    assert.NoError(t, err)
    assert.Contains(t, versions[0].Name,"3.1.202")
}



func TestGetVersions(t *testing.T) {
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _, _ = io.WriteString(w, "[{\"name\":\"3.1.202.yaml\",\"path\":\"sdk/dotnet/3.1.202.yaml\",\"sha\":\"ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"size\":1184,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet/3.1.202.yaml?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/blob/master/sdk/dotnet/3.1.202.yaml\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/blobs/ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"download_url\":\"https://raw.githubusercontent.com/alex-held/dev-env-registry/master/sdk/dotnet/3.1.202.yaml\",\"type\":\"file\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet/3.1.202.yaml?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/blobs/ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"html\":\"https://github.com/alex-held/dev-env-registry/blob/master/sdk/dotnet/3.1.202.yaml\"}}]")
    }
    api, ts := NewTestServerRegistryAPI(handler)
    defer ts.Close()

    versions, err := api.GetSDKVersions("dotnet")
    assert.NoError(t, err)
    assert.Contains(t, versions,"3.1.202")
}
