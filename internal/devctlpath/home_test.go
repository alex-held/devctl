package devctlpath

import (
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"testing"
)

func TestHomeFinder(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	var userHomeFn UserHomePathFinder
	var cacheFn CachePathFinder
	var configRootFn ConfigRootFinder
	var lazyFinder lazypathFinder

	customUserHomeFn := func() string {
		return "/h/o/m/e/user"
	}

	customCacheFn := func() string {
		return "/c/a/c/h/e/"
	}

	g.Describe("DevCtlConfigRoot", func() {

		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := filepath.Join(userHome, ".test_devctl")
			lazyFinder = NewLazyFinder("test_devctl", nil, nil, nil)
			actual := lazyFinder.DevCtlConfigRoot()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			expected := filepath.Join(customUserHomeFn(), ".test_devctl")
			lazyFinder = NewLazyFinder("test_devctl", customUserHomeFn, nil, nil)
			actual := lazyFinder.DevCtlConfigRoot()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := filepath.Join(userHome, ".test_devctl")
			lazyFinder = NewLazyFinder("test_devctl", userHomeFn, cacheFn, configRootFn)
			actual := lazyFinder.DevCtlConfigRoot()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := "/h/o/m/e/user/.test_devctl"
			lazyFinder = NewLazyFinder("test_devctl", nil, nil, func() string {
				return expected
			})
			actual := lazyFinder.DevCtlConfigRoot()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := "/h/o/m/e/user/.test_devctl"
			lazyFinder = NewLazyFinder(".test_devctl", nil, nil, func() string {
				return expected
			})
			actual := lazyFinder.DevCtlConfigRoot()
			Expect(actual).To(Equal(expected))
		})

	})

	/* CacheDir */
	g.Describe("Cache", func() {

		g.It("WHEN no pathFn set", func() {
			cacheDir, _ := os.UserCacheDir()
			expected := filepath.Join(cacheDir, "io.alexheld.test_devctl")
			lazyFinder = NewLazyFinder("test_devctl", nil, nil, nil)
			actual := lazyFinder.Cache()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN cacheFn set", func() {
			expected := filepath.Join(customCacheFn(), "io.alexheld.test_devctl")
			lazyFinder = NewLazyFinder("test_devctl", nil, customCacheFn, nil)
			actual := lazyFinder.Cache()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing path elems", func() {
			expected := filepath.Join(customCacheFn(), "io.alexheld.test_devctl/some/sub/dir")
			lazyFinder = NewLazyFinder("test_devctl", nil, customCacheFn, func() string {
				return expected
			})
			actual := lazyFinder.Cache("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customCacheFn(), "io.alexheld.test_devctl")
			lazyFinder = NewLazyFinder(".test_devctl", nil, customCacheFn, func() string {
				return expected
			})
			actual := lazyFinder.Cache()
			Expect(actual).To(Equal(expected))
		})

	})


	g.Describe("Bin", func() {

		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := filepath.Join(userHome, ".test_devctl/bin")
			lazyFinder = NewLazyFinder("test_devctl", nil, nil, nil)
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			expected := filepath.Join(customUserHomeFn(), ".test_devctl/bin")
			lazyFinder = NewLazyFinder("test_devctl", customUserHomeFn, nil, nil)
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := filepath.Join(userHome, ".test_devctl/bin")
			lazyFinder = NewLazyFinder("test_devctl", userHomeFn, cacheFn, configRootFn)
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := "/h/o/m/e/user/.test_devctl/bin"
			lazyFinder = NewLazyFinder("test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := "/h/o/m/e/user/.test_devctl/bin"
			lazyFinder = NewLazyFinder(".test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := "/h/o/m/e/user/.test_devctl/bin/some/sub/dir"
			lazyFinder = NewLazyFinder(".test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.Bin("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})



	g.Describe("SDK", func() {

		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := filepath.Join(userHome, ".test_devctl/sdks")
			lazyFinder = NewLazyFinder("test_devctl", nil, nil, nil)
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			expected := filepath.Join(customUserHomeFn(), ".test_devctl/sdks")
			lazyFinder = NewLazyFinder("test_devctl", customUserHomeFn, nil, nil)
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := filepath.Join(userHome, ".test_devctl/sdks")
			lazyFinder = NewLazyFinder("test_devctl", userHomeFn, cacheFn, configRootFn)
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := "/h/o/m/e/user/.test_devctl/sdks"
			lazyFinder = NewLazyFinder("test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := "/h/o/m/e/user/.test_devctl/sdks"
			lazyFinder = NewLazyFinder(".test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := "/h/o/m/e/user/.test_devctl/sdks/some/sub/dir"
			lazyFinder = NewLazyFinder(".test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.SDK("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})

	g.Describe("Config", func() {

		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := filepath.Join(userHome, ".test_devctl/config")
			lazyFinder = NewLazyFinder("test_devctl", nil, nil, nil)
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			expected := filepath.Join(customUserHomeFn(), ".test_devctl/config")
			lazyFinder = NewLazyFinder("test_devctl", customUserHomeFn, nil, nil)
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := filepath.Join(userHome, ".test_devctl/config")
			lazyFinder = NewLazyFinder("test_devctl", userHomeFn, cacheFn, configRootFn)
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := "/h/o/m/e/user/.test_devctl/config"
			lazyFinder = NewLazyFinder("test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := "/h/o/m/e/user/.test_devctl/config"
			lazyFinder = NewLazyFinder(".test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := "/h/o/m/e/user/.test_devctl/config/some/sub/dir"
			lazyFinder = NewLazyFinder(".test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.Config("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})

	g.Describe("Download", func() {

		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := filepath.Join(userHome, ".test_devctl/downloads")
			lazyFinder = NewLazyFinder("test_devctl", nil, nil, nil)
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			expected := filepath.Join(customUserHomeFn(), ".test_devctl/downloads")
			lazyFinder = NewLazyFinder("test_devctl", customUserHomeFn, nil, nil)
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := filepath.Join(userHome, ".test_devctl/downloads")
			lazyFinder = NewLazyFinder("test_devctl", userHomeFn, cacheFn, configRootFn)
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := "/h/o/m/e/user/.test_devctl/downloads"
			lazyFinder = NewLazyFinder("test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := "/h/o/m/e/user/.test_devctl/downloads"
			lazyFinder = NewLazyFinder(".test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := "/h/o/m/e/user/.test_devctl/downloads/some/sub/dir"
			lazyFinder = NewLazyFinder(".test_devctl", nil, nil, func() string {
				return "/h/o/m/e/user/.test_devctl"
			})
			actual := lazyFinder.Download("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})

}
