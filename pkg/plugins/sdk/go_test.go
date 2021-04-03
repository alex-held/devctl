package sdk

import (
	"path"
	"testing"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

type testPather struct {
	DevEnvConfigPath string
	SDKRoot          string
}

func (p *testPather) ConfigFilePath() string           { return p.DevEnvConfigPath }
func (p *testPather) ConfigRoot(elem ...string) string { return "" }
func (p *testPather) Config(elem ...string) string     { return "" }
func (p *testPather) Bin(elem ...string) string        { return "" }
func (p *testPather) Download(elem ...string) string   { return "" }
func (p *testPather) SDK(elem ...string) string        { return path.Join(p.SDKRoot, path.Join(elem...)) }
func (p *testPather) Cache(elem ...string) string      { return "" }

func TestGoSDKPlugin(t *testing.T) {
	g := goblin.Goblin(t)

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("devctl-sdkplugin-go", func() {
		var sut *devctl_sdkplugin_go
		var fs afero.Fs
		var pathr devctlpath.Pather

		g.BeforeEach(func() {
			fs = afero.NewMemMapFs()
			pathr = &testPather{
				DevEnvConfigPath: "/some/path/to/config.yaml",
				SDKRoot:          "/some/path/to/sdks",
			}
			sut = &devctl_sdkplugin_go{
				FS:     fs,
				Pather: pathr,
			}
		})

		g.It("Lists the installed sdks", func() {
			_ = fs.MkdirAll(pathr.SDK("go"), fileutil.PrivateDirMode)
			_ = fs.MkdirAll(pathr.SDK("go", "1.16"), fileutil.PrivateDirMode)
			_ = fs.MkdirAll(pathr.SDK("go", "1.16.2"), fileutil.PrivateDirMode)
			_ = fs.MkdirAll(pathr.SDK("go", "1.15"), fileutil.PrivateDirMode)
			_ = fs.MkdirAll(pathr.SDK("go", "1.14"), fileutil.PrivateDirMode)
			_ = fs.MkdirAll(pathr.SDK("go", "current"), fileutil.PrivateDirMode)
			defer func() { fs = afero.NewMemMapFs() }()

			expected := []string{"1.16", "1.16.2", "1.15", "1.14"}
			actual := sut.ListVersions()
			Expect(actual).Should(ContainElements(expected))
			Expect(actual).Should(HaveLen(len(expected)))
		})

		g.It("NewFunc creates a valid instance of the plugin", func() {
			actual := sut.NewFunc()
			Expect(actual).Should(Equal("devctl-sdkplugin-go"))
		})

		g.It("NewFunc creates a valid instance of the plugin", func() {
			actual := sut.NewFunc()
			Expect(actual).Should(Equal("devctl-sdkplugin-go"))
		})

		g.It("WHEN Download(<version>) is called => THEN the correct version gets getting downloaded", func() {
			sut.Download("1.16")
			downloadPath := pathr.Download("go", "1.16")
			dlDirExists, _ := afero.DirExists(fs, downloadPath)
			artifactName := path.Join(downloadPath, "golang-1.16.tar.gz")
			bytes, _ := afero.Exists(fs, artifactName)
			Expect(dlDirExists).Should(BeTrue())
			Expect(bytes).Should(BeNumerically(">=", 1))
		})


		g.It("WHEN Install(<version>) is called => THEN the correct version gets linked to current", func() {
			sut.Download("1.16")
			downloadPath := pathr.Download("go", "1.16")
			dlDirExists, _ := afero.DirExists(fs, downloadPath)
			artifactName := path.Join(downloadPath, "golang-1.16.tar.gz")
			bytes, _ := afero.Exists(fs, artifactName)
			Expect(dlDirExists).Should(BeTrue())
			Expect(bytes).Should(BeNumerically(">=", 1))
		})
	})
}
