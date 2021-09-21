package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-github/v35/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/alex-held/devctl/pkg/devctlpath"
)

func TestPlugins(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Plugin Suite")
}

var _ = Describe("Plugin", func() {
	Describe("Store", func() {
		tmpRoot := "/tmp/devctl"
		pather := devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			return tmpRoot
		}))
		manifest := pather.ConfigRoot("plugins.yaml")
		fs := afero.NewMemMapFs()
		storeP := &store{
			Pather: pather,
			Fs:     fs,
		}
		sut := Store(storeP)

		Context("no plugin manifest exists", func() {
			plugins, err := sut.List(SDK)
			It("returns no error", func() {
				fmt.Printf("Error:\t%+v\n", err)
				Expect(err).Should(Succeed())
			})
			It("returns empty list", func() { Expect(plugins).Should(BeEmpty()) })
		})

		Context("empty plugin manifest exists", func() {
			_, _ = fs.Create(manifest)
			plugins, err := sut.List(SDK)
			It("returns no error", func() { Expect(err).Should(Succeed()) })
			It("returns empty list", func() { Expect(plugins).Should(BeEmpty()) })
		})

		Context("plugin manifest with empty categories exists", func() {
			data, _ := afero.ReadFile(afero.NewOsFs(), "testdata/plugins_empty.yaml")
			_ = afero.WriteFile(fs, manifest, data, 0777)
			plugins, err := sut.List(SDK)

			It("returns no error", func() { Expect(err).Should(Succeed()) })
			It("returns empty list", func() { Expect(plugins).Should(BeEmpty()) })
		})

		Context("plugin manifest with content exists", func() {
			data, _ := afero.ReadFile(afero.NewOsFs(), "testdata/plugins_1.yaml")
			_ = afero.WriteFile(fs, manifest, data, 0777)
			plugins, err := sut.List(SDK)

			It("returns no error", func() { Expect(err).Should(Succeed()) })
			It("list contains sdk plugin '<root>/plugins/go.so'", func() { Expect(plugins).Should(ContainElement("/tmp/devctl/plugins/go.so")) })
		})

		When("no plugin manifest exists", func() {
			_ = fs.RemoveAll(manifest)
			registerErr := sut.Register(SDK, "scala")
			fi, statErr := fs.Stat(manifest)

			It("should return no error", func() { Expect(registerErr).Should(Succeed()) })
			It("creates the plugin manifest", func() {
				Expect(fi).ShouldNot(BeNil())
				Expect(statErr).Should(BeNil())
			})
		})

		When("empty plugin manifest exists", func() {
			err := sut.Register(SDK, "scala")
			file, err := sut.(*store).load()
			sdkPlugins, err := sut.(*store).List(SDK)

			It("not return error", func() { Expect(err).Should(Succeed()) })
			It("creates the SDK category", func() { Expect(file.SDK).ShouldNot(BeNil()) })
			It("creates the corresponding category with the registered plugin", func() {
				Expect(sdkPlugins).Should(ContainElement(ContainSubstring("scala")))
			})
		})

		When("plugin manifest with empty categories exists", func() {
			data, _ := afero.ReadFile(afero.NewOsFs(), "testdata/plugins_empty.yaml")
			_ = afero.WriteFile(fs, manifest, data, 0777)
			err := sut.Register(SDK, "scala")
			f, err := sut.(*store).load()

			It("not return error", func() { Expect(err).Should(Succeed()) })
			It("has one SDK plugin registered", func() { Expect(f.SDK).Should(HaveLen(1)) })
			It("appends the registered plugin to the corresponding category", func() {
				Expect(f.SDK["scala"]).Should(ContainSubstring("scala"))
			})
		})

		When("plugin manifest with content exists", func() {
			data, _ := afero.ReadFile(afero.NewOsFs(), "testdata/plugins_1.yaml")
			_ = afero.WriteFile(fs, manifest, data, 0777)
			plugins, err := sut.List(SDK)

			It("returns no error", func() { Expect(err).Should(Succeed()) })
			It("plugins have length 1", func() { Expect(plugins).Should(HaveLen(1)) })
			It("plugins has an entry with key 'go'", func() { Expect(plugins["go"]).ShouldNot(BeNil()) })
			It("plugin value 'go' contains plugin path", func() {
				Expect(plugins["go"]).Should(ContainSubstring(pather.ConfigRoot("plugins", "go.so")))
			})
		})
	})

	Describe("Kind", func() {
		Context("NewKind", func() {
			Context("Kind Name/Value is a valid value", func() {
				kind, err := NewKind(SDK.Name())
				It("should not return an error", func() { Expect(err).Should(Succeed()) })
				It("kind has correct value", func() { Expect(kind.String()).Should(Equal("SDK")) })
				It("kine has correct Name", func() { Expect(kind.Name()).Should(Equal("SDK")) })
				It("kine has correct Description", func() {
					Expect(kind.Description()).Should(Equal("installs updates and manages different sdks on your system"))
				})
			})

			Context("Kind Name/Value has an invalid value", func() {
				kind, err := NewKind("foo")
				It("should return an error", func() { Expect(err).Should(MatchError(fmt.Sprintf("'%s' is not a valid value for type", "foo"))) })
				It("kind has empty value", func() { Expect(kind.String()).Should(Equal("")) })
				It("kine has empty Name", func() { Expect(kind.Name()).Should(Equal("")) })
				It("kine has empty Description", func() { Expect(kind.Description()).Should(Equal("")) })
			})
		})
	})

	Describe("Query", func() {
		Context("constructing query parameter", func() {
			Context("initial default value", func() {
				actual := NewQuery()
				It("contains devctl-plugin constrain", func() { Expect(actual).Should(HavePrefix("q=topic:devctl-plugin")) })
			})
			Context("kind is specified", func() {
				actual := NewQuery(WithKind(SDK))
				It("query contains 'topic:devctl-[kind]-plugin' constraint", func() { Expect(actual).Should(Equal("q=topic:devctl-plugin+topic:devctl-sdk-plugin")) })
			})
			Context("plugin name is specified", func() {
				actual := NewQuery(WithName("go"))
				It("contains name constraint", func() { Expect(actual).Should(Equal("q=topic:devctl-plugin+go+in:name")) })
			})
			Context("all constraints are specified", func() {
				actual := NewQuery(WithKind(SDK), WithName("go"))
				It("query contains 'topic:devctl-[kind]-plugin' constraint", func() { Expect(actual).Should(ContainSubstring("+topic:devctl-sdk-plugin")) })
				It("contains name constraint", func() { Expect(actual).Should(ContainSubstring("+go+in:name")) })
			})
		})
	})

	Describe("Client", func() {
		Context("default client", func() {
			Context("github found no results", func() {
				var model SearchResponseModel
				var sut Client
				var mux *http.ServeMux
				var teardown func()
				var err error
				var plugins PluginSearchResults

				BeforeEach(func() {
					sut, mux, _, teardown = setup()
					model = SearchResponseModel{}

					mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
						json := model.ToJsonResponse()
						io.WriteString(w, json)
					})

					plugins, err = sut.Search(context.Background(), "topic:devctl-plugin")
				})

				It("should not return an error", func() { Expect(err).Should(Succeed()) })
				It("should return empty list of PluginResults", func() { Expect(plugins).Should(BeEmpty()) })

				AfterEach(func() {
					teardown()
				})
			})
			Context("github one sdk plugin as result", func() {

				var sut Client
				var mux *http.ServeMux
				var teardown func()
				var err error
				var plugins PluginSearchResults

				BeforeEach(func() {
					model := SearchResponseModel{
						{
							user:   "alex-held",
							name:   "devctl-sdkplugin-go",
							topics: []string{"devctl-plugin", "devctl-sdk-plugin", "devctl"},
						},
					}
					sut, mux, _, teardown = setup()
					mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
						json := model.ToJsonResponse()
						io.WriteString(w, json)
					})
					plugins, err = sut.Search(context.Background(), "topic:devctl-plugin")
				})

				It("should not return an error", func() { Expect(err).Should(Succeed()) })
				It("should return list of one PluginResult", func() { Expect(plugins).Should(HaveLen(1)) })
				It("should contain sdk plugin", func() {
					Expect(plugins[0]).Should(Equal(PluginSearchResult{
						Name:           "devctl-sdkplugin-go",
						Uri:            "https://api.github.com/repos/alex-held/devctl-sdkplugin-go",
						RepositoryName: "devctl-sdkplugin-go",
						Topics:         []string{"devctl-plugin", "devctl-sdk-plugin", "devctl"},
						Kind:           SDK,
					}))
				})

				AfterEach(func() {
					teardown()
				})
			})

		})
	})
})

