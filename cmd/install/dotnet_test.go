package install

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

type UserContextTest struct {
	PathGetter func() string
	Expected   string
	BaseTest
}
type DownloadTest struct {
	Url       string
	Directory string
	Filename  string
	BaseTest
}

type InstallTest struct {
	Commands  []Executable
	Expected  []string
	Directory string
	BaseTest
}
type BaseTest struct {
	Assertions *assert.Assertions
	FS         afero.Fs
}

type MockCommand struct {
	UnixCommand
}

func (cmd *MockCommand) execute() string {
	command := cmd.Render()
	return command
}

func TestDownload_Saves_RemoteFile_ToDirectory(t *testing.T) {
	test := newDownloadTest(t)
	test.Url = "https://raw.githubusercontent.com/alex-held/dev-env/develop/go.mod"
	test.Directory, _ = afero.TempDir(test.FS, GetDevEnvHome(), "")
	test.Filename = "go.mod"
	test.run()
}

func TestInstall_Runs_Install_Commands(t *testing.T) {
	test := newInstallTest(t, []string{
		"curl https://download.visualstudio.microsoft.com/download/pr/1016a722-2794-4381-88b8-29bf382901be/ea17a73205f9a7d33c2a4e38544935ba/dotnet-sdk-3.1.202-osx-x64.pkg >> /Users/dev/.dev-env/installers/dotnet-sdk-3.1.202-osx-x64.pkg",
		"sudo installer -verbose -pkg /Users/dev/.dev-env/installers/dotnet-sdk-3.1.202-osx-x64.pkg -target /Users/dev/.dev-env/sdk/dotnet/",
	}, &MockCommand{
		UnixCommand: UnixCommand{
			Template: "curl %s >> %s",
			Args: []string{
				"https://download.visualstudio.microsoft.com/download/pr/1016a722-2794-4381-88b8-29bf382901be/ea17a73205f9a7d33c2a4e38544935ba/dotnet-sdk-3.1.202-osx-x64.pkg",
				"/Users/dev/.dev-env/installers/dotnet-sdk-3.1.202-osx-x64.pkg",
			},
		},
	},
		&MockCommand{
			UnixCommand: UnixCommand{
				Template: "sudo installer -verbose -pkg %s -target %s",
				Args: []string{
					"/Users/dev/.dev-env/installers/dotnet-sdk-3.1.202-osx-x64.pkg",
					"/Users/dev/.dev-env/sdk/dotnet/",
				},
			},
		})

	test.run()
}

func newBaseTest(t *testing.T) BaseTest {
	return BaseTest{
		Assertions: assert.New(t),
		FS:         afero.NewMemMapFs(),
	}
}

func newDownloadTest(t *testing.T) DownloadTest {
	return DownloadTest{
		Url:       "",
		Directory: GetDevEnvHome(),
		BaseTest:  newBaseTest(t),
	}
}

func newInstallTest(t *testing.T, expected []string, commands ...Executable) InstallTest {
	return InstallTest{
		Commands:  commands,
		Expected:  expected,
		Directory: GetDevEnvHome(),
		BaseTest:  newBaseTest(t),
	}
}

func (test *InstallTest) run() {
	result := []string{}
	for i, cmd := range test.Commands {
		expected := test.Expected[i]
		actual := cmd.Execute()
		result = append(result, actual)
		test.Assertions.Equal(expected, actual)
	}
	test.Assertions.ElementsMatch(test.Expected, result)
}

func (test *DownloadTest) run() {
	expected := path.Join(test.Directory, test.Filename)
	file, err := Download(test.FS, test.Url, test.Directory)
	if err != nil {
		test.Assertions.Error(err)
	}
	test.Assertions.True(afero.Exists(test.FS, expected))
	test.Assertions.Equal(expected, *file)
}

func TestUserContext_GetSdks(t *testing.T) {
	home, _ := os.UserHomeDir()
	test := UserContextTest{
		PathGetter: GetSdks,
		Expected:   path.Join(home, ".dev-env", "sdk"),
	}
	test.run(t)
}

func TestUserContext_GetInstallers(t *testing.T) {
	home, _ := os.UserHomeDir()
	test := UserContextTest{
		PathGetter: GetInstallers,
		Expected:   path.Join(home, ".dev-env", "installers"),
	}
	test.run(t)
}

func TestUserContext_GetUserHome(t *testing.T) {
	home, _ := os.UserHomeDir()
	test := UserContextTest{
		PathGetter: GetUserHome,
		Expected:   home,
	}
	test.run(t)
}

func TestUserContext_GetDevEnvHome(t *testing.T) {

	test := UserContextTest{
		PathGetter: GetDevEnvHome,
		Expected:   path.Join(GetUserHome(), ".dev-env"),
	}
	test.run(t)
}

func (test *UserContextTest) run(t *testing.T) {
	actual := test.PathGetter()
	fmt.Printf("Expected: %s\nActual: %s\n", test.Expected, actual)
	assert.Equal(t, test.Expected, actual)
}
