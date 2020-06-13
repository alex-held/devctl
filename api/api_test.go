package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/github"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alex-held/dev-env/shared"
)

func init() {
	shared.BootstrapLogger(zerolog.TraceLevel)
}

type Test struct {
	InputPath      string
	ExpectedResult []interface{}
	ExpectedPath   string `json:"expected_path"`
	Body           string `json:"Body"`
}

func (t Test) MarshalZerologObject(e *zerolog.Event) {
	e.Str("inputPath", t.InputPath)
	arr := zerolog.Arr()
	for _, expected := range t.ExpectedResult {
		arr.Interface(expected)
	}
	e.Array("expected result", arr)
	e.Str("expectedPath", t.ExpectedPath)
	e.Interface("expectedResult", t.ExpectedPath)
}

func TestGetPKGs(t *testing.T) {
	tests := []Test{
		{
			"sdk", []interface{}{"dotnet"}, "/repos/alex-held/dev-env-registry/contents/sdk",
			"[{\"name\":\"dotnet\",\"path\":\"sdk/dotnet\",\"sha\":\"859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"size\":0,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"download_url\":null,\"type\":\"dir\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"html\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\"}}]",
		},
		{
			"sdk", []interface{}{"dotnet", "java"}, "/repos/alex-held/dev-env-registry/contents/sdk",
			"[{\"name\":\"dotnet\",\"path\":\"sdk/dotnet\",\"sha\":\"859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"size\":0,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"download_url\":null,\"type\":\"dir\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"html\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\"}},{\"name\":\"java\",\"path\":\"sdk/java\",\"sha\":\"859e4a060e287c06f777da09fbf8fe51dc4afc92\",\"size\":0,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/java?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/java\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc92\",\"download_url\":null,\"type\":\"dir\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/java?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc92\",\"html\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/java\"}}]",
		},
	}
	for _, test := range tests {
		api, mux, teardown := Setup()
		mux.HandleFunc(test.ExpectedPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			log.Trace().Interface("responseWriter", w).Interface("request", *r).Object("expected result", test).Send()
			_, _ = fmt.Fprint(w, test.Body)
		})

		result, err := api.GetPackages(test.InputPath)
		teardown()
		require.NoError(t, err)
		require.ElementsMatch(t, result, test.ExpectedResult)
		log.Info().Msg("Test Succeeded")
	}
}

func TestGetVersions(t *testing.T) {
	test := Test{
		"sdk/dotnet", []interface{}{"3.1.202"}, "/repos/alex-held/dev-env-registry/contents/sdk/dotnet",
		"[{\"name\":\"3.1.202.yaml\",\"path\":\"sdk/dotnet/3.1.202.yaml\",\"sha\":\"ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"size\":1184,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet/3.1.202.yaml?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/blob/master/sdk/dotnet/3.1.202.yaml\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/blobs/ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"download_url\":\"https://raw.githubusercontent.com/alex-held/dev-env-registry/master/sdk/dotnet/3.1.202.yaml\",\"type\":\"file\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet/3.1.202.yaml?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/blobs/ab4c4abcc41fed88d6e7ba56c8c8db094b448c04\",\"html\":\"https://github.com/alex-held/dev-env-registry/blob/master/sdk/dotnet/3.1.202.yaml\"}}]",
	}

	api, mux, teardown := Setup()
	defer teardown()
	mux.HandleFunc(test.ExpectedPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Trace().Interface("responseWriter", w).Interface("request", *r).Object("expected result", test).Send()
		_, _ = fmt.Fprint(w, test.Body)
	})

	result, err := api.GetPackageVersions(test.InputPath)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Interface("result version", result).Msg("Packages found")
	assert.ElementsMatch(t, result, test.ExpectedResult)
}

func Setup() (api GithubAPI, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	client := github.NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL + "/")
	api = NewTestGithubAPI(client)
	return api, mux, server.Close
}
