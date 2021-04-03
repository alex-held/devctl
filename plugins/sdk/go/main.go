package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blang/semver"
	"github.com/spf13/afero"
)

func Name() string {
	return "devctl-sdkplugin-go"
}

func ListVersions() (versions []string) {
	devctl_root := os.Getenv("DEVCTL_ROOT")
	sdk_dir := filepath.Join(devctl_root, "sdks")
	sdk_go_root := filepath.Join(sdk_dir, "go")

	fs := afero.NewOsFs()
	fileInfos, err := afero.ReadDir(fs, sdk_go_root)

	if err != nil {
		return versions
	}

	for _, fileInfo := range fileInfos {
		println(fileInfo.Name())
		if version, valid := isValidVersion(fileInfo.Name()); valid {
			versions = append(versions, version)
		}
	}

	return versions
}

func isValidVersion(dirname string) (version string, valid bool) {
	semver, err := semver.Parse(dirname)
	if err != nil {
		return "", false
	}
	return semver.String(), true
}

func Download(version string) {
	fmt.Printf("downloading go sdk version %s;", version)
}

func Install(version string) {
	fmt.Printf("installing go sdk version %s;", version)
}
