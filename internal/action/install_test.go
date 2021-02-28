package action

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"

	"github.com/alex-held/devctl/internal/system"

	"github.com/alex-held/devctl/internal/devctlpath"

	"github.com/spf13/afero"
)

func TestInstall_InstallF(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Installer.Install()", func() {
		var fixture *ActionTestFixture
		var err error
		var testdata []byte
		var fs afero.Fs
		var sut Install
		var expectedArchivePath, expectedSdkDir string

		g.JustBeforeEach(func() {
			fixture = SetupFixture(g)

			expectedArchivePath = fixture.pather.Download("scala", "2.13.4", "scala-2.13.4.zip")
			expectedSdkDir = fixture.pather.SDK("scala", "2.13.4")
			testdata, err = ioutil.ReadFile("testdata/scala-2.13.4.zip")
			if err != nil {
				g.Fatalf("error reading testdata; error=%v\n", err)
			}

			// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwin
			fixture.mux.HandleFunc(fmt.Sprintf("/broker/download/scala/2.13.4/%s", system.GetCurrent()), func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("content-type", "application/zip")
				w.Header().Add("content-length", fmt.Sprintf("%d", len(testdata)))
				_, e := io.Copy(w, bytes.NewBuffer(testdata))
				if e != nil {
					g.Fatalf("error writing testdata into http.Response; error=%v\n", err)
				}
			})

			fs = afero.NewMemMapFs()
			sut = *NewActions(WithFs(fs), WithSdkmanClient(fixture.client), WithPather(devctlpath.NewPather())).Install
		})

		g.AfterEach(func() {
			fixture.teardown()
		})

		g.It("Saves archive into archive folder", func() {
			_, err = sut.Install(fixture.context, "scala", "2.13.4")
			Expect(err).Should(BeNil(), "install failed; error=%v\n", err)
			testExists(g, fs, expectedArchivePath, "archive does not exist")
		})

		g.It("Extracts archive into sdk folder", func() {
			_, err = sut.Install(fixture.context, "scala", "2.13.4")
			Expect(err).Should(BeNil(), "install failed; error=%v\n", err)

			testExists(g, fs, expectedSdkDir, "sdk directory does not exist")
		})

		g.It("Returns sdk installation folder", func() {
			dir, err := sut.Install(fixture.context, "scala", "2.13.4")
			Expect(err).Should(BeNil(), "install failed; error=%v\n", err)
			Expect(dir).Should(Equal(expectedSdkDir), "actual sdk directory does not match expected; actual=%s; expected=%s", dir, expectedSdkDir)
		})
	})
}
