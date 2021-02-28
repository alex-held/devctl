package action

import (
	"archive/zip"
	"context"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/devctlpath"
)

// Install Provides the Install Action
type Install action

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
		err = extractFile(fs, extractIntoDir, file, func(s string) {
			files = append(files, s)
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func extractFile(fs afero.Fs, extractIntoDir string, file *zip.File, appender func(string)) (err error) {
	rc, err := file.Open()
	if err != nil {
		return err
	}

	defer rc.Close()
	if err != nil {
		return err
	}

	sanitized := sanitize(file.Name)
	extractToFile := filepath.Join(extractIntoDir, sanitized)
	if !strings.HasSuffix(file.Name, "/") {
		appender(extractToFile)
	}
	fileWriteTo, err := fs.Create(extractToFile)
	if err != nil {
		return err
	}
	b := make([]byte, file.UncompressedSize64)
	_, err = rc.Read(b)
	if err != io.EOF && err != nil {
		return err
	}
	_, err = fileWriteTo.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func sanitize(name string) (sanitized string) {
	name = strings.ReplaceAll(name, `\`, `/`)
	sanitized = path.Clean(name)
	sanitized = strings.TrimPrefix(sanitized, "/")

	for strings.HasPrefix(name, "../") {
		sanitized = sanitized[len("../"):]
	}
	return sanitized
}

func (i *Install) Install(ctx context.Context, sdk, version string) (dir string, err error) {
	archive, err := i.Actions.Download.Download(ctx, sdk, version)
	extractIntoDir := devctlpath.SDKsPath(sdk, version)
	if err != nil {
		return extractIntoDir, errors.Wrapf(err, "error executing download action; sdk=%s; version=%s\n", sdk, version)
	}

	_, err = extractArchive(i.Options.Fs, extractIntoDir, archive)
	if err != nil {
		return extractIntoDir, errors.Wrapf(err,
			"unable to extract files from zip; archive=%s; dir=%s\n",
			archive.Name(),
			extractIntoDir)
	}

	return extractIntoDir, nil
}
