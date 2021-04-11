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
	"github.com/onsi/gomega/types"

	"github.com/alex-held/devctl/cli/internal/testutils"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

const (
	rootPath = "/tmp/.devctl"
	version  = "1.51"
)

var _ = Describe("Go SDK Plugin", func() {

	var (
		wd, _         = os.Getwd()
		rootArgument  = path.Join(wd, "cmd/devctl")
		versionSdkDir string
		currentPath   string
		fs            vfs.VFS
		pathr         devctlpath.Pather
		sut           *GoUseCmd
	)

	BeforeEach(func() {
		fs = vfs.New(memoryfs.New())
		pathr = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			return rootPath
		}))
		sut = &GoUseCmd{
			path: pathr,
			fs:   fs,
		}
		versionSdkDir = pathr.SDK("go", version)
		currentPath = pathr.SDK("go", "current")
	})

	Context("USE <version>", func() {

		When("no @current version has been installed", func() {

			BeforeEach(func() {
				_ = fs.MkdirAll(versionSdkDir, os.ModePerm)
				_ = fs.MkdirAll(pathr.SDK("go", version, "src"), os.ModePerm)
				_ = fs.MkdirAll(pathr.SDK("go", version, "doc"), os.ModePerm)
				_ = fs.MkdirAll(pathr.SDK("go", version, "bin"), os.ModePerm)
			})

			It("The new SDK Version is symlinked to @current Version ", func() {
				err := sut.ExecuteCommand(context.Background(), rootArgument, []string{"use", version})
				Expect(err).Should(BeNil())
				linkDest, _ := fs.Readlink(currentPath)
				Expect(linkDest).Should(Equal(versionSdkDir))
				// ExpectFoldersOrdering(fs, pathr.SDK("go", "current"), []string{"bin", "src", "doc"}, nil, false)
				// Expect(err).Should(BeNil())
			})
		})

		Context("a broken symlink exists for @current", func() {

			BeforeEach(func() {
				_ = fs.MkdirAll(pathr.SDK("go"), os.ModePerm)
				_ = fs.MkdirAll(currentPath, os.ModePerm)
				err := fs.Symlink(currentPath, pathr.SDK("go", "19.5"))
				Expect(err).Should(BeNil())
			})

			It("removes broken symlink and replaces it with <version>", func() {
				err := sut.ExecuteCommand(context.Background(), rootArgument, []string{"use", "1.16.3"})
				Expect(err).Should(BeNil())

				currentDir, err := fs.Open(pathr.SDK("go", "current"))
				Expect(err).Should(BeNil())
				Expect(currentDir).Should(And(BeASymlink(fs), Not(BeADirectory())))
				newDir, err := fs.Open(pathr.SDK("go", "1.16.3"))
				Expect(err).Should(BeNil())
				Expect(newDir).Should(And(BeADirectory(), Not(BeASymlink(fs))))

				oldDir, err := fs.Open(pathr.SDK("go", "19.5"))
				Expect(err).Should(BeNil())
				Expect(oldDir).Should(Not(Or(BeADirectory(), BeASymlink(fs))))
			})
		})
	})

})

/*

	Context("Go SDK Plugin - Use", func() {
		BeforeEach(func() {
			v.MkdirAll(goSdkDir, os.ModePerm)
			fs.MkdirAll(path.Join(versionSdkDir, "bin"), os.ModePerm)
			fs.MkdirAll(path.Join(versionSdkDir, "src"), os.ModePerm)
			fs.MkdirAll(path.Join(versionSdkDir, "doc"), os.ModePerm)
		})

		It("The new SDK Version is symlinked to @current Version ", func() {
			err := sut.ExecuteCommand(context.Background(), "/Users/dev/go/src/github.com/alex-held/devctl/cmd/devctl", []string{"use", version})
			ExpectFolders(fs, path.Join(goSdkDir, "current"), []string{"bin", "src", "doc"}, nil)
			Expect(err).Should(BeNil())
		})
	})*/

