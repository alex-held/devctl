package repo

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/downloader"
	"github.com/alex-held/devctl/internal/meta"
)

const (
	testfile            = "testdata/local-index.yaml"
	annotationstestfile = "testdata/local-index-annotations.yaml"
	chartmuseumtestfile = "testdata/chartmuseum-index.yaml"
	unorderedTestfile   = "testdata/local-index-unordered.yaml"
	testRepo            = "test-repo"
	indexWithDuplicates = `
apiVersion: v1
entries:
  nginx:
    - urls:
        - https://charts.helm.sh/stable/nginx-0.2.0.tgz
      name: nginx
      description: string
      version: 0.2.0
      home: https://github.com/something/else
      digest: "sha256:1234567890abcdef"
  nginx:
    - urls:
        - https://charts.helm.sh/stable/alpine-1.0.0.tgz
        - http://storage2.googleapis.com/kubernetes-charts/alpine-1.0.0.tgz
      name: alpine
      description: string
      version: 1.0.0
      home: https://github.com/something
      digest: "sha256:1234567890abcdef"
`
)

func TestIndexFile_SortEntries(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("IndexFile", func() {
		var i *IndexFile

		g.BeforeEach(func() {
			i = NewIndexFile()
		})

		g.It("MustAdd", func() {
			for _, x := range []struct {
				md       *meta.Metadata
				filename string
				baseURL  string
			}{
				{&meta.Metadata{APIVersion: "v2", Name: "clipper", Version: "0.1.0"}, "clipper-0.1.0.tgz", "http://example.com/charts"},
				{&meta.Metadata{APIVersion: "v2", Name: "cutter", Version: "0.1.1"}, "cutter-0.1.1.tgz", "http://example.com/charts"},
				{&meta.Metadata{APIVersion: "v2", Name: "cutter", Version: "0.1.0"}, "cutter-0.1.0.tgz", "http://example.com/charts"},
				{&meta.Metadata{APIVersion: "v2", Name: "cutter", Version: "0.2.0"}, "cutter-0.2.0.tgz", "http://example.com/charts"},
				{&meta.Metadata{APIVersion: "v2", Name: "setter", Version: "0.1.9+alpha"}, "setter-0.1.9+alpha.tgz", "http://example.com/charts"},
				{&meta.Metadata{APIVersion: "v2", Name: "setter", Version: "0.1.9+beta"}, "setter-0.1.9+beta.tgz", "http://example.com/charts"},
			} {
				if err := i.MustAdd(x.md, x.filename, x.baseURL); err != nil {
					Expect(err).Should(BeNil(), fmt.Sprintf("unexpected error adding to index: %s", err))
				}
			}

			i.SortEntries()

			Expect(i.APIVersion).Should(Equal(APIVersionV1), "Expected API version v1")
			Expect(i.Entries).Should(HaveLen(3), fmt.Sprintf("Expected 3 charts. Got %d", len(i.Entries)))
			Expect(i.Entries["clipper"][0].Name).Should(Equal("clipper"), fmt.Sprintf("Expected clipper, got %s", i.Entries["clipper"][0].Name))

			Expect(i.Entries["cutter"]).Should(HaveLen(3), "Expected three cutters.")
			Expect(i.Entries["cutter"][0].Version).Should(Equal("0.2.0"), fmt.Sprintf("Unexpected first version: %s", i.Entries["cutter"][0].Version))

			Expect(i.Entries["clipper"][0].Name).Should(Equal("clipper"), fmt.Sprintf("Expected clipper, got %s", i.Entries["clipper"][0].Name))

			cv, err := i.Get("setter", "0.1.9")
			Expect(err).Should(BeNil())
			Expect(cv.Metadata.Version).Should(ContainSubstring("0.1.9"), fmt.Sprintf("Unexpected version: %s", cv.Metadata.Version))

			cv, err = i.Get("setter", "0.1.9+alpha")
			Expect(err).Should(BeNil())
			Expect(cv.Metadata.Version).Should(Equal("0.1.9+alpha"), "Expected version: 0.1.9+alpha")
		})
	})
}

