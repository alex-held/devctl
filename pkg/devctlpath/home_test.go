package devctlpath

import (
	"fmt"
	. "math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"

	"github.com/alex-held/devctl/internal/system"

	"github.com/coreos/etcd/pkg/stringutil"
	_ "github.com/onsi/gomega/matchers"

	"github.com/alex-held/devctl/pkg/logging"
)

func TestHomeFinder(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	var lazyFinder Pather

	customUserHomeFn := func() string {
		return "/h/o/m/e/user"
	}

	customCacheFn := func() string {
		return "/c/a/c/h/e/"
	}

	const testAppPrefix = "test_devctl"
	var testAppPrefixWithLeadingDot = fmt.Sprintf(".%s", testAppPrefix)
	var customConfigRoot = fmt.Sprintf("/h/o/m/e/user/%s", testAppPrefix)

	/* ConfigRoot */
	g.Describe("ConfigRoot", func() {
		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := resolveConfigSubDir(userHome, testAppPrefix)
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.ConfigRoot()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			userHome := customUserHomeFn()
			expected := resolveConfigSubDir(userHome, testAppPrefixWithLeadingDot)
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithUserHomeFn(customUserHomeFn))
			actual := lazyFinder.ConfigRoot()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := customConfigRoot
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return expected
			}))
			actual := lazyFinder.ConfigRoot()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := customConfigRoot
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return expected
			}))
			actual := lazyFinder.ConfigRoot()
			Expect(actual).To(Equal(expected))
		})
	})

	/* CacheDir */
	g.Describe("Cache", func() {
		g.It("WHEN no pathFn set", func() {
			_ = os.Setenv("XDG_CACHE_HOME", "/tmp/cache")
			cacheDir, _ := os.UserCacheDir()
			expected := filepath.Join(cacheDir, "io.alexheld.test_devctl")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.Cache()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN cacheFn set", func() {
			expected := filepath.Join(customCacheFn(), "io.alexheld.test_devctl")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithCachePathFn(customCacheFn))
			actual := lazyFinder.Cache()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing path elems", func() {
			expected := filepath.Join(customCacheFn(), "io.alexheld.test_devctl/some/sub/dir")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithCachePathFn(customCacheFn), WithConfigRootFn(func() string {
				return expected
			}))
			actual := lazyFinder.Cache("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customCacheFn(), "io.alexheld.test_devctl")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithCachePathFn(customCacheFn), WithConfigRootFn(func() string {
				return expected
			}))
			actual := lazyFinder.Cache()
			Expect(actual).To(Equal(expected))
		})
	})

	/* Bin */
	g.Describe("Bin", func() {
		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "bin")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			userHome := customUserHomeFn()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "bin")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithUserHomeFn(customUserHomeFn))
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := filepath.Join(customConfigRoot, "bin")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customConfigRoot, "bin")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := filepath.Join(customConfigRoot, "/bin/some/sub/dir")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Bin("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})

	/* SDK */
	g.Describe("SDK", func() {
		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "sdks")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			userHome := customUserHomeFn()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "sdks")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithUserHomeFn(customUserHomeFn))
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := filepath.Join(customConfigRoot, "sdks")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customConfigRoot, "sdks")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := filepath.Join(customConfigRoot, "sdks/some/sub/dir")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.SDK("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})

	/* Config */
	g.Describe("Config", func() {
		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "config")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			userHome := customUserHomeFn()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "config")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithUserHomeFn(customUserHomeFn))
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := filepath.Join(customConfigRoot, "config")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customConfigRoot, "config")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := filepath.Join(customConfigRoot, "config/some/sub/dir")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Config("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})

	/* Download */
	g.Describe("Download", func() {
		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "downloads")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			userHome := customUserHomeFn()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "downloads")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithUserHomeFn(customUserHomeFn))
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := filepath.Join(customConfigRoot, "downloads")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customConfigRoot, "downloads")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := filepath.Join(customConfigRoot, "downloads/some/sub/dir")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Download("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})
}

type ExpectedPaths struct {
	user    string
	prefix  string
	cache   string
	cfgFile string
}

func (e ExpectedPaths) ConfigRoot(elem ...string) string {
	return filepath.Join(filepath.Join(e.user, e.prefix), filepath.Join(elem...))
}

func (e ExpectedPaths) Config(elem ...string) string {
	return filepath.Join(e.ConfigRoot(), "config", filepath.Join(elem...))
}

func (e ExpectedPaths) Bin(elem ...string) string {
	return filepath.Join(e.ConfigRoot(), "bin", filepath.Join(elem...))
}

func (e ExpectedPaths) Download(elem ...string) string {
	return filepath.Join(e.ConfigRoot(), "downloads", filepath.Join(elem...))
}

func (e ExpectedPaths) SDK(elem ...string) string {
	return filepath.Join(e.ConfigRoot(), "sdks", filepath.Join(elem...))
}

func (e ExpectedPaths) Cache(elem ...string) string {
	return filepath.Join(e.cache, fmt.Sprintf("io.alexheld%s", e.prefix), filepath.Join(elem...))
}

