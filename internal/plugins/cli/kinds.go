package cli

/*
import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Kinder interface {
	Kind() string
}

type downloader struct {
	Fs     afero.Fs
	Writer ProgressWriter
	Client *http.Client
	URL    string
	Dest   string
}

type ProgressWriter struct {
	mtx *sync.Mutex
	ProgressC chan<- int
	Writer    io.Writer
	buffer    *bytes.Buffer
	BufWriter bufio.Writer

}

func (w ProgressWriter) ReadFrom(r io.Reader) (n int64, err error) {

}

func NewProgressWriter(progress chan<- int) *ProgressWriter {

	var b *bytes.Buffer
	// TODO: remove var if possible
	return &ProgressWriter{
		ProgressC: progress,
		buffer:    b,
	}
}

func (w ProgressWriter) Read(p []byte) (n int, err error)  {
	aa, err := w.buffer.Read(p)


}
func (w ProgressWriter) Write(p []byte) (written int, err error) {

	total := len(p)
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				_ = fmt.Errorf("RECOVERED %v", r)
			}
		}
	}()


	if w.buffer == nil {
		w.buffer = bytes.NewBuffer(make([]byte, total))
	}


	for i, b := range p {
		w.buffer.WriteByte(b)
		_, err := w.Writer.Write([]byte{b})
		if err != nil {
			return written, err
		}
		written++
		percent := int(float64(i) / float64(total) * 100)
		w.ProgressC <- percent
	}
cop
	return written, nil
}

func (d *downloader) Download(w io.Writer, ctx context.Context, progress chan<- int) (err error) {
	done := make(chan int64)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.URL, http.NoBody)
	if err != nil {
		return err
	}
	resp, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-done:
				break
			default:

				f, err := os.Open(d.Dest)
				if err != nil {
					return
				}
				s, err := f.Stat()
				if err != nil {
					return
				}
				percentage := int(float64(s.Size()) / float64(size) * 100.0)
				progress <- percentage
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return fmt.Errorf("(%d) unable to download package %s", resp.StatusCode, filepath.Base(d.Dest))
	}

	written, err := io.Copy(d.Writer, resp.Body)
	done <- written
	n, err := io.Copy(w, resp.Body)
	done <- n

	return err
}

type Downloader interface {
	Download(w io.Writer, progress chan<- int) (err error)
}
*/
