package registry

import (
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)


type SDKAPITest struct {
    body string
    expectedSDKs []string
}
type SDKAPITests []SDKAPITest

func TestGetSDK(t *testing.T) {
    tests := SDKAPITests{
        {
            body:             "[{\"name\":\"dotnet\",\"path\":\"sdk/dotnet\",\"sha\":\"859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"size\":0,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"download_url\":null,\"type\":\"dir\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"html\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\"}},{\"name\":\"java\",\"path\":\"sdk/java\",\"sha\":\"859e4a060e287c06f777da09fbf8fe51dc4afc92\",\"size\":0,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/java?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/java\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc92\",\"download_url\":null,\"type\":\"dir\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/java?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc92\",\"html\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/java\"}}]",
            expectedSDKs: []string{"dotnet", "java"},
        },
        {
            body:             "[{\"name\":\"dotnet\",\"path\":\"sdk/dotnet\",\"sha\":\"859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"size\":0,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"download_url\":null,\"type\":\"dir\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"html\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\"}}]",
            expectedSDKs: []string{"dotnet"},
        },
    }
    tests.Run(t,  SDKApi.GetSDKs)
}




func (test SDKAPITests) Run(t *testing.T, sut interface{}) {
    createTestRegistry := func (test SDKAPITest, t *testing.T, ) ( RegistryAPI,  *httptest.Server)  {
        handler := func(w http.ResponseWriter, r *http.Request) {
            expectedPath := fmt.Sprintf("/content/sdk")
            assert.Equal(t, expectedPath, r.URL.Path)
            w.Header().Set("Content-Type", "application/json")
            _, _ = io.WriteString(w, test.body)
        }
        server := httptest.NewServer(http.HandlerFunc(handler))
        return GithubRegistryApiClient{
            baseUrl: server.URL + "repos/alex-held/dev-env-registry",
        }, server
    }
    action := func(apiTest SDKAPITest) {
        api, ts := createTestRegistry(apiTest, t)
        defer ts.Close()
        switch actualSut := sut.(type) {
        case func(api SDKApi) ([]string, error):
            actual, err := actualSut(api)
            assert.NoError(t, err)
            for _, sdk := range apiTest.expectedSDKs {
                assert.Contains(t, actual, sdk)
            }
        default:
            t.FailNow()
        }
    }

    for _, test := range test {
        action(test)
    }
}
