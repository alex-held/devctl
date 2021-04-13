package downloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var handlerPath = "/dl/go1.16.3.darwin-amd64.tar.gz"
var dlDir = "/tmp/.devctl/downloads/go/1.16.3"
var dlFilename = "go1.16.3.darwin-amd64.tar.gz"

const testArchive = "testdata/go1.16.3.darwin-amd64.tar.gz"

func setup(t *testing.T, w io.Writer, args ...interface{}) (data []byte, fs afero.Fs, file afero.File, sut *downloader, ctx context.Context) {
	t.Helper()
	var err error

	var aData []byte
	var aFs afero.Fs
	var aCtx context.Context

	if len(args) > 0 {
		for _, arg := range args {
			switch v := arg.(type) {
			case context.Context:
				aCtx = v
			case []byte:
				aData = v
			case afero.Fs:
				aFs = v
			}
		}
	}

	if aData != nil {
		data = aData
	} else {
		var testData, err = ioutil.ReadFile(testArchive) //nolint:govet
		if err != nil {
			t.Fatal(err)
		}
		data = bytes.Repeat(testData, 1000)
	}

	expectedSize := len(data)
	contentLength := fmt.Sprintf("%d", expectedSize)

	if aCtx != nil {
		ctx = aCtx
	} else {
		ctx = context.Background()
	}

	mux := http.NewServeMux()
	mux.HandleFunc(handlerPath, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Length", contentLength)
		_, e := w.Write(data)
		if e != nil {
			t.Fatal(e)
		}
	})

	s := httptest.NewServer(mux)
	if aFs != nil {
		fs = aFs
	} else {
		fs = afero.NewMemMapFs()
	}

	err = fs.MkdirAll(dlDir, 0700)
	if err != nil {
		t.Fatal(err)
	}

	file, err = fs.Create(path.Join(dlDir, dlFilename))
	if err != nil {
		t.Fatal(err)
	}
	sut = NewDownloader(s.URL+handlerPath, "testing ...", file, w)

	return data, fs, file, sut, ctx
}

func TestDownloader_Writes_Downloaded_Bytes_Into_Configured_OutWriter(t *testing.T) {
	data, fs, _, sut, ctx := setup(t, os.Stdout)

	err := sut.Download(ctx)
	assert.NoError(t, err)
	actualData, err := afero.ReadFile(fs, path.Join(dlDir, dlFilename))
	assert.NoError(t, err)
	assert.Equal(t, data, actualData)
}

func TestDownloader_Cancels_Download_When_Context_Is_Done(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Millisecond) //nolint:govet
	_, _, _, sut, _ := setup(t, os.Stdout, ctx)                             //nolint:dogsled

	err := sut.Download(ctx)
	assert.Error(t, err)
}
