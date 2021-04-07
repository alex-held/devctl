package sdk

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/alex-held/devctl/pkg/devctlpath"
	"github.com/alex-held/devctl/pkg/plugins"
)

type testRuntimeInfoGetter struct {
	RuntimeInfo plugins.RuntimeInfo
}

func (g testRuntimeInfoGetter) Get() (ri plugins.RuntimeInfo) {
	return g.RuntimeInfo
}

type testPather struct {
	DevEnvConfigPath string
	SDKRoot          string
	DownloadPath     string
}

func (p *testPather) ConfigFilePath() string           { return p.DevEnvConfigPath }
func (p *testPather) ConfigRoot(elem ...string) string { return "" }
func (p *testPather) Config(elem ...string) string     { return "" }
func (p *testPather) Bin(elem ...string) string        { return "" }
func (p *testPather) Download(elem ...string) string {
	return path.Join(p.DownloadPath, path.Join(elem...))
}
func (p *testPather) SDK(elem ...string) string   { return path.Join(p.SDKRoot, path.Join(elem...)) }
func (p *testPather) Cache(elem ...string) string { return "" }

func TestName(t *testing.T) {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)

	mux.HandleFunc("/dl/go1.16.darwin-amd64.tar.gz", func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(w, "go1.16")
		w.WriteHeader(http.StatusOK)
	})

	uri := ts.URL
	fmt.Println(uri)

	client := http.Client{}
	resp, _ := client.Get(uri + "/dl/go1.16.darwin-amd64.tar.gz")
	assert.Equal(t, 200, resp.StatusCode)
	defer ts.Close()
}

