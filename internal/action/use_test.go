package action

import (
	"os"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/config"

	"github.com/alex-held/devctl/internal/devctlpath"
	"github.com/alex-held/devctl/pkg/logging"
)

func TestUse_UseSDKf(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("UseSDK", func() {
		var fs afero.Fs
		var pather devctlpath.Pather
		var logger logging.Log
		var fixture *ActionTestFixture
		var sut *Use

		const sdk = "scala"
		const useVersion = "1.13.4"
		const oldVersion = "1.10.0"

		g.JustBeforeEach(func() {
			fs = afero.NewOsFs()
			logger = logging.NewLogger()
			tmp, err := afero.TempDir(fs, "", "devctl_use_usesdk_test")
			if err != nil {
				g.Fail(err)
			}
			teardown := func() {
				_ = fs.RemoveAll(tmp)
			}
			pather = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
				return tmp
			}))

			fixture = SetupFixtureDeps(g, fs, pather, logger, teardown)
			sut = fixture.actions.Use
		})

		g.AfterEach(func() {
			fixture.teardown()
		})

		g.It("THEN updates config with sdk.current = version and link installed sdks/<sdk>/<version> to sdks/<sdk>/current", func() {
			oldVersionPath := pather.SDK(sdk, oldVersion)
			useVersionPath := pather.SDK(sdk, useVersion)
			currentPath := pather.SDK(sdk, "current")
			dirs := []string{oldVersionPath, useVersionPath}

			// create necessary dirs
			SetupFs(g, fs, dirs, map[string]string{
				oldVersionPath: currentPath,
			})

			fixture.SetupConfig(func(c *config.Config) {
				c.Sdks[sdk] = config.SdkConfig{
					SDK:     sdk,
					Current: oldVersionPath,
					Installations: map[string]string{
						oldVersion: oldVersionPath,
						useVersion: useVersionPath,
					},
				}
			})

			err := sut.UseSDK("scala", "1.13.4")
			Expect(err).Should(BeNil())
			readlink, err := os.Readlink(currentPath)
			Expect(err).Should(BeNil())
			Expect(readlink).Should(Equal(useVersionPath))
		})
	})
}
