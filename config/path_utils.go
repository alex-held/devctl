package config

import (
	"fmt"
	"os"
	"path"
)

func GetUserHome() string {
	userHome, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("Error resolving $HOME\n", err.Error())
		os.Exit(1)
	}
	return userHome
}

func GetDevEnvHome() string { return path.Join(GetUserHome(), ".dev-env") }
func GetSdks() string       { return path.Join(GetDevEnvHome(), "sdk") }
func GetInstallers() string { return path.Join(GetDevEnvHome(), "installers") }
func GetManifests() string  { return path.Join(GetDevEnvHome(), "manifests") }
