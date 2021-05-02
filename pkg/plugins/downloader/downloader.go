package downloader

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/alex-held/devctl/internal/plugins"
)

func NewDownloader(url, desc string, dlWriter io.Writer, progressWriter io.Writer) *downloader {
	return &downloader{
		URL:            url,
		ProgressWriter: progressWriter,
		DlWriter:       dlWriter,
		DownloadDesc:   desc,
	}
}

type downloader struct {
	URL            string
	DlWriter       io.Writer
	ProgressWriter io.Writer
	DownloadDesc   string
}

func (d downloader) Download(ctx context.Context) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.URL, http.NoBody)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	size, err := getContentSize(resp, err)
	if err != nil {
		return err
	}

	var progressBar = plugins.NewProgress(d.ProgressWriter, size, d.DownloadDesc)
	var multi = io.MultiWriter(progressBar, d.DlWriter)
	_, err = io.Copy(multi, resp.Body)
	return err
}

func getContentSize(resp *http.Response, _ error) (int, error) {
	headers := resp.Header
	cl := headers.Get("Content-Length")
	size, err := strconv.Atoi(cl)
	if err != nil {
		return 0, err
	}
	return size, nil
}
