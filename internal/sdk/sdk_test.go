package sdk

import (
	"path"
	"testing"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/franela/goblin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/pkg/mocks"
)

const sdkpath = "/devctl/sdks"

func TestList(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	c := gomock.NewController(t)
	mockPather := mocks.NewMockPather(c)
	mockPather.EXPECT().SDK(gomock.Any()).Return(sdkpath).AnyTimes()

	g.Describe("ListSDK", func() {
		var sut SDKActions
		var fs afero.Fs

		g.BeforeEach(func() {
			fs = afero.NewMemMapFs()
			err := fs.MkdirAll(sdkpath, fileutil.PrivateDirMode)
			if err != nil {
				g.Fail(err)
			}
			sut = SDKActions{
				Pather: mockPather,
				FS:     fs,
			}
		})

		g.It("WHEN no sdks installed => THEN returns no sdks", func() {
			sdks, err := sut.List()
			Expect(err).Should(BeNil())
			Expect(sdks).Should(BeEmpty())
		})

		g.It("WHEN sdks installed => THEN returns these sdks", func() {
			expectedSdks := []string{"go", "scala", "haskell"}
			for _, sdk := range expectedSdks {
				err := fs.MkdirAll(path.Join(sdkpath, sdk), fileutil.PrivateDirMode)
				if err != nil {
					g.Fail(err)
				}
			}
			sdks, err := sut.List()
			Expect(err).Should(BeNil())
			Expect(sdks).Should(ContainElements(expectedSdks))
		})
	})
}

func TestListVersions(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	c := gomock.NewController(t)
	mockPather := mocks.NewMockPather(c)
	mockPather.EXPECT().SDK(gomock.Eq("go")).Return(path.Join(sdkpath, "go")).AnyTimes()

	g.Describe("ListSDKVersions", func() {
		var sut SDKActions
		var fs afero.Fs
		var gosdkPath string
		g.BeforeEach(func() {
			fs = afero.NewMemMapFs()
			err := fs.MkdirAll(gosdkPath, fileutil.PrivateDirMode)
			if err != nil {
				g.Fail(err)
			}
			sut = SDKActions{
				Pather: mockPather,
				FS:     fs,
			}
		})

		g.It("WHEN sdk not installed => THEN returns no versions", func() {
			sdks, err := sut.ListVersions("go")
			Expect(err).Should(BeNil())
			Expect(sdks).Should(BeEmpty())
		})

		g.It("WHEN no sdk version installed => THEN returns no versions", func() {
			sdks, err := sut.ListVersions("go")
			gosdkPath = path.Join(sdkpath, "go")
			err = fs.MkdirAll(gosdkPath, fileutil.PrivateDirMode)
			if err != nil {
				g.Fail(err)
			}
			Expect(err).Should(BeNil())
			Expect(sdks).Should(BeEmpty())
		})

		g.It("WHEN sdk versions installed => THEN returns these versions", func() {
			gosdkPath = path.Join(sdkpath, "go")
			expectedVersions := []string{"1.15", "1.16", "1.16.2"}
			err := fs.MkdirAll(gosdkPath, fileutil.PrivateDirMode)
			if err != nil {
				g.Fail(err)
			}
			for _, sdk := range expectedVersions {
				err := fs.MkdirAll(path.Join(gosdkPath, sdk), fileutil.PrivateDirMode)
				if err != nil {
					g.Fail(err)
				}
			}
			sdks, err := sut.ListVersions("go")
			Expect(err).Should(BeNil())
			Expect(sdks).Should(ContainElements(expectedVersions))
		})
	})
}
