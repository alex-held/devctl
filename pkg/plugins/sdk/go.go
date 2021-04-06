package sdk

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/pkg/devctlpath"
	"github.com/alex-held/devctl/pkg/plugins"
)

type devctlSdkpluginGo struct {
	FS     afero.Fs
	Pather devctlpath.Pather
}

func (p *devctlSdkpluginGo) NewFunc() plugins.SDKPlugin { return &devctlSdkpluginGo{} }

func (p *devctlSdkpluginGo) Name() string {
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

	fmt.Printf("found versions: %v", versions)
	return versions
}

func (p *devctlSdkpluginGo) isValidVersion(dirname string) (version string, valid bool) {
	_, err := semver.ParseTolerant(dirname)
	if err != nil {
		return "", false
	}
	return dirname, true
}

func (p *devctlSdkpluginGo) Download(version string) {
	fmt.Printf("downloading go sdk version %s;", version)
}

func (p *devctlSdkpluginGo) Install(version string) {
	fmt.Printf("installing go sdk version %s;", version)
}
