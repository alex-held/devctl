package action

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alex-held/devctl/internal/devctlpath"
	"github.com/alex-held/devctl/internal/sdkman"
	"github.com/alex-held/devctl/internal/testutils"
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func setup() (client *sdkman.Client, logger *logrus.Logger, mux *http.ServeMux, out bytes.Buffer, teardown testutils.Teardown) {
	logger = testutils.NewLogger(&out)

	mux = http.NewServeMux()
	fs := afero.NewMemMapFs()

	apiHandler := http.NewServeMux()
	const baseURLPath = "/2"
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintln(os.Stderr, "FAIL: ClientIn.BaseURL path prefix is not preserved in the request URL:")
		_, _ = fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		http.Error(w, "ClientIn.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)

	client = sdkman.NewSdkManClient(
		sdkman.URLOptions(server.URL+"/2"),
		sdkman.FileSystemOption(fs),
		sdkman.HTTPClientOption(&http.Client{}),
	)

	teardown = func() {
		server.Close()
	}
	return client, logger, mux, out, teardown
}

func TestInstall_InstallF(t *testing.T) {

	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Installer.Install()", func() {

		var client *sdkman.Client
		var logger *logrus.Logger
		var mux *http.ServeMux
		var teardown testutils.Teardown
		var err error
		var testdata []byte
		var fs afero.Fs
		var i Install
		var ctx context.Context
		var expectedArchivePath, expectedSdkDir string

		g.JustBeforeEach(func() {
			ctx = context.Background()
			expectedArchivePath = devctlpath.DownloadPath("scala", "2.13.4", "scala-2.13.4.zip")
			expectedSdkDir = devctlpath.SDKsPath("scala", "2.13.4")
			client, logger, mux, _, teardown = setup()
			testdata, err = ioutil.ReadFile("testdata/scala-2.13.4.zip")
			if err != nil {
				g.Fatalf("error reading testdata; error=%v\n", err)
			}

			// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwin
			mux.HandleFunc("/broker/download/scala/2.13.4/darwin", func(w http.ResponseWriter, r *http.Request) {

				w.Header().Add("content-type", "application/zip")
				w.Header().Add("accept-ranges", "actualDownloadContent")
				w.Header().Add("content-length", fmt.Sprintf("%d", len(testdata)))
				n, e := io.Copy(w, bytes.NewBuffer(testdata))
				if e != nil {
					g.Fatalf("error writing testdata into http.Response; error=%v\n", err)
				}
				logger.
					WithField("length", n).
					Warnln("written testdata into http.Response")
			})

			fs = afero.NewMemMapFs()
			i = Install{
				fs:     fs,
				client: client,
			}
		})

		g.AfterEach(func() {
			teardown()
		})

		g.It("Saves archive into archive folder", func() {

			err = i.Install(ctx, "scala", "2.13.4")
			Expect(err).Should(BeNil(), "install failed; error=%v\n", err)

			testExists(g, fs, expectedArchivePath, "archive does not exist")
		})

		g.It("Extracts archive into sdk folder", func() {
			err = i.Install(ctx, "scala", "2.13.4")
			Expect(err).Should(BeNil(), "install failed; error=%v\n", err)

			testExists(g, fs, expectedSdkDir, "sdk directory does not exist")
		})
	})
}

func testExists(g *goblin.G, fs afero.Fs, expected, msg string) {
	g.Helper()
	exists, err := afero.Exists(fs, expected)
	if err != nil {
		g.Fatalf("error occurred while testing whether file or dir exists; path=%s; error=%v\n", expected, err)
	}
	Expect(exists).Should(BeTrue(), "%s; path=%s", msg, expected)
}
