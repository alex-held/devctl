package sdkman

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/pkg/aarch"
	"github.com/alex-held/devctl/pkg/testutils"
)

const baseURLPath = "/2"

func setup() (client *Client, logger *logrus.Logger, mux *http.ServeMux, out bytes.Buffer, teardown testutils.Teardown) { //nolint:lll
	logger = testutils.NewLogger(&out)

	mux = http.NewServeMux()
	fs := afero.NewMemMapFs()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintln(os.Stderr, "FAIL: ClientIn.BaseURL path prefix is not preserved in the request URL:")
		_, _ = fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		http.Error(w, "ClientIn.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)

	client = NewSdkManClient(
		URLOptions(server.URL+"/2"),
		FileSystemOption(fs),
		HTTPClientOption(&http.Client{}),
	)

	teardown = func() {
		server.Close()
	}
	return client, logger, mux, out, teardown
}

func testMethod(t testing.TB, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestSdkmanClient_ListCandidates(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Client", func() {
		var client *Client
		var logger *logrus.Logger
		var mux *http.ServeMux
		var out bytes.Buffer
		var teardown testutils.Teardown
		var ctx context.Context

		g.Describe("Download", func() {
			g.JustBeforeEach(func() {
				client, logger, mux, out, teardown = setup()
				ctx = context.Background()
			})

			g.AfterEach(func() {
				out.Reset()
				teardown()
			})

			g.It("Lists available sdk", func() {
				mux.HandleFunc("/candidates/all", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					_, _ = fmt.Fprint(w, "ant,asciidoctorj,ballerina,bpipe,btrace,ceylon,concurnas,crash,cuba,cxf,doctoolchain,dotty,gaiden,glide,gradle,gradleprofiler,grails,groovy,groovyserv,http4k,infrastructor,java,jbake,jbang,karaf,kotlin,kscript,layrry,lazybones,leiningen,maven,micronaut,mulefd,mvnd,sbt,scala,spark,springboot,sshoogr,test,tomcat,vertx,visualvm") //nolint:lll
				})

				candidates, resp, err := client.ListSdks.ListAllSDK(ctx)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				logger.WithField("length", len(candidates)).Debug(candidates)
				Expect(candidates).To(HaveLen(43))
				Expect(candidates).To(ConsistOf(strings.Split("ant,asciidoctorj,ballerina,bpipe,btrace,ceylon,concurnas,crash,cuba,cxf,doctoolchain,dotty,gaiden,glide,gradle,gradleprofiler,grails,groovy,groovyserv,http4k,infrastructor,java,jbake,jbang,karaf,kotlin,kscript,layrry,lazybones,leiningen,maven,micronaut,mulefd,mvnd,sbt,scala,spark,springboot,sshoogr,test,tomcat,vertx,visualvm", ","))) //nolint:lll
			})
		})
	})
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

func handleTestdata(logger *logrus.Logger, testdataPath string) http.HandlerFunc {
	testdataContent, err := ioutil.ReadFile(testdataPath)
	if err != nil {
		logger.
			WithField("path", testdataPath).
			Fatalln("Unable to read testdata.")
	}

	return func(responseWriter http.ResponseWriter, request *http.Request) {
		length, err := responseWriter.Write(testdataContent)
		if err != nil {
			logger.
				WithError(err).
				Errorln("Unable to write testdata to http.Response.Body")
		}
		logger.
			WithField("length", length).
			WithField("content", testdataContent).
			Infoln("testdata written to http.Response.Body")
	}
}

func TestClient_Download(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Client", func() {
		expectedTestDataPath := os.ExpandEnv("$HOME/go/src/github.com/alex-held/devctl/internal/sdkman/testdata/scala-1.8")
		expectedDownloadPath := os.ExpandEnv("$HOME/.devctl/archives/scala/1.8/scala-1.8")

		var client *Client
		var logger *logrus.Logger
		var mux *http.ServeMux
		var out bytes.Buffer
		var teardown testutils.Teardown
		var ctx context.Context

		RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

		g.Describe("Download", func() {
			g.JustBeforeEach(func() {
				client, logger, mux, out, teardown = setup()
				ctx = context.Background()
			})

			g.AfterEach(func() {
				out.Reset()
				teardown()
			})

			g.It("WHEN no problems => THEN downloads SDK to local path", func() {
				expectedDownloadContent, _ := ioutil.ReadFile(expectedTestDataPath)
				logger.
					WithField("path", expectedDownloadContent).
					WithField("content", expectedDownloadContent).
					Debugln("loading expected-download-content from ./testdata")

				logger.
					WithField("path", expectedTestDataPath).
					Debugln("Expected Download Path")

				combinedHandler := combine(
					handleTestdata(logger, expectedTestDataPath),
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Add("content-type", "application/zip")
						w.Header().Add("accept-ranges", "actualDownloadContent")
						w.Header().Add("content-length", "23015564")
						testMethod(t, r, "GET")
					},
				)

				// https://api.sdkman.io/2/broker/download/scala/1.8/darwinx64
				mux.HandleFunc("/broker/download/scala/1.8/darwinx64", combinedHandler)

				download, resp, err := client.Download.DownloadSDK(ctx, expectedDownloadPath, "scala", "1.8", aarch.MacOsx)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				logger.WithField("path", download.Path).Infoln("Actual Download Path")

				actualDownloadContent, err := ioutil.ReadAll(download.Reader)
				Expect(err).To(BeNil())
				Expect(actualDownloadContent).To(Equal(expectedDownloadContent))
				Expect(download.Path).To(Equal(expectedDownloadPath))
			})
		})
	})
}

func _(t *testing.T, r *http.Request, want io.Reader) {
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
