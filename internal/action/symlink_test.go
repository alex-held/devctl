package action

import (
	"os"
	"os/exec"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/logging"

	"github.com/alex-held/devctl/internal/devctlpath"
)

//nolint:gocognit
func TestSymlink_LinkCurrentSDKF(t *testing.T) {
	g := goblin.Goblin(t)

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("GIVEN not yet linked", func() {
		var fs afero.Fs
		var fixture *ActionTestFixture
		var sut *Symlink
		var logger *logging.Logger
		const sdk = "scala"
		const version = "2.13.4"
		var installedSdkPath string

		g.JustBeforeEach(func() {
			fs = afero.NewOsFs()
			logger = logging.NewLogger(
				logging.WithVerbose(true),
				logging.WithFormatter(&logrus.TextFormatter{
					ForceColors:      true,
					DisableTimestamp: true,
					PadLevelText:     true,
					QuoteEmptyFields: true,
				}),
			)
			dir, err := afero.TempDir(fs, "", "devctl_symlink_test")
			if err != nil {
				g.Fatalf("failed to create tmp dir in setup; error=%v", err)
			}
			logger.Tracef("Tempdir=%s", dir)
			teardown := func() {
				err := fs.RemoveAll(dir)
				if err != nil {
					g.Fatalf("failed to clean up temp dir in teardown; tmpdir=%s; error=%v", dir, err)
				}
			}

			pather := devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
				return dir
			}))
			installedSdkPath = pather.SDK(sdk, version)
			e := fs.MkdirAll(installedSdkPath, 0700|os.ModeDir)
			if e != nil {
				g.Fatalf("failed to create sdks/scala/2.13.4 dir in setup; tmp=%s; error=%v", installedSdkPath, e)
			}
			fixture = SetupFixtureDeps(g, fs, pather, logger, teardown)
			sut = fixture.actions.Symlink
		})

		g.AfterEach(func() {
			fixture.teardown()
		})

		g.It("THEN symlinks sdks/scala/2.13.4 to sdks/scala/current", func() {
			current, err := sut.LinkCurrentSDK(sdk, version)
			Expect(err).Should(BeNil())
			Expect(current).ShouldNot(BeNil())
		})
	})

	g.Describe("GIVEN current already linked", func() {
		var fs afero.Fs
		var fixture *ActionTestFixture
		var sut *Symlink
		var logger *logging.Logger
		const sdk = "scala"
		const version = "2.13.4"
		const alreadyLinkedVersion = "2.10.0"
		var installedSdkPath string
		var alreadyLinkedSdkPath string

		g.AfterEach(func() {
			fixture.teardown()
		})

		g.JustBeforeEach(func() {
			fs = afero.NewOsFs()
			logger = logging.NewLogger(
				logging.WithVerbose(true),
				logging.WithFormatter(&logrus.TextFormatter{
					ForceColors:      true,
					DisableTimestamp: true,
					PadLevelText:     true,
					QuoteEmptyFields: true,
				}),
			)
			dir, err := afero.TempDir(fs, "", "devctl_symlink_test")
			if err != nil {
				g.Fatalf("failed to create tmp dir in setup; error=%v", err)
			}
			logger.Tracef("Tempdir=%s", dir)
			teardown := func() {
				err = fs.RemoveAll(dir)
				if err != nil {
					g.Fatalf("failed to clean up temp dir in teardown; tmpdir=%s; error=%v", dir, err)
				}
			}

			pather := devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
				return dir
			}))

			var mkdirPaths []string
			currentPath := pather.SDK(sdk, "current")
			alreadyLinkedSdkPath = pather.SDK(sdk, alreadyLinkedVersion)
			installedSdkPath = pather.SDK(sdk, version)
			mkdirPaths = append(mkdirPaths, alreadyLinkedSdkPath)

			for _, p := range mkdirPaths {
				e := fs.MkdirAll(p, 0700|os.ModeDir)
				if e != nil {
					g.Fatalf("failed to create sdks/scala/2.13.4 dir in setup; tmp=%s; error=%v", installedSdkPath, e)
				}
			}
			cmd := exec.Command("ln", "-s", alreadyLinkedSdkPath, currentPath)
			err = cmd.Run()
			if err != nil {
				g.Fatalf("failed to setup pre-existing symlink; alreadyLinkedSdk=%s; current=%s", alreadyLinkedVersion, currentPath)
			}

			fixture = SetupFixtureDeps(g, fs, pather, logger, teardown)
			sut = fixture.actions.Symlink
		})

		g.It("THEN replaces existing symlink with symlink from sdks/scala/2.13.4 to sdks/scala/current", func() {
			current, err := sut.LinkCurrentSDK(sdk, version)
			Expect(err).Should(BeNil())
			readlink, err := os.Readlink(current)
			Expect(err).Should(BeNil())
			Expect(readlink).ShouldNot(Equal(alreadyLinkedSdkPath))
			Expect(readlink).Should(Equal(installedSdkPath))
			Expect(current).ShouldNot(BeNil())
		})
	})
}
