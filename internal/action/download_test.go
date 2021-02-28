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
	"github.com/spf13/afero"
)

func TestDownloader_DownloadF(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Download.Download()", func() {
		var fixture *ActionTestFixture
		var err error
		var archive afero.File
		var testdata []byte
		var sut Download
		var expectedArchivePath string

		g.JustBeforeEach(func() {
			fixture = SetupFixture()

			expectedArchivePath = fixture.pather.Download("scala", "2.13.4", "scala-2.13.4.zip")
			testdata, err = ioutil.ReadFile("testdata/scala-2.13.4.zip")
			if err != nil {
				g.Fatalf("error reading testdata; error=%v\n", err)
			}

			// https://api.sdkman.io/2/broker/download/scala/2.13.4/darwin
			fixture.mux.HandleFunc("/broker/download/scala/2.13.4/darwin", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("content-type", "application/zip")
				w.Header().Add("content-length", fmt.Sprintf("%d", len(testdata)))
				_, e := io.Copy(w, bytes.NewBuffer(testdata))
				if e != nil {
					g.Fatalf("error writing testdata into http.Response; error=%v\n", err)
				}
			})

			sut = *NewActions(WithFs(fixture.fs), WithSdkmanClient(fixture.client)).Download
		})

		g.AfterEach(func() {
			fixture.teardown()
		})

		g.It("Downloaded archive has correct size", func() {
			archive, err = sut.Download(fixture.context, "scala", "2.13.4")
			Expect(err).Should(BeNil(), "download failed; error=%v\n", err)

			stat, e := archive.Stat()
			if e != nil {
				g.Fatalf("failed to get archive-file file stats; archive=%s; error=%v\n", archive.Name(), e)
			}
			actualSize := stat.Size()
			Expect(actualSize).Should(Equal(int64(len(testdata))))
		})

		g.It("Downloaded archive has correct path", func() {
			archive, err = sut.Download(fixture.context, "scala", "2.13.4")
			Expect(err).Should(BeNil(), "download failed; error=%v\n", err)
			Expect(archive.Name()).Should(Equal(expectedArchivePath))
		})

		g.It("Downloaded archive exists on file system", func() {
			_, err = sut.Download(fixture.context, "scala", "2.13.4")
			Expect(err).Should(BeNil(), "download failed; error=%v\n", err)
			testExists(g, fixture.fs, expectedArchivePath, "archive does not exist")
		})
	})
}
