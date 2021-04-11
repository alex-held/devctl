package golang

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/mandelsoft/vfs/pkg/memoryfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/alex-held/devctl/cli/internal/testutils"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

const (
	rootPath = "root"
	version  = "1.51"
)

var _ = Describe("go-plugin USE", func() {
	var (
		vs              vfs.VFS
		versionSdkDir   string
		expectedCurrent string
		pp              devctlpath.Pather
		sut             *GoUseCmd
	)

	BeforeEach(func() {
		vs = vfs.New(memoryfs.New())
		pp = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			return rootPath
		}))

		sut = &GoUseCmd{
			path: pp,
			fs:   vs,
		}
		versionSdkDir = pp.SDK("go", version)
		expectedCurrent = pp.SDK("go", "current")
	})

	Context("USE <version>", func() {
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
				_ = vs.MkdirAll(expectedCurrent, os.ModePerm)
			})

			It("replaces @current symlink with which links to <version>", func() {

				By("Symlinking /root/sdks/go/current -> /root/sdks/go/1.16.3 \n" +
					"ln -s -v -F  /root/sdks/go/1.16.3  /root/sdks/go/current")

				Expect(vs.Symlink(expectedCurrent, pp.SDK("go", "19.5"))).To(Succeed())

				Expect(sut.ExecuteCommand(context.Background(), "devctl", []string{"use", "1.16.3"})).To(Succeed())

				currentFi, err := vs.Lstat(expectedCurrent)
				Expect(expectedCurrent).Should(And(BeASymlink(vs), Not(BeADirectoryFs(vs))))
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
		})
	})

})

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
				ExpectFolders(fs, path.Join(goSdkDir, "current"), []string{"bin", "src", "doc"}, nil)
			})
		})
	})
}

func ExpectFolders(fs vfs.FileSystem, path string, names []string, err error) {
	ExpectFoldersOrdering(fs, path, names, err, false)
}

func ExpectFoldersOrdering(fs vfs.FileSystem, path string, names []string, err error, forceOrder bool) {
	var f, openErr = fs.Open(path)
	if err == nil {
		Expect(openErr).To(BeNil())
	} else {
		Expect(openErr).Should(Equal(err))
	}

	var nms, dirErr = f.Readdirnames(0)
	Expect(dirErr).Should(BeNil())

	if names == nil {
		names = []string{}
	}
	found := append(names, nms...)
	if forceOrder {
		Expect(found).To(Equal(names))
	} else {
		Expect(found).To(ContainElements(names))
	}
}
