package sdk

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/pkg/constants"
	"github.com/alex-held/devctl/pkg/devctlpath"
	"github.com/alex-held/devctl/pkg/plugins"
)

type devctlSdkpluginGo struct {
	FS                afero.Fs
	Pather            devctlpath.Pather
	HTTPClient        http.Client
	BaseURI           string
	Context           context.Context
	RuntimeInfoGetter plugins.RuntimeInfoGetter
}

func (p *devctlSdkpluginGo) NewFunc() plugins.SDKPlugin {
	return &devctlSdkpluginGo{
		FS:                afero.NewOsFs(),
		Pather:            devctlpath.DefaultPather(),
		HTTPClient:        http.Client{},
		BaseURI:           "https://golang.org",
		Context:           context.Background(),
		RuntimeInfoGetter: plugins.OSRuntimeInfoGetter{},
	}
}

func (*devctlSdkpluginGo) Name() string {
	return "devctl-sdkplugin-go"
}

func (p *devctlSdkpluginGo) ListVersions() (versions []string) {
	sdkGoRoot := p.Pather.SDK("go")
	fileInfos, err := afero.ReadDir(p.FS, sdkGoRoot)

	if err != nil {
		return versions
	}
	for _, fileInfo := range fileInfos {
		dirname := fileInfo.Name()
		if dirname == "current" || fileInfo.Mode().Type() == os.ModeSymlink {
			continue
		}

		if version, valid := p.isValidVersion(fileInfo.Name()); valid {
			versions = append(versions, version)
		}
	}

	fmt.Printf("found versions: %v\n", versions)
	return versions
}

func (p *devctlSdkpluginGo) isValidVersion(dirname string) (version string, valid bool) {
	_, err := semver.ParseTolerant(dirname)
	if err != nil {
		return "", false
	}
	return dirname, true
}

func (p *devctlSdkpluginGo) Download(version string) (err error) {
	dlPath := p.Pather.Download("go", version)
	err = p.FS.MkdirAll(dlPath, fileutil.PrivateDirMode)
	if err != nil {
		return errors.Wrapf(err, "failed creating go sdk download path; version=%s; err=%v", version, err)
	}
	filename := fmt.Sprintf("go%s.darwin-amd64.tar.gz", version)
	dlURI := fmt.Sprintf("%s/dl/%s", p.BaseURI, filename)
	req, err := http.NewRequestWithContext(p.Context, http.MethodGet, dlURI, http.NoBody)
	if err != nil {
		return errors.Wrapf(err, "failed downloading go sdk archive; version=%s; err=%v", version, err)
	}
	response, err := p.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed downloading go sdk archive; version=%s; err=%v", version, err)
	}
	defer response.Body.Close()
	filePath := path.Join(dlPath, filename)
	err = afero.WriteReader(p.FS, filePath, response.Body)
	if err != nil {
		return errors.Wrapf(err, "failed writing go sdk archive; version=%s; err=%v", version, err)
	}
	return nil
}

func (p *devctlSdkpluginGo) Link(version string) (err error) {
	sdkPath := p.Pather.SDK("go", version)
	current := p.Pather.SDK("go", "current")
	if ok, _ := afero.DirExists(p.FS, current); ok {
		_ = p.FS.Remove(current)
	}

	symlink := exec.Command("ln", "-s", sdkPath, current)
	err = symlink.Run()
	if err != nil {
		return errors.Wrapf(err, "failed linking go sdk %s; src=%s; dest=%s\n", version, sdkPath, current)
	}
	return nil
}

func (p *devctlSdkpluginGo) archiveName(version string) string {
	ri := p.RuntimeInfoGetter.Get()
	return fmt.Sprintf("go%s.%s-%s.tar.gz", version, ri.OS, ri.Arch)
}

func (p *devctlSdkpluginGo) Extract(version string) (err error) {
	archiveName := p.archiveName(version)
	archivePath := p.Pather.Download("go", version, archiveName)
	sdkPath := p.Pather.SDK("go", version)
	archive, err := p.FS.Open(archivePath)
	if err != nil {
		return errors.Wrapf(err, "failed to open go sdk archive=%s\n", archivePath)
	}
	err = p.FS.MkdirAll(p.Pather.SDK("go", version), fileutil.PrivateDirMode)
	if err != nil {
		return errors.Wrapf(err, "failed to Extract go sdk %s; dest=%s; archive=%s\n", version, sdkPath, archivePath)
	}
	err = UnTarGzip(archive, sdkPath, GoSDKUnarchiveRenamer())
	if err != nil {
		return errors.Wrapf(err, "failed to Extract go sdk %s; dest=%s; archive=%s\n", version, sdkPath, archivePath)
	}
	return err
}

func GoSDKUnarchiveRenamer() Renamer {
	return func(p string) string {
		parts := strings.Split(p, string(filepath.Separator))
		parts = parts[1:]
		newPath := strings.Join(parts, string(filepath.Separator))
		return newPath
	}
}

func (p *devctlSdkpluginGo) InstallE(version string) (err error) {
	if exists, _ := afero.DirExists(p.FS, p.Pather.SDK("go", version)); exists {
		err = p.Link(version)
		if err != nil {
			return errors.Wrapf(err, "failed to install go sdk %s", version)
		}
	}
	if exists, _ := afero.Exists(p.FS, p.Pather.Download("go", version, p.archiveName(version))); !exists {
		err = p.Download(version)
		if err != nil {
			return errors.Wrapf(err, "failed to install go sdk %s", version)
		}
	}
	err = p.Extract(version)
	if err != nil {
		return errors.Wrapf(err, "failed to install go sdk %s", version)
	}
	err = p.Link(version)
	if err != nil {
		return errors.Wrapf(err, "failed to install go sdk %s", version)
	}
	return nil
}

func (p *devctlSdkpluginGo) Install(version string) {
	err := p.InstallE(version)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(constants.Failure)
	}
}

type Renamer func(p string) string

//nolint:gocognit
func UnTarGzip(file io.Reader, target string, renamer Renamer) error {
	gr, _ := gzip.NewReader(file)
	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		filename := header.Name
		if renamer != nil {
			filename = renamer(filename)
		}

		p := filepath.Join(target, filename)
		fi := header.FileInfo()

		if fi.IsDir() {
			if e := os.MkdirAll(p, fi.Mode()); e != nil {
				return e
			}
			continue
		}
		file, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tr)
		if err != nil {
			return err
		}
	}
	return nil
}
