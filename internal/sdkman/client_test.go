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
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/logging"

	"github.com/alex-held/devctl/internal/system"
	"github.com/alex-held/devctl/internal/testutils"
)

const baseURLPath = "/2"

func setup() (client *Client, logger *logging.Logger, mux *http.ServeMux, teardown testutils.Teardown) {
	logger = testutils.NewLogger()

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
	return client, logger, mux, teardown
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
		var logger *logging.Logger
		var mux *http.ServeMux
		var teardown testutils.Teardown
		var ctx context.Context

		g.Describe("Download", func() {
			g.JustBeforeEach(func() {
				client, logger, mux, teardown = setup()
				ctx = context.Background()
			})

			g.AfterEach(func() {
				logger.Output.Reset()
				teardown()
			})

			g.It("Lists available sdk", func() {
				mux.HandleFunc("/candidates/all", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					_, _ = fmt.Fprint(w, "ant,asciidoctorj,ballerina,bpipe,btrace,ceylon,concurnas,crash,cuba,cxf,doctoolchain,dotty,gaiden,glide,gradle,gradleprofiler,grails,groovy,groovyserv,http4k,infrastructor,java,jbake,jbang,karaf,kotlin,kscript,layrry,lazybones,leiningen,maven,micronaut,mulefd,mvnd,sbt,scala,spark,springboot,sshoogr,test,tomcat,vertx,visualvm")
				})

				candidates, resp, err := client.ListSdks.ListAllSDK(ctx)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				logger.WithField("length", len(candidates)).Debug(candidates)
				Expect(candidates).To(HaveLen(43))
				Expect(candidates).To(ConsistOf(strings.Split("ant,asciidoctorj,ballerina,bpipe,btrace,ceylon,concurnas,crash,cuba,cxf,doctoolchain,dotty,gaiden,glide,gradle,gradleprofiler,grails,groovy,groovyserv,http4k,infrastructor,java,jbake,jbang,karaf,kotlin,kscript,layrry,lazybones,leiningen,maven,micronaut,mulefd,mvnd,sbt,scala,spark,springboot,sshoogr,test,tomcat,vertx,visualvm", ",")))
			})
		})
	})
}

func TestClient_Download(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Client", func() {
		g.Describe("Download", func() {
			expectedDownloadPath := "/tmp/downloads/scala/1.8/scala-1.8.zip"
			expectedTestDataPath := os.ExpandEnv("testdata/scala-1.8")

			var client *Client
			var logger *logging.Logger
			var mux *http.ServeMux
			var _ bytes.Buffer
			var teardown testutils.Teardown
			var ctx context.Context

			g.JustBeforeEach(func() {
				client, logger, mux, teardown = setup()
				ctx = context.Background()
			})

			g.AfterEach(func() {
				teardown()
			})

			g.It("WHEN no problems => THEN downloads SDK to local path", func() {
				expectedDownloadContent, err := ioutil.ReadFile(expectedTestDataPath)
				expectedContentBuffer := bytes.NewBuffer(expectedDownloadContent)

				if err != nil {
					errMessage := fmt.Sprintf("problem reading the testata. testdata-path: %s; error: %+v\n", expectedTestDataPath, err)
					_, _ = os.Stderr.WriteString(errMessage)
					t.Fatal(errMessage)
				}
				logger.
					WithField("path", expectedTestDataPath).
					WithField("content", expectedContentBuffer.String()).
					Warnln("loading expected-download-content from testdata")

				logger.
					WithField("path", expectedDownloadPath).
					Warnln("Expected Download Path")

				// https://api.sdkman.io/2/broker/download/scala/1.8/darwin
				mux.HandleFunc("/broker/download/scala/1.8/darwin", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("content-type", "application/zip")
					w.Header().Add("accept-ranges", "actualDownloadContent")
					w.Header().Add("content-length", fmt.Sprintf("%d", expectedContentBuffer.Len()))
					n, e := io.Copy(w, expectedContentBuffer)
					if e != nil {
						logger.
							WithError(e).
							Fatalln("error writing testdata into http.Response")
					}
					logger.
						WithField("length", n).
						Warnln("written testdata into http.Response")

					testMethod(t, r, "GET")
				})

				download, err := client.Download.DownloadSDK(ctx, "scala", "1.8", system.Darwin)

				Expect(err).To(BeNil())
				Expect(download).ShouldNot(BeNil())
				Expect(download.Buffer.Bytes()).Should(Equal(expectedDownloadContent))
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
