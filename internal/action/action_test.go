package action

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/devctlpath"
	"github.com/alex-held/devctl/internal/sdkman"
	"github.com/alex-held/devctl/internal/testutils"
)

func testExists(g *goblin.G, fs afero.Fs, expected, msg string) {
	g.Helper()
	exists, err := afero.Exists(fs, expected)
	if err != nil {
		g.Fatalf("error occurred while testing whether file or dir exists; path=%s; error=%v\n", expected, err)
	}
	Expect(exists).Should(BeTrue(), "%s; path=%s", msg, expected)
}

type ActionTestFixture struct {
	actions  *Actions
	fs       afero.Fs
	logger   *logrus.Logger
	mux      *http.ServeMux
	out      bytes.Buffer
	teardown testutils.Teardown
	context  context.Context
	client   *sdkman.Client
	pather   devctlpath.Pather
}

func SetupFixture() (fixture *ActionTestFixture) {
	const baseURLPath = "/2"
	var out bytes.Buffer

	logger := testutils.NewLogger(&out)
	pather := devctlpath.NewPather()
	mux := http.NewServeMux()
	fs := afero.NewMemMapFs()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintln(os.Stderr, "FAIL: ClientIn.BaseURL path prefix is not preserved in the request URL:")
		_, _ = fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		http.Error(w, "ClientIn.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)
	teardown := func() {
		server.Close()
	}
	client := sdkman.NewSdkManClient(
		sdkman.URLOptions(server.URL+baseURLPath),
		sdkman.FileSystemOption(fs),
		sdkman.HTTPClientOption(&http.Client{}),
	)

	actions := NewActions(WithFs(fs), WithSdkmanClient(client))

	fixture = &ActionTestFixture{
		actions:  actions,
		logger:   logger,
		mux:      mux,
		pather:   pather,
		out:      out,
		fs:       fs,
		teardown: teardown,
		context:  context.Background(),
		client:   client,
	}
	return fixture
}

func TestNewActions(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	fs := afero.NewMemMapFs()
	pather := devctlpath.NewPather()
	client := sdkman.NewSdkManClient()

	g.Describe("NewActions", func() {
		g.It("WithFs", func() {
			actions := NewActions(WithFs(fs))
			Expect(actions.Fs).Should(Equal(fs))
		})

		g.It("WithPather", func() {
			actions := NewActions(WithPather(pather))
			Expect(actions.Pather).Should(Equal(pather))
		})

		g.It("WithSdkmanClient", func() {
			actions := NewActions(WithSdkmanClient(client))
			Expect(actions.Client).Should(Equal(client))
		})
	})
}
