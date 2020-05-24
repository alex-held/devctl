package config

import (
	"fmt"
	"os"
	"path"
)

type PathFactory interface {
	GetUserHome() string
	GetDevEnvHome() string
	GetSdks() string
	GetManifests() string
}

type DefaultPathFactory struct {
	UserHomeOverride *string
	DevEnvOverride   *string
}

func GetUserHome() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error resolving $HOME\n", err.Error())
		os.Exit(1)
	}
	return userHome
}

func (fac *DefaultPathFactory) GetUserHome() string {
	if fac.UserHomeOverride != nil {
		return *fac.UserHomeOverride
	}
	return GetUserHome()
}

func (fac *DefaultPathFactory) GetDevEnvHome() string {
	var prefix = ".dev-env"
	if fac.DevEnvOverride != nil {
		prefix = *fac.DevEnvOverride
	}
	return path.Join(GetUserHome(), prefix)
}

func (fac *DefaultPathFactory) GetSdks() string { return path.Join(fac.GetDevEnvHome(), "sdk") }
func (fac *DefaultPathFactory) GetManifests() string {
	return path.Join(fac.GetDevEnvHome(), "manifests")
}
