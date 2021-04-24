package plugins

import (
	"github.com/alex-held/devctl/pkg/devctlpath"
	"github.com/alex-held/devctl/pkg/testutils/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/spf13/afero"
)


var _ = Describe("Store", func() {

	var sut Store
	var pather devctlpath.Pather
	var fs afero.Fs
	var manifest string

	BeforeEach(func() {
		pather = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			return "/tmp/devctl"
		}))
		fs = afero.NewMemMapFs()
		sut = &store{
			Pather: pather,
			Fs:     fs,
		}
		manifest = pather.ConfigRoot("plugins.yaml")
	})

	Context("List", func() {

		When("no plugin manifest exists", func() {

			It("returns empty list", func() {
				plugins, err := sut.List(SDK)
				Expect(err).Should(Succeed())
				Expect(plugins).Should(BeEmpty())
			})
		})

		When("empty plugin manifest exists", func() {

			BeforeEach(func() {
				_, _ = fs.Create(manifest)
			})

			It("returns empty list", func() {
				plugins, err := sut.List(SDK)
				Expect(err).Should(Succeed())
				Expect(plugins).Should(BeEmpty())
			})
		})

		When("plugin manifest with empty categories exists", func() {

			BeforeEach(func() {
				data, _ := afero.ReadFile(afero.NewOsFs(), "testdata/plugins_empty.yaml")
				_ = afero.WriteFile(fs, manifest, data, 0777)
			})

			It("returns empty list", func() {
				plugins, err := sut.List(SDK)
				Expect(err).Should(Succeed())
				Expect(plugins).Should(BeEmpty())
			})
		})

		When("plugin manifest with content exists", func() {

			BeforeEach(func() {
				data, _ := afero.ReadFile(afero.NewOsFs(), "testdata/plugins_1.yaml")
				_ = afero.WriteFile(fs, manifest, data, 0777)
			})

			It("lists sdk plugins", func() {
				plugins, err := sut.List(SDK)
				Expect(err).Should(Succeed())
				Expect(plugins).Should(Equal([]string{"go"}))
			})
		})
	})

	Context("Register", func() {

		When("no plugin manifest exists", func() {

			It("creates the plugin manifest", func() {
				err := sut.Register(SDK, "scala")
				Expect(err).Should(Succeed())
				Expect(manifest).Should(matchers.BeAnExistingFileFs(fs))
			})
		})

		When("empty plugin manifest exists", func() {

			BeforeEach(func() {
				_, _ = fs.Create(manifest)
			})

			It("creates the corresponding category with the registered plugin", func() {
				err := sut.Register(SDK, "scala")
				Expect(err).Should(Succeed())
				Expect(sut).Should(ContainsPluginForKind(SDK, "scala"))
			})
		})

		When("plugin manifest with empty categories exists", func() {

			BeforeEach(func() {
				data, _ := afero.ReadFile(afero.NewOsFs(), "testdata/plugins_empty.yaml")
				_ = afero.WriteFile(fs, manifest, data, 0777)
			})

			It("appends the registered plugin to the corresponding category", func() {
				err := sut.Register(SDK, "scala")
				Expect(err).Should(Succeed())
				Expect(sut).Should(ContainsPluginForKind(SDK, "scala"))
			})
		})

		When("plugin manifest with content exists", func() {

			BeforeEach(func() {
				data, _ := afero.ReadFile(afero.NewOsFs(), "testdata/plugins_1.yaml")
				_ = afero.WriteFile(fs, manifest, data, 0777)
			})

			It("appends the registered plugin to the corresponding category", func() {
				plugins, err := sut.List(SDK)
				Expect(err).Should(Succeed())
				Expect(plugins).Should(Equal([]string{"go"}))
			})
		})

	})
})

func ContainsPluginForKind(kind Kind, name string) types.GomegaMatcher {
	type wrapper struct {
		Plugins []string
		Err     error
	}
	return WithTransform(func(store *store) wrapper {
		plugins, err := store.List(kind)
		return wrapper{
			Plugins: plugins,
			Err:     err,
		}
	}, SatisfyAll(
		WithTransform(func(w wrapper) error { return w.Err }, SatisfyAll(Succeed())),
		WithTransform(func(w wrapper) []string { return w.Plugins }, SatisfyAll(ContainElement(name))),
	))
}

func assertPluginManifestExists(fs afero.Fs, path string)  {
	Expect(path).Should(matchers.BeAnExistingFileFs(fs))
}