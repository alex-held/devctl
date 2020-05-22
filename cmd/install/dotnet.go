package install

import (
	"fmt"
	"github.com/alex-held/dev-env/manifest"
	"github.com/spf13/afero"
	"io"
	"net/http"
	"path"
)

type dotnet struct {
	Versions map[string]Installable
}

type Installable interface {
	Install(m *manifest.Manifest) error
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
