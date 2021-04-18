package golang

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/alex-held/devctl/cli/internal/testutils"
	"github.com/alex-held/devctl/pkg/mocks"
	plugins2 "github.com/alex-held/devctl/pkg/plugins"
	"github.com/alex-held/devctl/pkg/system"
	"github.com/gobuffalo/plugins"
	"github.com/mandelsoft/vfs/pkg/memoryfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

const (
	rootPath = "/root"
	version  = "1.51"
)

type NamedNoOpPlugin struct {
	Name string
	plugins2.NoOpPlugin
}

type testArchiveHttpHandler struct {
	ArchiveBytes []byte // embed: testdata/archive.tar.gz
}

func (t testArchiveHttpHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(t.ArchiveBytes)))
	_, e := w.Write(t.ArchiveBytes)
	if e != nil {
		Panic()
	}
}

func (p *NamedNoOpPlugin) PluginName() string { return p.Name }

var _ = Describe("go-plugin USE", func() {
	var (
		vs              vfs.VFS
		versionSdkDir   string
		expectedCurrent string
		pp              devctlpath.Pather
		sut             *GoUseCmd
		dlCmd           *GoDownloadCmd
		linkerCmd       *GoLinkerCmd
		installerCmd    *GoInstallCmd
		srvr            *httptest.Server
		mux             *http.ServeMux
		riGetter        system.RuntimeInfoGetter
	)

	AfterSuite(func() {
		srvr.Close()
	})

	BeforeEach(func() {

		riGetter = mocks.MockRuntimeInfoGetter{
			RuntimeInfo: system.RuntimeInfo{
				OS:   "darwin",
				Arch: "amd64",
			},
		}
		mux = http.NewServeMux()
		mux.Handle("/", testArchiveHttpHandler{})

		srvr = httptest.NewServer(mux)
		vs = vfs.New(memoryfs.New())
		pp = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			return rootPath
		}))
		linkerCmd = &GoLinkerCmd{
			path: pp,
			fs:   vs,
		}
		installerCmd = &GoInstallCmd{
			runtime: riGetter,
			path:    pp,
			Fs:      vs,
		}
		dlCmd = &GoDownloadCmd{
			Fs:      vs,
			BaseUri: srvr.URL,
			Pather:  pp,
			Runtime: riGetter,
		}
		sut = &GoUseCmd{
			plugins: []plugins.Plugin{
				installerCmd,
				linkerCmd,
				dlCmd,
			},
			path: pp,
			fs:   vs,
		}
		versionSdkDir = pp.SDK("go", version)
		expectedCurrent = pp.SDK("go", "current")
	})

	Context("USE <version>", func() {
		When("using the TaskRunner", func() {
			BeforeEach(func() {
				_ = vs.MkdirAll(versionSdkDir, os.ModePerm)
				_ = vs.MkdirAll(pp.SDK("go", version, "src"), os.ModePerm)
				_ = vs.MkdirAll(pp.SDK("go", version, "doc"), os.ModePerm)
				_ = vs.MkdirAll(pp.SDK("go", version, "bin"), os.ModePerm)

				feeder := plugins.Feeder(func() []plugins.Plugin {
					return []plugins.Plugin{
						&NamedNoOpPlugin{
							Name: GoDownloadCmdName,
						},
					}
				})

				sut.WithPlugins(feeder)
			})

			It("resolves the correct plugins", func() {
				runner := sut.CreateTaskRunner(GoDownloadCmdName)
				desc := runner.Describe()
				Expect(desc).Should(ContainSubstring(GoDownloadCmdName))
			})

		})

		When("no @current version has been installed", func() {
			BeforeEach(func() {
				_ = vs.MkdirAll(versionSdkDir, os.ModePerm)
				_ = vs.MkdirAll(pp.SDK("go", version, "src"), os.ModePerm)
				_ = vs.MkdirAll(pp.SDK("go", version, "doc"), os.ModePerm)
				_ = vs.MkdirAll(pp.SDK("go", version, "bin"), os.ModePerm)
			})

			It("The new SDK Version is symlinked to @current Version ", func() {
				Expect(sut.ExecuteCommand(context.Background(), "devctl", []string{"use", version})).To(Succeed())
				linkDest, _ := vs.Readlink(expectedCurrent)
				Expect(linkDest).Should(Equal(versionSdkDir))
			})
		})

		When("@current has already a sdk configured", func() {

			BeforeEach(func() {
				_ = vs.MkdirAll(pp.SDK("go"), os.ModePerm)
				_ = vs.MkdirAll(pp.Download("go"), os.ModePerm)
				_ = vs.MkdirAll(pp.SDK("go", "19.5"), os.ModePerm)
				_ = vs.Symlink(pp.SDK("go", "19.5"), expectedCurrent)
			})

			It("replaces @current symlink with which links to <version>", func() {
				Expect(sut.ExecuteCommand(context.Background(), "devctl", []string{"use", "1.16.3"})).To(Succeed())
				Expect(expectedCurrent).Should(BeASymlink(vs))
				actual, err := vs.Readlink(expectedCurrent)
				Expect(err).Should(Succeed())
				Expect(actual).ShouldNot(Equal("1.16.3"))
			})

			It("replaces @current symlink with which links to <version>", func() {
				Expect(sut.ExecuteCommand(context.Background(), "devctl", []string{"use", "1.16.3"})).To(Succeed())
				currentFi, err := vs.Lstat(expectedCurrent)
				statCurrentFi, err := vs.Stat(expectedCurrent)
				_ = statCurrentFi

				Expect(expectedCurrent).Should(BeASymlink(vs))
				Expect(currentFi).ShouldNot(BeNil())
				Expect(err).Should(BeNil())

				expectedNewVersion := pp.SDK("go", "1.16.3")
				newVersionFi, err := vs.Lstat(expectedNewVersion)
				Expect(err).Should(BeNil())
				Expect(newVersionFi).ShouldNot(BeNil())
				Expect(expectedNewVersion).Should(BeADirectoryFs(vs))

				expectedOldVersion := pp.SDK("go", "19.5")
				oldFi, err := vs.Lstat(expectedOldVersion)
				Expect(oldFi).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				Expect(expectedOldVersion).Should(Or(BeADirectoryFs(vs), BeASymlink(vs)))
			})

			It("removes symlink from old to current", func() {
				Expect(sut.ExecuteCommand(context.Background(), "devctl", []string{"use", "1.16.3"})).To(Succeed())

				Expect(pp.SDK("go", "19.5")).Should(BeADirectoryFs(vs))
				_, err := vs.Readlink(pp.SDK("go", "19.5"))
				Expect(err).ShouldNot(Succeed())
			})
		})
	})
})

