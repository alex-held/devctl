package action

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/alex-held/devctl/internal/devctlpath"
	"github.com/alex-held/devctl/internal/sdkman"
	"github.com/alex-held/devctl/internal/system"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Install struct {
	fs     afero.Fs
	client *sdkman.Client
}

func extractArchive(fs afero.Fs, extractIntoDir string, f afero.File) (files []string, err error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	reader, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return nil, err
	}

	for _, file := range reader.File {
		readCloser, err := file.Open()
		if err != nil {
			return files, err
		}
		extractToFile := filepath.Join(extractIntoDir, file.Name)
		if !strings.HasSuffix(file.Name, "/") {
			files = append(files, extractToFile)
		}
		fileWriteTo, err := fs.Create(extractToFile)
		if err != nil {
			_ = readCloser.Close()
			return files, err
		}
		bytes := make([]byte, file.UncompressedSize64)
		_, err = readCloser.Read(bytes)
		if err != io.EOF && err != nil {
			_ = readCloser.Close()
			return files, err
		}
		_, err = fileWriteTo.Write(bytes)
		if err != nil {
			_ = readCloser.Close()
			return files, err
		}
		_ = readCloser.Close()
	}
	return files, nil
}

func saveArchive(fs afero.Fs, buf bytes.Buffer, path string) (file afero.File, err error) {
	exist, err := afero.Exists(fs, path)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot check whether downloaded archive already exists; archive=%s\n", path)
	}
	if exist {
		return nil, nil
	}
	archive, err := fs.Create(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create file; path=%s", path)
	}
	n, err := io.Copy(archive, &buf)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to copy http content into archive file; written %d bytes\n", n)
	}
	return archive, nil
}

func (i *Install) Install(ctx context.Context, sdk, version string) error {

	// download zip archive
	dl, err := i.client.Download.DownloadSDK(ctx, sdk, version, system.GetCurrent())
	if err != nil {
		return errors.Wrap(err, "error downloading sdk from api.sdkman.io")
	}

	// copy content to zip archive in devctl download path
	archiveName := fmt.Sprintf("%s-%s.zip", sdk, version)
	archiveDir := devctlpath.DownloadPath(sdk, version)
	archivePath := filepath.Join(archiveDir, archiveName)
	archive, err := saveArchive(i.fs, dl.Buffer, archivePath)
	if err != nil {
		return errors.Wrapf(err, "unable to save http content to zip file; path=%s\n", archivePath)
	}

	// extract zip archive from devctl download path into devctl sdk path
	extractIntoDir := devctlpath.SDKsPath(sdk, version)
	_, err = extractArchive(i.fs, extractIntoDir, archive)

	if err != nil {
		return errors.Wrapf(err, "unable to extract files from zip; archive=%s; dir=%s\n", archivePath, extractIntoDir)
	}

	return nil
}
