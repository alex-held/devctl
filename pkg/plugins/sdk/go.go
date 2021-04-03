package sdk

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

type devctl_sdkplugin_go struct {
	FS     afero.Fs
	Pather devctlpath.Pather
}

func (p *devctl_sdkplugin_go) NewFunc() interface{ SDKPlugin } { return &devctl_sdkplugin_go{} }

func (p *devctl_sdkplugin_go) Name() string {
	return "devctl-sdkplugin-go"
}

func (p *devctl_sdkplugin_go) ListVersions() (versions []string) {
	sdk_go_root := p.Pather.SDK("go")
	fileInfos, err := afero.ReadDir(p.FS, sdk_go_root)

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

	fmt.Printf("found versions: %v", versions)
	return versions
}

func (p *devctl_sdkplugin_go) isValidVersion(dirname string) (version string, valid bool) {
	_, err := semver.ParseTolerant(dirname)
	if err != nil {
		return "", false
	}
	return dirname, true
}

func (p *devctl_sdkplugin_go) Download(version string) {
	fmt.Printf("downloading go sdk version %s;", version)
}

func (p *devctl_sdkplugin_go) Install(version string) {
	fmt.Printf("installing go sdk version %s;", version)
}