func (e ExpectedPaths) ConfigFilePath() string {
	return e.ConfigRoot(e.cfgFile)
}

type TC struct {
	SystemUnderTest Pather
	Exp             *ExpectedPaths
}

func (t *TC) Generate(_ *Rand, _ int) reflect.Value {
	randStrings := stringutil.UniqueStrings(5, 5)

	var user = randStrings[0]
	var prefix = "." + randStrings[1]
	var cache = randStrings[2]
	var cfgFile = randStrings[3] + ".yaml"

	var lpF Pather = &lazypathFinder{
		cfgName: cfgFile,
		lp:      lazypath(prefix),
		finder: finder{
			GetUserHomeFn: func() string {
				return user
			},
			GetCachePathFn: func() string {
				return cache
			},
			GetConfigRootFn: func() string {
				return filepath.Join(user, prefix)
			},
		},
	}

	expected := ExpectedPaths{user, prefix, cache, cfgFile}

	var testCase = &TC{
		SystemUnderTest: lpF,
		Exp:             &expected,
	}

	value := reflect.ValueOf(testCase)
	return value
}

func (t *TC) Pather() Pather           { return t.SystemUnderTest }
func (t *TC) Expected() *ExpectedPaths { return t.Exp }

func NewQuickCheckConfig(iterations, scale int) *quick.Config {
	var source = NewSource(time.Now().Unix())
	r := New(source)
	return &quick.Config{
		MaxCount:      iterations,
		MaxCountScale: float64(scale),
		Rand:          r,
	}
}

type GoblinG struct {
	TestingTB testing.TB
	*goblin.G
	QuickCheckConfig *quick.Config
	Logger           logging.Log
}

func TestPather_QuickCheck(t *testing.T) {
	g := GoblinG{
		TestingTB:        t,
		G:                goblin.Goblin(t),
		QuickCheckConfig: NewQuickCheckConfig(5000, 100),
		Logger:           logging.NewLogger(logging.WithOutputs(), logging.WithLevel(logging.LogLevelError)),
	}

	tt := map[string]struct {
		underTestFn func(Pather) string
		expectedFn  func(*ExpectedPaths) string
	}{
		"Pather.ConfigFilePath()": {
			underTestFn: func(p Pather) string { return p.ConfigFilePath() },
			expectedFn:  func(p *ExpectedPaths) string { return p.ConfigFilePath() },
		},
		"Pather.ConfigRoot()": {
			underTestFn: func(p Pather) string { return p.ConfigRoot() },
			expectedFn:  func(p *ExpectedPaths) string { return p.ConfigRoot() },
		},
		"Pather.Config()": {
			underTestFn: func(p Pather) string { return p.Config() },
			expectedFn:  func(p *ExpectedPaths) string { return p.Config() },
		},
		"Pather.Bin()": {
			underTestFn: func(p Pather) string { return p.Bin() },
			expectedFn:  func(p *ExpectedPaths) string { return p.Bin() },
		},
		"Pather.Download()": {
			underTestFn: func(p Pather) string { return p.Download() },
			expectedFn:  func(p *ExpectedPaths) string { return p.Download() },
		},
		"Pather.SDK()": {
			underTestFn: func(p Pather) string { return p.SDK() },
			expectedFn:  func(p *ExpectedPaths) string { return p.SDK() },
		},
		"Pather.Cache()": {
			underTestFn: func(p Pather) string { return p.Cache() },
			expectedFn:  func(p *ExpectedPaths) string { return p.Cache() },
		},
	}

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("QuickCheck", func() {
		for scenario, test := range tt {
			g.It(scenario, func() {
				f := g.GetQuickCheckFunc(scenario, test.underTestFn, test.expectedFn)
				Expect(quick.Check(f, g.QuickCheckConfig)).To(BeNil(), "!! Error occurred: \n\n")
			})
		}
	})
}

func (g *GoblinG) GetQuickCheckFunc(scenario string, functionUnderTest func(Pather) string, expected func(*ExpectedPaths) string) func(*TC) bool {
	return func(tc *TC) bool {
		expected := expected(tc.Expected())
		actual := functionUnderTest(tc.Pather())
		if testing.Verbose() {
			g.Logger.Infof("Scenario: %s, %v, %v", scenario, tc.SystemUnderTest, actual)
		}
		return Expect(actual).To(Equal(expected))
	}
}

func resolveConfigSubDir(home, prefix string, elem ...string) (path string) {
	arch := system.GetCurrent()
	prefix = "." + strings.TrimPrefix(prefix, ".")
	var cfgRoot string
	switch {
	case arch.IsLinux():
		cfgRoot = filepath.Join(home, ".config", prefix)
	case arch.IsDarwin():
		cfgRoot = filepath.Join(home, prefix)
	default:
		// windows not yet supported
		panic(fmt.Errorf("the current os is not yet supported; os=%s", arch))
	}
	return filepath.Join(cfgRoot, filepath.Join(elem...))
}
