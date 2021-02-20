// +build !windows

package cli

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/alex-held/devctl/internal/system"

	"github.com/bxcodec/faker/v3"
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

func GetTestGetHomeFunc(os, username string) ConfigGetter {
	return func() string {
		switch os {
		case system.OsDarwin:
			return fmt.Sprintf(" /Users/%s", username)
		case system.OsWindows:
			panic(errors.Errorf("windows not supported yet.."))
		default: // linux
			return fmt.Sprintf("/home/%s", username)
		}
	}
}

func GetTestGetEnvFunc(user string) EnvGetter {
	return func(e string) string {
		if e == "HOME" {
			return fmt.Sprintf("/home/%s", user)
		}
		panic(errors.Errorf("No test setup configured for env var: $%s", e))
	}
}

func setupHomeTest(os string) (expectedHome string, hf HomeFinder) {
	username := faker.Username()
	prefix := fmt.Sprintf("devctl_test_%s", os)

	expectedHome = GetExpectedHome(os, username, prefix)

	hf = NewHomeFinderForOS(os, prefix, GetTestGetHomeFunc(os, username), GetTestGetEnvFunc(username))

	return expectedHome, hf
}

func GetExpectedHome(os, username, prefix string) (expectedHome string) {
	switch os {
	case "darwin":
		expectedHome = fmt.Sprintf("%s/.%s", GetTestGetHomeFunc(os, username)(), prefix)
	case "windows":
		panic(errors.Errorf("windows not supported yet.."))
	default:
		expectedHome = fmt.Sprintf("%s/.config/.%s", GetTestGetHomeFunc(os, username)(), prefix)
	}
	return expectedHome
}

func TestHomeFinder(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("HomeFinder", func() {
		var expectedHome string
		var hf HomeFinder

		g.Describe("GIVEN GOOS=linux", func() {
			g.JustBeforeEach(func() {
				expectedHome, hf = setupHomeTest("linux")
			})
			g.AfterEach(func() {
				expectedHome, hf = "<nil>", nil
			})

			g.It("WHEN HomeFinder.Home() => THEN returns /home/<username>/.config/.<app-prefix>", func() {
				actual := hf.Home()
				Expect(actual).To(Equal(expectedHome))
			})

			g.It("WHEN HomeFinder.BinDir() => THEN returns /home/<username>/.config/.<app-prefix>/bin", func() {
				expected := filepath.Join(expectedHome, "bin")
				actual := hf.BinDir()
				Expect(actual).To(Equal(expected))
			})

			g.It("WHEN HomeFinder.ConfigDir() => THEN returns /home/<username>/.config/.<app-prefix>/config", func() {
				expected := filepath.Join(expectedHome, "config")
				actual := hf.ConfigDir()
				Expect(actual).To(Equal(expected))
			})

			g.It("WHEN HomeFinder.DownloadsDir() => THEN returns /home/<username>/.config/.<app-prefix>/downloads", func() {
				expected := filepath.Join(expectedHome, "downloads")
				actual := hf.DownloadsDir()
				Expect(actual).To(Equal(expected))
			})

			g.Describe("WHEN HomeFinder.SDKDir(sdk string)", func() {
				for _, sdk := range []string{"scala", "go", "java"} {
					testcase := fmt.Sprintf("THEN returns /Users/<username>/.<app-prefix>/sdks/%s", sdk)
					g.It(testcase, func() {
						expected := filepath.Join(expectedHome, "sdks", sdk)
						actual := hf.SDKDir(sdk)
						Expect(actual).To(Equal(expected))
					})
				}
			})

			g.It("WHEN HomeFinder.SDKRoot() => THEN returns /home/<username>/.config/.<app-prefix>/sdks", func() {
				expected := filepath.Join(expectedHome, "sdks")
				actual := hf.SDKRoot()
				Expect(actual).To(Equal(expected))
			})

			g.It("WHEN HomeFinder.LogDir() => THEN returns /home/<username>/.config/.<app-prefix>/logs", func() {
				expected := filepath.Join(expectedHome, "logs")
				actual := hf.LogDir()
				Expect(actual).To(Equal(expected))
			})
		})

		g.Describe("GIVEN GOOS=darwin", func() {
			g.BeforeEach(func() {
				expectedHome, hf = setupHomeTest("darwin")
			})

			g.It("WHEN HomeFinder.Home() => THEN returns /Users/<username>/.<app-prefix>", func() {
				actual := hf.Home()
				Expect(actual).To(Equal(expectedHome))
			})

			g.It("WHEN HomeFinder.BinDir() => THEN returns /Users/<username>/.<app-prefix>/bin", func() {
				expected := filepath.Join(expectedHome, "bin")
				actual := hf.BinDir()
				Expect(actual).To(Equal(expected))
			})

			g.It("WHEN HomeFinder.ConfigDir() => THEN returns /Users/<username>/.<app-prefix>/config", func() {
				expected := filepath.Join(expectedHome, "config")
				actual := hf.ConfigDir()
				Expect(actual).To(Equal(expected))
			})

			g.It("WHEN HomeFinder.DownloadsDir() => THEN returns /Users/<username>/.<app-prefix>/downloads", func() {
				expected := filepath.Join(expectedHome, "downloads")
				actual := hf.DownloadsDir()
				Expect(actual).To(Equal(expected))
			})

			g.Describe("WHEN HomeFinder.SDKDir(sdk string)", func() {
				for _, sdk := range []string{"scala", "go", "java"} {
					testcase := fmt.Sprintf("THEN returns /Users/<username>/.<app-prefix>/sdks/%s", sdk)
					g.It(testcase, func() {
						expected := filepath.Join(expectedHome, "sdks", sdk)
						actual := hf.SDKDir(sdk)
						Expect(actual).To(Equal(expected))
					})
				}
			})

			g.It("WHEN HomeFinder.SDKRoot() => THEN returns /Users/<username>/.<app-prefix>/sdks", func() {
				expected := filepath.Join(expectedHome, "sdks")
				actual := hf.SDKRoot()
				Expect(actual).To(Equal(expected))
			})

			g.It("WHEN HomeFinder.LogDir() => THEN returns /Users/<username>/.<app-prefix>/logs", func() {
				expected := filepath.Join(expectedHome, "logs")
				actual := hf.LogDir()
				Expect(actual).To(Equal(expected))
			})
		})
	})
}
