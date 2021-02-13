package sdkman

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alex-held/devctl/pkg/testutils"
	"go.uber.org/zap"
	"golang.org/x/exp/errors/fmt"

	"github.com/alex-held/devctl/pkg/aarch"
	// "github.com/mitchellh/go-testing-interface"
	"github.com/spf13/afero"
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

func NewTestSdkManClient(doFunc HTTPDoFunc) (client *Client, fs afero.Fs, ctx context.Context) {
	fs = afero.NewMemMapFs()
	ctx = context.Background()

	c := &Client{
		context: ctx,
		httpClient: &httpClientStub{
			DoFunc: doFunc,
		},
		fs: fs,
	}

	c.common.client = c
	c.Download = (*DownloadService)(&c.common)
	c.ListSdks = (*ListAllSDKService)(&c.common)

	return c, fs, ctx
}

const baseURLPath = "/2"

func setup(t *testing.T) (client *Client, fs afero.Fs, mux *http.ServeMux, spy *testutils.LogSpy, teardown testutils.Teardown) {
	spy, teardown = testutils.SetupTestLogger(t)

	mux = http.NewServeMux()
	fs = afero.NewMemMapFs()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintln(os.Stderr, "FAIL: ClientIn.BaseURL path prefix is not preserved in the request URL:")
		_, _ = fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		http.Error(w, "ClientIn.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)

	client = NewSdkManClient(
		UrlOptions(server.URL+"/2"),
		FileSystemOption(fs),
		HttpClientOption(&http.Client{}),
	)

	return client, fs, mux, spy, teardown.CombineInto(func() {
		server.Close()
	})
}

func testMethod(t testing.TB, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestSdkmanClient_ListCandidates(t *testing.T) {
	client, _, mux, spy, teardown := setup(t)
	defer teardown()

	mux.HandleFunc("/candidates/all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(spy, r, "GET")
		_, _ = fmt.Fprint(w, "ant,asciidoctorj,ballerina,bpipe,btrace,ceylon,concurnas,crash,cuba,cxf,doctoolchain,dotty,gaiden,glide,gradle,gradleprofiler,grails,groovy,groovyserv,http4k,infrastructor,java,jbake,jbang,karaf,kotlin,kscript,layrry,lazybones,leiningen,maven,micronaut,mulefd,mvnd,sbt,scala,spark,springboot,sshoogr,test,tomcat,vertx,visualvm")
	})

	candidates, _, err := client.ListSdks.ListAllSDK(context.Background())
	require.NoError(spy, err)
	assert.Equal(spy, len(candidates), 43)
}

func combine(handlers ...http.HandlerFunc) http.HandlerFunc {
	if handlers == nil {
		handlers = []http.HandlerFunc{}
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		for _, handlerFunc := range handlers {
			handlerFunc(writer, request)
		}
	}
}

func handleTestdata(t *testing.T, testdataPath string) http.HandlerFunc {
	z := zap.S()

	testdataContent, err := ioutil.ReadFile(testdataPath)
	if err != nil {
		z.With("path", testdataPath).
			Fatalf("Unable to read testdata.")
	}

	return func(responseWriter http.ResponseWriter, request *http.Request) {
		length, err := responseWriter.Write(testdataContent)
		if err != nil {
			z.With(zap.Error(err)).Errorf("Unable to write testdata to http.Response.Body")
		}
		z.With(zap.Int("length", length), zap.ByteString("content", testdataContent)).Info("testdata written to http.Response.Body")
	}

}

func TestSdkmanClient_Download(t *testing.T) {
	client, fs, mux, spy, teardown := setup(t)
	ctx := context.Background()
	defer teardown()

	const expectedTestDataPath = "/Users/dev/go/src/github.com/alex-held/devctl/internal/sdkman/testdata/scala-1.8"
	expectedDownloadContent, _ := ioutil.ReadFile(expectedTestDataPath)
	z := zap.S()
	z.Infof("File contents: %s", expectedDownloadContent)

	// https://api.sdkman.io/2/broker/download/java/1.8/darwinx64
	//	mux.HandleFunc("/broker/download/scala/1.8/darwinx64", handleTestdata(t, "/Users/dev/go/src/github.com/alex-held/devctl/internal/sdkman/testdata/scala-1.8"))
	mux.HandleFunc("/2/broker/download/scala/1.8/darwinx64", handleTestdata(t, "/Users/dev/go/src/github.com/alex-held/devctl/internal/sdkman/testdata/scala-1.8"))
	mux.HandleFunc("broker/download/scala/1.8/darwinx64", handleTestdata(t, "/Users/dev/go/src/github.com/alex-held/devctl/internal/sdkman/testdata/scala-1.8"))

	testdataHandlerFunc := handleTestdata(t, "/Users/dev/go/src/github.com/alex-held/devctl/internal/sdkman/testdata/scala-1.8")
	combinedHandler := combine(
		testdataHandlerFunc,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("content-type", "application/zip")
			w.Header().Add("accept-ranges", "bytes")
			w.Header().Add("content-length", "23015564")
			w.WriteHeader(http.StatusOK)

			testMethod(spy, r, "GET")
		},
	)

	mux.HandleFunc("/broker/download/scala/1.8/darwinx64", combinedHandler)

	/*
		mux.HandleFunc("/broker/download/scala/1.8/darwinx64", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")

			w.Header().Add("content-type", "application/zip")
			w.Header().Add("accept-ranges", "bytes")
			w.Header().Add("content-length", "23015564")

			w.WriteHeader(200)
			writeLen, err := w.Write(expectedDownloadContent)
			if err != nil {
				fmt.Printf("Error writing mock response body: Error=%+v\n", err)
				//		t.Fatalf("Error writing mock response body: Error=%+v", err)
			}
			fmt.Printf("Written response body with length: %d\n", writeLen)

			//	testBody(t, r, bytes.NewReader(expectedDownloadContent))
		})
	*/

	expectedDownloadPath := os.ExpandEnv("$HOME/.devctl/archives/scala/1.8/scala-1.8")
	z.With(zap.String("path", expectedTestDataPath)).Info("Expected Download")

	download, _, err := client.Download.DownloadSDK(ctx, expectedDownloadPath, "scala", "1.8", aarch.MacOsx)

	require.NoError(t, err)

	actualDownloadPath := download.Path
	assert.Equal(t, expectedDownloadPath, actualDownloadPath)

	actualDownloadBytes, err := afero.ReadFile(fs, expectedDownloadPath)
	assert.Equal(spy, expectedDownloadContent, actualDownloadBytes)
}

func testBody(t *testing.T, r *http.Request, want io.Reader) {
	t.Helper()
	got, err := r.GetBody()
	if err != nil {
		t.Errorf("Error while accessing request body: %v", err)
	}

	gotBytes, err := ioutil.ReadAll(got)
	gotString := string(gotBytes)
	if err != nil {
		panic(err)
	}

	wantBytes, err := ioutil.ReadAll(want)
	wantString := string(wantBytes)
	if err != nil {
		panic(err)
	}

	if gotString != wantString {
		t.Errorf("Request body: %v, want %v", gotString, wantString)
	}
}
