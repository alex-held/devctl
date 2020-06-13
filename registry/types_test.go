package registry

import (
	"net/http/httptest"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGetUrlReturnsCorrectUrlWhenDefaultAPI(t *testing.T) {
	api := NewRegistryAPI()
	client := api.(GithubRegistryApiClient)
	expected := "https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk"
	actual := client.getUrl("contents", "sdk")
	assert.Equal(t, expected, actual.String())
}

func TestGetGitUrlReturnsCorrectUrlWhenDefaultAPI(t *testing.T) {
	api := NewRegistryAPI()
	client := api.(GithubRegistryApiClient)
	expected := "https://api.github.com/repos/alex-held/dev-env-registry/git/sdk"
	actual := client.getGitUrl("sdk")
	assert.Equal(t, expected, actual.String())
}

func TestGetContentsUrlReturnsCorrectUrlWhenDefaultAPI(t *testing.T) {
	api := NewRegistryAPI()
	client := api.(GithubRegistryApiClient)
	expected := "https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet"
	actual := client.getContentUrl("sdk", "dotnet")
	assert.Equal(t, expected, actual.String())
}

type TestAPI interface {
	VersionAPI
	SDKApi
	API
	Close()
}

type TestRegistryApiClient struct {
	API
	VersionAPI
	SDKApi
	client GithubRegistryApiClient
	server *httptest.Server
}

func (testClient TestRegistryApiClient) Close() {
	testClient.server.Close()
}