/*
var _ = Describe("memory filesystem", func() {
	var fs vfs.FileSystem

	BeforeEach(func() {
		fs = memoryfs.New()
	})

	test.StandardTest(memoryfs.New)

	Context("rename", func() {
		BeforeEach(func() {
			fs.MkdirAll("d1/d1n1/d1n1a", os.ModePerm)
			fs.MkdirAll("d1/d1n2", os.ModePerm)
		})
		It("rename top level", func() {
			Expect(fs.Rename("/d1", "d2")).To(Succeed())
			ExpectFolders(fs, "d2/", []string{"d1n1", "d1n2"}, nil)
		})
		It("rename sub level", func() {
			Expect(fs.Rename("/d1/d1n1", "d2")).To(Succeed())
			ExpectFolders(fs, "d2/", []string{"d1n1a"}, nil)
		})
		It("fail rename root", func() {
			Expect(fs.Rename("/", "d2")).To(Equal(errors.New("cannot rename root dir")))
		})
		It("fail rename to existent", func() {
			Expect(fs.Rename("/d1/d1n1", "d1/d1n2")).To(Equal(os.ErrExist))
		})

		It("rename link", func() {
			Expect(fs.MkdirAll("d2", os.ModePerm)).To(Succeed())
			Expect(fs.Symlink("/d1/d1n1", "d2/link")).To(Succeed())
			Expect(fs.Rename("d2/link", "d2/new")).To(Succeed())
			ExpectFolders(fs, "d2", []string{"new"}, nil)
			ExpectFolders(fs, "d2/new", []string{"d1n1a"}, nil)
		})
	})
})*/

func TestGoUseCmd_Link(t *testing.T) {
	const rootPath = "/home/user1/.devctl"
	const version = "1.0.0"

	RegisterFailHandler(Fail)
	RunSpecs(t, "TestGoUseCmd_Link")

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
				fs.MkdirAll(goSdkDir, os.ModePerm)
				fs.MkdirAll(path.Join(versionSdkDir, "bin"), os.ModePerm)
				fs.MkdirAll(path.Join(versionSdkDir, "src"), os.ModePerm)
				fs.MkdirAll(path.Join(versionSdkDir, "doc"), os.ModePerm)
			})

			It("The new GOROOT is symlinked to $SDKPATH/go/current", func() {
				var err = sut.ExecuteCommand(nil, "", []string{version})
				Expect(err).Should(BeNil())
				ExpectFolders(fs, path.Join(goSdkDir, "current"), []string{"bin", "src", "doc"}, nil)
			})
		})
	})
}

func SymlinkTestcase(fs vfs.FileSystem) {
	Context("symlinks", func() {
		BeforeEach(func() {
			fs.MkdirAll("d1/d1n1/d1n1a", os.ModePerm)
			fs.MkdirAll("d1/d1n2", os.ModePerm)
			fs.MkdirAll("d2/d2n1", os.ModePerm)
			fs.MkdirAll("d2/d2n2", os.ModePerm)
		})

		It("creates link", func() {
			Expect(fs.Symlink("/d1/d1n1", "d2/link")).To(Succeed())
			ExpectFolders(fs, "d2", []string{"d2n1", "d2n2", "link"}, nil)
			Expect(fs.Readlink("/d2/link")).To(Equal("/d1/d1n1"))
			ExpectFolders(fs, "d2/link", []string{"d1n1a"}, nil)
		})

		It("lstat link", func() {
			Expect(fs.Symlink("/d1/d1n1", "d2/link")).To(Succeed())
			fi, err := fs.Lstat("d2/link")
			Expect(err).To(Succeed())
			Expect(fi.Mode() & os.ModeType).To(Equal(os.ModeSymlink))
		})

		It("stat link", func() {
			Expect(fs.Symlink("/d1/d1n1", "d2/link")).To(Succeed())
			fi, err := fs.Stat("d2/link")
			Expect(err).To(Succeed())
			Expect(fi.Mode() & os.ModeType).To(Equal(os.ModeDir))
		})

		It("remove link", func() {
			Expect(fs.Symlink("/d1/d1n1", "d2/link")).To(Succeed())
			Expect(fs.Remove("d2/link")).To(Succeed())
			ExpectFolders(fs, "d1", []string{"d1n1", "d1n2"}, nil)
			ExpectFolders(fs, "d2", []string{"d2n1", "d2n2"}, nil)
		})

		Context("eval", func() {
			It("plain", func() {
				Expect(fs.Symlink("/d1/d1n1", "d2/link")).To(Succeed())
				ExpectFolders(fs, "d2/link", []string{"d1n1a"}, nil)
			})
			It("dotdot", func() {
				Expect(fs.Symlink("/d1/d1n1", "d2/link")).To(Succeed())
				ExpectFolders(fs, "d2/link/..", []string{"d1n1", "d1n2"}, nil)
			})
			It("dotdot in link", func() {
				Expect(fs.Symlink("../d1", "d2/link")).To(Succeed())
				ExpectFolders(fs, "d2/link", []string{"d1n1", "d1n2"}, nil)
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

//BeADirectory succeeds if a file exists and is a directory.
//Actual must be a string representing the abs path to the file being checked.
func BeASymlink(fs vfs.VFS) types.GomegaMatcher {
	return &testutils.BeASymlinkMatcher{
		Fs: fs,
	}
}
