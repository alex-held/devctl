package sdkman

import (
	"net/http"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

type optionsTestCase struct {
	option ClientOption
	want   ClientConfig
}

func TestURLOptions(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("ClientOption", func() {
		c := &ClientConfig{}

		g.JustBeforeEach(func() {
			c = &ClientConfig{}
		})

		g.Describe("URLOptions", func() {
			g.It("WHEN baseUrl not set => THEN set baseURL", func() {
				tc := optionsTestCase{
					option: URLOptions("https://base.url/v1"),
					want: ClientConfig{
						httpClient: nil,
						fs:         nil,
						baseURL:    "https://base.url/v1",
					},
				}
				actual := tc.option(c)
				Expect(*actual).To(Equal(tc.want))
			})

			g.It("WHEN baseUrl already set => THEN overwrite BaseUrl", func() {
				tc := optionsTestCase{
					option: URLOptions("new-value"),
					want: ClientConfig{
						httpClient: nil,
						fs:         nil,
						baseURL:    "new-value",
					},
				}
				c.baseURL = "overwrite-me"
				actual := tc.option(c)
				Expect(*actual).To(Equal(tc.want))
			})
		})

		g.Describe("HTTPClientOption", func() {
			g.It("WHEN http.Client is nil => THEN use http.DefaultClient", func() {
				tc := optionsTestCase{
					option: HTTPClientOption(nil),
					want: ClientConfig{
						httpClient: http.DefaultClient,
					},
				}

				actual := tc.option(c)
				Expect(*actual).To(Equal(tc.want))
			})

			g.It("WHEN http.Client already set => THEN overwrite it", func() {
				client := &http.Client{}
				tc := optionsTestCase{
					option: HTTPClientOption(client),
					want: ClientConfig{
						httpClient: client,
					},
				}
				c.httpClient = http.DefaultClient
				actual := tc.option(c)
				Expect(actual.httpClient).Should(Equal(client))
				Expect(*actual).To(Equal(tc.want))
			})

			g.It("WHEN http.Client not set => THEN set it", func() {
				client := &http.Client{}
				tc := optionsTestCase{
					option: HTTPClientOption(client),
					want: ClientConfig{
						httpClient: client,
					},
				}
				actual := tc.option(c)
				Expect(actual.httpClient).Should(Equal(client))
				Expect(*actual).To(Equal(tc.want))
			})
		})

		g.Describe("FileSystemOption", func() {
			g.It("WHEN afero.Fs is nil => THEN use afero.OsFs", func() {
				tc := optionsTestCase{
					option: FileSystemOption(nil),
					want: ClientConfig{
						fs: afero.NewOsFs(),
					},
				}

				actual := tc.option(c)
				Expect(*actual).To(Equal(tc.want))
				Expect(actual.fs).To(Equal(tc.want.fs))
				Expect(actual.fs.Name()).To(Equal(tc.want.fs.Name()))
			})

			g.It("WHEN afero.Fs already set => THEN overwrite it", func() {
				fs := afero.NewMemMapFs()
				tc := optionsTestCase{
					option: FileSystemOption(fs),
					want: ClientConfig{
						fs: fs,
					},
				}
				c.fs = afero.NewOsFs()
				actual := tc.option(c)
				Expect(actual.fs).Should(Equal(fs))
				Expect(actual.fs.Name()).Should(Equal(fs.Name()))
				Expect(*actual).To(Equal(tc.want))
			})

			g.It("WHEN afero.Fs not set => THEN set it", func() {
				fs := afero.NewMemMapFs()
				tc := optionsTestCase{
					option: FileSystemOption(fs),
					want: ClientConfig{
						fs: fs,
					},
				}
				actual := tc.option(c)
				Expect(actual.fs).Should(Equal(fs))
				Expect(actual.fs.Name()).Should(Equal(fs.Name()))
				Expect(*actual).To(Equal(tc.want))
			})
		})
	})
}