/*
func TestGoUseCmd_Link(t *testing.T) {
	const rootPath = "/home/user1/.devctl"
	const version = "1.0.0"

	goSdkDir := path.Join(rootPath, "sdks", "go")
	versionSdkDir := path.Join(goSdkDir, version)

	var memFs vfs.FileSystem
	var fs vfs.VFS
	var sut *GoUseCmd

	BeforeEach(func() {
		memFs = memoryfs.New()
		fs = vfs.New(memFs)

		sut = &GoUseCmd{
			path: devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
				return rootPath
			})),
			fs: fs,
		}
	})

	Context("Go SDK Plugin - Use", func() {

		When("no 'current' symlink in go sdk directory", func() {
			BeforeEach(func() {
				_ = fs.MkdirAll(goSdkDir, os.ModePerm)
				_ = fs.MkdirAll(path.Join(versionSdkDir, "bin"), os.ModePerm)
				_ = fs.MkdirAll(path.Join(versionSdkDir, "src"), os.ModePerm)
				_ = fs.MkdirAll(path.Join(versionSdkDir, "doc"), os.ModePerm)
			})

			It("The new GOROOT is symlinked to $SDKPATH/go/current", func() {
				var err = sut.ExecuteCommand(context.Background(), "devctl", []string{version})
				Expect(err).Should(BeNil())
				currejnt
				Expect(fs.)
				//	ExpectFolders(fs, path.Join(goSdkDir, "current"), []string{"bin", "src", "doc"}, nil)
			})
		})
	})
}
*/
