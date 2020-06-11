package registry

import (
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)


type VersionAPITest struct {
    sdk string
    body string
    expectedPath string
    expectedVersions []string
}
type VersionAPITests []VersionAPITest

func TestGetVersionFiles(t *testing.T) {
    tests := VersionAPITests{
        {
            sdk:              "dotnet",
            body:             "[{\"name\":\"3.1.202.yaml\",\"path\":\"sdk/dotnet/3.1.202.yaml\",\"sha\":\"ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"size\":1184,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet/3.1.202.yaml?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/blob/master/sdk/dotnet/3.1.202.yaml\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/blobs/ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"download_url\":\"https://raw.githubusercontent.com/alex-held/dev-env-registry/master/sdk/dotnet/3.1.202.yaml\",\"type\":\"file\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet/3.1.202.yaml?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/blobs/ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"html\":\"https://github.com/alex-held/dev-env-registry/blob/master/sdk/dotnet/3.1.202.yaml\"}}]",
            expectedPath: "/repos/alex-held/dev-env-registry/content/sdk/dotnet",
            expectedVersions: []string{"3.1.202"},
        },
    }
    tests.Run(t,  VersionAPI.GetSDKVersionFiles)
}

func TestGetVersions(t *testing.T) {
    tests := VersionAPITests{
        {
            sdk:              "dotnet",
            body:             "[{\"name\":\"3.1.202.yaml\",\"path\":\"sdk/dotnet/3.1.202.yaml\",\"sha\":\"ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"size\":1184,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet/3.1.202.yaml?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/blob/master/sdk/dotnet/3.1.202.yaml\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/blobs/ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"download_url\":\"https://raw.githubusercontent.com/alex-held/dev-env-registry/master/sdk/dotnet/3.1.202.yaml\",\"type\":\"file\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet/3.1.202.yaml?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/blobs/ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"html\":\"https://github.com/alex-held/dev-env-registry/blob/master/sdk/dotnet/3.1.202.yaml\"}}]",
            expectedPath: "/repos/alex-held/dev-env-registry/content/sdk/dotnet",
            expectedVersions: []string{"3.1.202"},
        },
    }
    tests.Run(t,  VersionAPI.GetSDKVersions)
}



func (test VersionAPITests) Run(t *testing.T, sut interface{}) {
    createTestRegistry := func (tst VersionAPITest, t *testing.T, ) ( RegistryAPI,  *httptest.Server)  {
        handler := func(w http.ResponseWriter, r *http.Request) {
            assert.Equal(t, tst.expectedPath, r.URL.Path)
            w.Header().Set("Content-Type", "application/json")
            _, _ = io.WriteString(w, tst.body)
        }
        server := httptest.NewServer(http.HandlerFunc(handler))
        return GithubRegistryApiClient{
            baseUrl: server.URL,
        }, server
    }
    action := func(apiTest VersionAPITest) {
        api, ts := createTestRegistry(apiTest, t)
        defer ts.Close()
        switch actualSut := sut.(type) {
        case func(api VersionAPI, s string) ([]string, error):
            actual, err := actualSut(api, apiTest.sdk)
            assert.NoError(t, err)
            for _, version := range apiTest.expectedVersions {
                assert.Contains(t, actual, version)
            }
        case func(api VersionAPI , s string) ([]GitHubFile, error):
            actual, err := actualSut(api, apiTest.sdk)
            assert.NoError(t, err)
            var versions []string

            for _, file := range actual {
                versions = append(versions, file.Name)
            }
            for _, version := range apiTest.expectedVersions {
                expectedFileName := version + ".yaml"
                assert.Contains(t, versions, expectedFileName)
            }
        default:
            t.FailNow()
        }
    }

    for _, test := range test {
        action(test)
    }
}