func TestLoadIndex(t *testing.T) {
	t.Parallel()
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("IndexFile", func() {
		tests := []struct {
			Name     string
			Filename string
		}{
			{
				Name:     "regular index file",
				Filename: testfile,
			},
			{
				Name:     "chartmuseum index file",
				Filename: chartmuseumtestfile,
			},
		}

		for _, tc := range tests {
			tc := tc

			// g.Timeout(time.Millisecond * 50)
			g.It(tc.Name, func() {
				i, err := LoadIndexFile(tc.Filename)
				if err != nil {
					t.Fatal(err)
				}
				verifyLocalIndex(i)
			})
		}
	})
}

//nolint:gocognit
func verifyLocalIndex(i *IndexFile) {
	Expect(i.Entries).Should(HaveLen(3), fmt.Sprintf("Expected 3 entries in index file but got %d", len(i.Entries)))

	alpine, ok := i.Entries["alpine"]
	Expect(ok).Should(BeTrue(), "'alpine' section not found.")
	Expect(alpine).Should(HaveLen(1), fmt.Sprintf("'alpine' should have 1 chart, got %d", len(i.Entries["alpine"])))

	nginx, ok := i.Entries["nginx"]
	Expect(ok).Should(BeTrue())
	Expect(nginx).Should(HaveLen(2), "Expected 2 nginx entries")

	expects := []*downloader.SDKVersion{
		{
			Metadata: &meta.Metadata{
				APIVersion:  "v2",
				Name:        "alpine",
				Description: "string",
				Version:     "1.0.0",
				Home:        "https://github.com/something",
			},
			URLs: []string{
				"https://charts.helm.sh/stable/alpine-1.0.0.tgz",
				"http://storage2.googleapis.com/kubernetes-charts/alpine-1.0.0.tgz",
			},
		},
		{
			Metadata: &meta.Metadata{
				APIVersion:  "v2",
				Name:        "nginx",
				Description: "string",
				Version:     "0.2.0",
				Home:        "https://github.com/something/else",
			},
			URLs: []string{
				"https://charts.helm.sh/stable/nginx-0.2.0.tgz",
			},
		},
		{
			Metadata: &meta.Metadata{
				APIVersion:  "v2",
				Name:        "nginx",
				Description: "string",
				Version:     "0.1.0",
				Home:        "https://github.com/something",
			},
			URLs: []string{
				"https://charts.helm.sh/stable/nginx-0.1.0.tgz",
			},
		},
	}

	tests := []*downloader.SDKVersion{alpine[0], nginx[0], nginx[1]}

	for i, tt := range tests {
		expect := expects[i]

		if tt.Name != expect.Name {
			Expect(tt.Name).Should(Equal(expect.Name))
		}
		if tt.Description != expect.Description {
			Expect(tt.Description).Should(Equal(expect.Description))
		}
		if tt.Version != expect.Version {
			Expect(tt.Version).Should(Equal(expect.Version))
		}

		if tt.Home != expect.Home {
			Expect(tt.Home).Should(Equal(expect.Home))
		}

		for i, url := range tt.URLs {
			if url != expect.URLs[i] {
				Expect(url).Should(Equal(expect.URLs[i]))
			}
		}
	}
}

func TestNewIndexFile(t *testing.T) {
	g := goblin.Goblin(t)

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("NewIndexFile", func() {
		g.It("WHEN", func() {
			before := time.Now()
			actual := NewIndexFile()
			Expect(actual.Generated).To(BeTemporally(">", before))
			Expect(actual.APIVersion).Should(Equal(APIVersionV1), "the default api version should be 'v1'")
			Expect(actual.Annotations).Should(BeEmpty(), "new index file should not contain annotations")
			Expect(actual.Entries).Should(BeEmpty(), "new index file should not contain entries")
		})
	})
}