func TestGoSDKPlugin(t *testing.T) {
	//	os.Args = append(os.Args, "-goblin.timeout=10000ms")

	g := goblin.Goblin(t)

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("devctl-sdkplugin-go", func() {
		var sut *devctlSdkpluginGo
		var fs afero.Fs
		var pathr devctlpath.Pather
		var tmpPath string
		var mux *http.ServeMux
		var teardown func()

		g.BeforeEach(func() {
			fs = afero.NewOsFs()
			mux = http.NewServeMux()
			tmpPath, _ = afero.TempDir(fs, "/tmp", "devctl")
			server := httptest.NewServer(mux)

			teardown = func() {
				server.Close()
				_ = fs.RemoveAll(tmpPath)
			}

			pathr = &testPather{
				DevEnvConfigPath: path.Join(tmpPath, "config.yaml"),
				SDKRoot:          path.Join(tmpPath, "sdks"),
				DownloadPath:     path.Join(tmpPath, "downloads"),
			}
			sut = &devctlSdkpluginGo{
				FS:         fs,
				Pather:     pathr,
				HttpClient: http.Client{},
				BaseUri:    server.URL,
				RuntimeInfoGetter: testRuntimeInfoGetter{RuntimeInfo: plugins.RuntimeInfo{
					OS:   "darwin",
					Arch: "amd64",
				}},
			}
		})

		g.After(func() {
			teardown()
		})

		g.It("Lists the installed sdks", func() {
			_ = fs.MkdirAll(pathr.SDK("go"), fileutil.PrivateDirMode)
			_ = fs.MkdirAll(pathr.SDK("go", "1.16"), fileutil.PrivateDirMode)
			_ = fs.MkdirAll(pathr.SDK("go", "1.16.2"), fileutil.PrivateDirMode)
			_ = fs.MkdirAll(pathr.SDK("go", "1.15"), fileutil.PrivateDirMode)
			_ = fs.MkdirAll(pathr.SDK("go", "1.14"), fileutil.PrivateDirMode)
			_ = fs.MkdirAll(pathr.SDK("go", "current"), fileutil.PrivateDirMode)

			expected := []string{"1.16", "1.16.2", "1.15", "1.14"}
			actual := sut.ListVersions()
			Expect(actual).Should(ContainElements(expected))
			Expect(actual).Should(HaveLen(len(expected)))
		})

		g.It("NewFunc creates a valid instance of the plugin", func() {
			actual := sut.NewFunc()
			Expect(actual.Name()).Should(Equal("devctl-sdkplugin-go"))
		})

		g.It("WHEN Download(<version>) is called => THEN the correct version gets getting downloaded", func() {
			downloadPath := pathr.Download("go", "1.16")
			artifactName := "go1.16.darwin-amd64.tar.gz"
			artifactPath := path.Join(downloadPath, artifactName)

			mux.HandleFunc("/dl/"+artifactName, func(w http.ResponseWriter, req *http.Request) {
				_, _ = io.WriteString(w, artifactName)
			})
			err := sut.Download("1.16")
			Expect(err).Should(BeNil())
			dlDirExists, _ := afero.DirExists(fs, downloadPath)
			artifactExists, _ := afero.Exists(fs, artifactPath)
			Expect(dlDirExists).Should(BeTrue())
			Expect(artifactExists).Should(BeTrue())
			artifactBytes, err := afero.ReadFile(fs, artifactPath)
			if err != nil {
				t.Errorf(err.Error())
			}
			Expect(artifactBytes).Should(Equal([]byte(artifactName)))
		})

		g.It("WHEN Extract(<version>) is called => THEN the go sdk tarball gets extracted to the correct dir", func() {
			dlDir := pathr.Download("go", "1.16")
			archiveName := "go1.16.darwin-amd64.tar.gz"
			err := fs.MkdirAll(dlDir, fileutil.PrivateDirMode)
			if err != nil {
				t.Fatal(err)
			}
			testdataFile, err := fs.Open("testdata/go1.16.darwin-amd64.tar.gz")
			if err != nil {
				t.Fatal(err)
			}
			archiveFile, err := fs.Create(path.Join(dlDir, archiveName))
			if err != nil {
				t.Fatal(err)
			}

			_, err = io.Copy(archiveFile, testdataFile)
			if err != nil {
				t.Fatal(err)
			}

			err = sut.Extract("1.16")
			Expect(err).Should(BeNil())

			sdkPath := pathr.SDK("go", "1.16")
			fis, err := afero.ReadDir(fs, sdkPath)
			if err != nil {
				t.Fatal(err)
			}
			Expect(fis).Should(HaveLen(3))
		})

		g.It("WHEN Link(<version>) is called => THEN sdks/go/<version> is symlinked to sdks/go/current", func() {
			versionSDKPath := pathr.SDK("go", "1.16")
			_ = fs.MkdirAll(versionSDKPath, fileutil.PrivateDirMode)
			f, _ := fs.Create(path.Join(versionSDKPath, "1.16"))
			_, _ = f.WriteString("1.16")
			_ = f.Close()
			currentSDKPath := pathr.SDK("go", "current")

			err := sut.Link("1.16")
			Expect(err).Should(BeNil())

			currentSDKDirExists, _ := afero.DirExists(fs, currentSDKPath)
			bytes, err := afero.ReadFile(fs, path.Join(currentSDKPath, "1.16"))
			if err != nil {
				t.Fatal(err)
			}
			Expect(currentSDKDirExists).Should(BeTrue())
			Expect(string(bytes)).Should(Equal("1.16"))
		})

		g.It("WHEN Install(<version>) is called => THEN the correct version gets linked to current", func() {
			artifactName := "go1.16.darwin-amd64.tar.gz"
			mux.HandleFunc("/dl/"+artifactName, func(w http.ResponseWriter, req *http.Request) {
				bytes, _ := afero.ReadFile(fs, "testdata/go1.16.darwin-amd64.tar.gz")
				_, _ = w.Write(bytes)
			})
			err := sut.InstallE("1.16")
			Expect(err).Should(BeNil())

			currentSDKPath := pathr.SDK("go", "current")
			currentSDKDirExists, _ := afero.DirExists(fs, currentSDKPath)
			Expect(currentSDKDirExists).Should(BeTrue())
		})
	})
}