type pluginResponseModel struct {
	name   string
	topics []string
	user   string
}

type SearchResponseModel []pluginResponseModel

func (srm SearchResponseModel) ToJsonResponse() string {
	sb := strings.Builder{}
	itemCount := len(srm)

	if itemCount == 0 {
		sb.WriteString(`{
			"total_count": 0,
			"incomplete_results": true,
			"items": []
		}`)
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("{\n"))
	sb.WriteString(fmt.Sprintf("\t\"total_count\": %d,\n", itemCount))
	sb.WriteString(fmt.Sprintf("\t\"incomplete_results\": false,\n"))
	sb.WriteString(fmt.Sprintf("\t\"items\": [\n"))

	for i, m := range srm {
		item := "\t" + jsonItem(m, i, itemCount)
		sb.WriteString(item)
	}

	sb.WriteString("\t]\n")
	sb.WriteString("}\n")
	jsonResp := sb.String()
	return jsonResp
}

func jsonItem(m pluginResponseModel, current, itemCount int) string {
	endLine := ""
	if current > itemCount {
		endLine = ",\n"
	} else {
		endLine = "\n"
	}

	topicsJsonBytes, _ := json.Marshal(&m.topics)
	topics := string(topicsJsonBytes)

	master := `{
            "id": 363272734,
            "node_id": "MDEwOlJlcG9zaXRvcnkzNjMyNzI3MzQ=",
            "name": "#{repo_name}#",
            "full_name": "#{user_name}#/#{repo_name}#",
            "private": false,
			"topics": #{topics}#,
            "owner": {
                "login": "#{user_name}#",
                "id": 50153092,
                "node_id": "MDQ6VXNlcjUwMTUzMDky",
                "avatar_url": "https://avatars.githubusercontent.com/u/50153092?v=4",
                "gravatar_id": "",
                "url": "https://api.github.com/users/#{user_name}#",
                "html_url": "https://github.com/#{user_name}#",
                "followers_url": "https://api.github.com/users/#{user_name}#/followers",
                "following_url": "https://api.github.com/users/#{user_name}#/following{/other_user}",
                "gists_url": "https://api.github.com/users/#{user_name}#/gists{/gist_id}",
                "starred_url": "https://api.github.com/users/#{user_name}#/starred{/owner}{/repo}",
                "subscriptions_url": "https://api.github.com/users/#{user_name}#/subscriptions",
                "organizations_url": "https://api.github.com/users/#{user_name}#/orgs",
                "repos_url": "https://api.github.com/users/#{user_name}#/repos",
                "events_url": "https://api.github.com/users/#{user_name}#/events{/privacy}",
                "received_events_url": "https://api.github.com/users/#{user_name}#/received_events",
                "type": "User",
                "site_admin": false
            },
            "html_url": "https://github.com/#{user_name}#/#{repo_name}#",
            "description": null,
            "fork": false,
            "url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#",
            "forks_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/forks",
            "keys_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/keys{/key_id}",
            "collaborators_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/collaborators{/collaborator}",
            "teams_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/teams",
            "hooks_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/hooks",
            "issue_events_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/issues/events{/number}",
            "events_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/events",
            "assignees_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/assignees{/user}",
            "branches_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/branches{/branch}",
            "tags_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/tags",
            "blobs_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/git/blobs{/sha}",
            "git_tags_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/git/tags{/sha}",
            "git_refs_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/git/refs{/sha}",
            "trees_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/git/trees{/sha}",
            "statuses_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/statuses/{sha}",
            "languages_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/languages",
            "stargazers_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/stargazers",
            "contributors_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/contributors",
            "subscribers_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/subscribers",
            "subscription_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/subscription",
            "commits_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/commits{/sha}",
            "git_commits_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/git/commits{/sha}",
            "comments_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/comments{/number}",
            "issue_comment_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/issues/comments{/number}",
            "contents_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/contents/{+path}",
            "compare_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/compare/{base}...{head}",
            "merges_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/merges",
            "archive_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/{archive_format}{/ref}",
            "downloads_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/downloads",
            "issues_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/issues{/number}",
            "pulls_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/pulls{/number}",
            "milestones_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/milestones{/number}",
            "notifications_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/notifications{?since,all,participating}",
            "labels_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/labels{/name}",
            "releases_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/releases{/id}",
            "deployments_url": "https://api.github.com/repos/#{user_name}#/#{repo_name}#/deployments",
            "created_at": "2021-04-30T22:24:57Z",
            "updated_at": "2021-04-30T22:44:16Z",
            "pushed_at": "2021-04-30T22:25:28Z",
            "git_url": "git://github.com/#{user_name}#/#{repo_name}#.git",
            "ssh_url": "git@github.com:#{user_name}#/#{repo_name}#.git",
            "clone_url": "https://github.com/#{user_name}#/#{repo_name}#.git",
            "svn_url": "https://github.com/#{user_name}#/#{repo_name}#",
            "homepage": "",
            "size": 0,
            "stargazers_count": 0,
            "watchers_count": 0,
            "language": null,
            "has_issues": true,
            "has_projects": true,
            "has_downloads": true,
            "has_wiki": true,
            "has_pages": false,
            "forks_count": 0,
            "mirror_url": null,
            "archived": false,
            "disabled": false,
            "open_issues_count": 0,
            "license": null,
            "forks": 0,
            "open_issues": 0,
            "watchers": 0,
            "default_branch": "master",
            "score": 1.0
        }`

	replaced := strings.ReplaceAll(master, "#{repo_name}#", m.name)
	//	fmt.Printf("replaced repo_name '%s'\t[1/3]:\t%s\n", m.name, replaced)
	replaced = strings.ReplaceAll(replaced, "#{user_name}#", m.user)
	//	fmt.Printf("replaced user_name '%s'\t[2/3]:\t%s\n", m.user, replaced)
	replaced = strings.ReplaceAll(replaced, "#{topics}#", topics)
	//	fmt.Printf("replaced topics '%s'\t[3/3]:\t%s\n", topics, replaced)

	return replaced + endLine
}

const (
// baseURLPath is a non-empty Client.BaseURL path to use during tests,
// to ensure relative URLs are used for all endpoints. See issue #752.
// baseURLPath = "/api-v3"
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (c Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	//	apiHandler := http.NewServeMux()
	// apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	baseURL, _ := url.Parse(server.URL + "/")
	c = NewClient(baseURL, nil)
	return c, mux, server.URL, server.Close
}

func TestRepo(t *testing.T) {
	username := "alex-held"
	name := "devctl-sdkplugin-go"
	topics := []string{"devctl-plugin", "go", "plugin", "devctl-sdk-plugin", "sdk"}
	var repo = github.Repository{
		Owner: &github.User{
			Name: &username,
		},
		Name:   &name,
		Topics: topics,
	}
	kind, err := ParseKindFromRepo(repo)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, SDK, kind)
}
