package meta

import (
	"os"
	path2 "path"
	"strings"
)

const (
	DEVENV_HOME = "DEVENV_HOME"
)

type Source struct {
	URL    *string `yaml:"url"`
	Folder *string `yaml:"folder"`
	Sha256 *string `yaml:"sha256"`
}

type Meta struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	InstallRoot string   `yaml:"install_root"`
	LinkRoot    string   `yaml:"link_root"`
	Sources     []Source `yaml:"sources"`
	Install     []string `yaml:"install"`
	Link        []string `yaml:"link"`
	Homepage    string   `yaml:"homepage"`
	Summary     string   `yaml:"summary"`
}

func normalize(content string) (result string) {
	result = strings.Trim(content, "\n")
	result = strings.Trim(content, "")
	result = strings.Trim(content, "\n")
	return result
}

func NewRemoteArchiveSource(sha string, url string) Source {
	return Source{
		URL:    &url,
		Sha256: &sha,
	}
}

func (u Meta) GetUserHome(paths ...string) (home string) {
	home, _ = os.UserHomeDir()
	if paths == nil || len(paths) == 0 {
		return home
	}
	pathSegments := append([]string{home}, paths...)
	devenvSubDirectory := path2.Join(pathSegments...)
	return devenvSubDirectory
}

func (u Meta) GetPkgPath(name string, version string) string {
	pathSegments := []string{"pkg", name, version}
	return u.Home(pathSegments)
}

// TryGetDevEnvHome Tries to get the DEVENV_HOME environment variable
// If DEVENV_HOME is unset, the method returns nil!
func (u Meta) TryGetDevEnvHome() *string {
	home := os.Getenv(DEVENV_HOME)
	if &home != nil {
		return &home
	}
	return nil
}

func (u Meta) Home(paths []string) (home string) {
	home = *u.TryGetDevEnvHome()
	if &home == nil {
		home = u.GetUserHome(paths...)
		return home
	}
	segments := append([]string{home}, paths...)
	subDirectory := path2.Join(segments...)
	return subDirectory
}
