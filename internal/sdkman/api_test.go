package sdkman

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type httpClientStub struct {
	StubResponse   *http.Response
	StubError      error
	InvokedRequest *http.Request
	DoFunc         HTTPDoFunc
}

func (h *httpClientStub) Do(req *http.Request) (response *http.Response, err error) {
	if h.DoFunc != nil {
		return h.DoFunc(req)
	}
	h.InvokedRequest = req
	if h.StubError != nil {
		return nil, h.StubError
	}
	return h.StubResponse, nil
}

func NewTestSdkManClient(doFunc HTTPDoFunc) Client {
	return &sdkmanClient{
		baseURL: "https://api.sdkman.io",
		version: "2",
		httpClient: &httpClientStub{
			DoFunc: doFunc,
		},
	}
}

func TestSdkmanClient_ListCandidates(t *testing.T) {
	client := NewTestSdkManClient(func(req *http.Request) (*http.Response, error) {
		responseWriter := &bytes.Buffer{}

		_, _ = responseWriter.WriteString("ant,asciidoctorj,ballerina,bpipe,btrace,ceylon,concurnas,crash,cuba,cxf,doctoolchain,dotty,gaiden,glide,gradle,gradleprofiler,grails,groovy,groovyserv,http4k,infrastructor,java,jbake,jbang,karaf,kotlin,kscript,layrry,lazybones,leiningen,maven,micronaut,mulefd,mvnd,sbt,scala,spark,springboot,sshoogr,test,tomcat,vertx,visualvm") // nolint

		return &http.Response{
			StatusCode:    200,
			Body:          ioutil.NopCloser(responseWriter),
			ContentLength: int64(responseWriter.Len()),
		}, nil
	})

	candidates, err := client.ListCandidates()
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(candidates), 1)
}
