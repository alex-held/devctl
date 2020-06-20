package config

import (
	"fmt"
	"os"
	"path"
)

type PathFactory interface {
	GetUserHome() string
	GetDevEnvHome() string
	GetPkgRoot() string
	GetPkgDir(name string, version string) string
	GetManifests() string
}

func NewDefaultPathFactory() DefaultPathFactory {
	return DefaultPathFactory{
		UserHomeOverride: nil,
		DevEnvDirectory:  ".devenv",
	}
}

type DefaultPathFactory struct {
	UserHomeOverride *string
	DevEnvDirectory  string
}

func GetUserHome() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error resolving $HOME\n", err.Error())
		os.Exit(1)
	}
	return userHome
}

func (fac *DefaultPathFactory) GetPkgDir(name string, version string) string {
	devenv := fac.GetPkgRoot()
	pkgDir := path.Join(devenv, name, version)
	return pkgDir
}

func (fac *DefaultPathFactory) GetUserHome() string {
	if fac.UserHomeOverride != nil {
		return *fac.UserHomeOverride
	}
	return GetUserHome()
}

func (fac *DefaultPathFactory) GetDevEnvHome() string {
	return path.Join(fac.GetUserHome(), fac.DevEnvDirectory)
}

func (fac *DefaultPathFactory) GetPkgRoot() string { return path.Join(fac.GetDevEnvHome(), "sdk") }

func (fac *DefaultPathFactory) GetManifests() string {
	return path.Join(fac.GetDevEnvHome(), "manifests")
}
