package install

import (
	"fmt"
	"github.com/spf13/afero"
	"io"
	"net/http"
	"os"
	"path"
)

type dotnet struct {
	Versions map[string]Installable
}

func NewDotnet() dotnet {
	sdk := dotnet{}
	sdk.Versions["3.2.102"] = &Installation{
		Version: "3.2.102",
		Commands: []string{
			"curl https://download.visualstudio.microsoft.com/download/pr/1016a722-2794-4381-88b8-29bf382901be/ea17a73205f9a7d33c2a4e38544935ba/dotnet-sdk-3.1.202-osx-x64.pkg >> /Users/dev/.dev-env/installers/dotnet-sdk-3.1.202-osx-x64.pkg",
			"sudo installer -verbose -pkg /Users/dev/.dev-env/installers/dotnet-sdk-3.1.202-osx-x64.pkg -target /Users/dev/.dev-env/sdk/dotnet/",
		},
	}
	return sdk
}

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

type Installable interface {
	Install(directory string) (string, error)
}

type Installation struct {
	Version  string
	Commands []string
}

func (installation *Installation) Install(directory string) (string, error) {
	for _, command := range installation.Commands {
		println(command)
	}

	return "", nil
}

type Artifact struct {
	Url string
}

type UnixCommand struct {
	Template string
	Args     []string
}

func Download(fs afero.Fs, url string, directory string) (*string, error) {

	remoteFilename := path.Base(url)
	downloadFilename := path.Join(directory, remoteFilename)

	response, _ := http.Get(url)
	defer response.Body.Close()

	// Check server response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", response.Status)
	}

	// Create file
	out, _ := fs.Create(downloadFilename)

	// Copy response into file
	_, err := io.Copy(out, response.Body)
	if err != nil {
		return nil, err
	}

	return &downloadFilename, nil
}

type Renderable interface {
	Render() string
}

func (cmd *UnixCommand) Render() string {
	args := []interface{}{}
	for _, arg := range cmd.Args {
		args = append(args, arg)
	}
	return fmt.Sprintf(cmd.Template, args...)
}

func (cmd *UnixCommand) Execute() string {
	command := cmd.Render()
	return command
}

type Executable interface {
	Execute() string
}

func (dn *dotnet) Uninstall(version string) {
	/* commands := []UnixCommand

	   version = "1.0.1"
	   sudo
	   rm - rf/usr/local/share/dotnet/sdk/$version
	   sudo
	   rm - rf/usr/local/share/dotnet/shared/Microsoft.NETCore.App/$version
	   sudo
	   rm - rf/usr/local/share/dotnet/shared/Microsoft.AspNetCore.All/$version
	   sudo
	   rm - rf/usr/local/share/dotnet/shared/Microsoft.AspNetCore.App/$version
	   sudo
	   rm - rf/usr/local/share/dotnet/host/fxr/$version*/
}