func TestIndexFile_MustAdd(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	tts := []struct {
		md       *meta.Metadata
		filename string
		baseURL  string
	}{

		{&meta.Metadata{APIVersion: "v2", Name: "clipper", Version: "0.1.0"}, "clipper-0.1.0.tgz", "http://example.com/charts"},
		{&meta.Metadata{APIVersion: "v2", Name: "alpine", Version: "0.1.0"}, "/home/charts/alpine-0.1.0.tgz", "http://example.com/charts"},
		{&meta.Metadata{APIVersion: "v2", Name: "deis", Version: "0.1.0"}, "/home/charts/deis-0.1.0.tgz", "http://example.com/charts/"},
	}

	g.Describe("MustAdd", func() {
		g.It("WHEN IndexFile.MustAdd(...) called 3 times => THEN adds 3 Entries", func() {
			i := NewIndexFile()

			for _, x := range tts {
				err := i.MustAdd(x.md, x.filename, x.baseURL)
				Expect(err).Should(BeNil(), fmt.Sprintf("unexpected error adding to index: %s", err))
			}

			Expect(i.Entries).Should(HaveLen(3))
			Expect(i.Entries["clipper"][0].URLs[0]).Should(Equal("http://example.com/charts/clipper-0.1.0.tgz"), fmt.Sprintf("Expected http://example.com/charts/clipper-0.1.0.tgz, got %s", i.Entries["clipper"][0].URLs[0]))
			Expect(i.Entries["alpine"][0].URLs[0]).Should(Equal("http://example.com/charts/alpine-0.1.0.tgz"), fmt.Sprintf("Expected http://example.com/charts/alpine-0.1.0.tgz, got %s", i.Entries["alpine"][0].URLs[0]))
			Expect(i.Entries["deis"][0].URLs[0]).Should(Equal("http://example.com/charts/deis-0.1.0.tgz"), fmt.Sprintf("Expected http://example.com/charts/deis-0.1.0.tgz, got %s", i.Entries["deis"][0].URLs[0]))

			// test error condition
			err := i.MustAdd(&meta.Metadata{}, "error-0.1.0.tgz", "")
			Expect(err).ShouldNot(BeNil(), "expected error adding to index")
		})
	})

	g.Describe("IndexFile.MustAdd", func() {
		var i *IndexFile

		g.BeforeEach(func() {
			i = NewIndexFile()
		})

		g.It("WHEN ", func() {
			md := &meta.Metadata{
				Name:    "scala",
				Version: "2.13.4",
				Home:    "http://www.scala-lang.org",
				Description: `Scala is a programming language for general software applications. Scala has
							full support for functional programming and a very strong static type system.
							This allows programs written in Scala to be very concise and thus smaller in
							size than other general-purpose programming languages.`,
				APIVersion: meta.APIVersionV2,
				Type:       "sdk",
			}

			err := i.MustAdd(md, "scala.sdk", "https://api.sdkman.io/2/broker/download")
			if err != nil {
				fmt.Printf("%v\n", err)
				uerr := errors.Unwrap(err)
				fmt.Printf("%v", uerr.Error())
			}

			Expect(err).Should(BeNil())
			Expect(i.Entries).Should(HaveLen(1))
		})
	})
}

func TestIndexWrite(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	var fs afero.Fs
	g.Describe("IndexFile.Write", func() {
		g.BeforeEach(func() {
			fs = afero.NewMemMapFs()
		})

		g.It("", func() {
			i := NewIndexFile()
			if err := i.MustAdd(&meta.Metadata{APIVersion: "v2", Name: "clipper", Version: "0.1.0"}, "clipper-0.1.0.tgz", "http://example.com/charts"); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			dir, err := afero.TempDir(fs, "", "helm-tmp")
			Expect(err).Should(BeNil())

			p := filepath.Join(dir, "test")
			err = i.WriteFile(fs, p, 0600)
			Expect(err).Should(BeNil())

			got, err := afero.ReadFile(fs, p)
			Expect(err).Should(BeNil())
			Expect(string(got)).Should(ContainSubstring("clipper-0.1.0.tgz"), "Index files doesn't contain expected content")
		})
	})
}

// errors.Is(uerr, meta.ValidationError(""))
