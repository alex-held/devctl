package plugins

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/go-github/v35/github"
	"github.com/pkg/errors"
)

type PluginSearchResults []PluginSearchResult
type PluginSearchResult struct {
	Name           string
	Uri            string
	RepositoryName string
	Topics         []string
	Kind           Kind
}

var (
	ErrNoValidTopics = errors.New("topics did not contain a valid topic matching 'devctl-<PLUGIN_KIND>-plugin'")
	ErrInvalidTopic  = errors.New("failed to parse string as a plugin kind that is not known")
)

type Client interface {
	SearchOpts(ctx context.Context, option ...QueryOption) (res PluginSearchResults, err error)
	Search(ctx context.Context, query string) (res PluginSearchResults, err error)
}

type client struct {
	GithubClient *github.Client
	BaseUri      *url.URL
	OutWriter    io.Writer
}

func NewClient(baseURL *url.URL, httpClient *http.Client) Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	c := &client{
		BaseUri:      baseURL,
		GithubClient: github.NewClient(httpClient),
		OutWriter:    &bytes.Buffer{},
	}
	c.GithubClient.BaseURL = baseURL

	return c
}

type Query struct {
	name string
	kind Kind
}

type QueryOption func(query *Query) *Query

func WithKind(kind Kind) QueryOption {
	return func(query *Query) *Query {
		query.kind = kind
		return query
	}
}

func WithName(name string) QueryOption {
	return func(query *Query) *Query {
		query.name = name
		return query
	}
}

func NewQuery(opts ...QueryOption) string {
	q := &Query{}
	for _, opt := range opts {
		q = opt(q)
	}
	query := q.String()
	return query
}
func (p *Query) String() string {
	query := "q=topic:devctl-plugin"
	if p == nil {
		return query
	}
	if IsValidKind(p.kind) {
		query += fmt.Sprintf("+topic:devctl-%s-plugin", strings.ToLower(p.kind.String()))
	}
	if p.name != "" {
		query += fmt.Sprintf("+%s+in:name", strings.ToLower(p.name))
	}
	return query
}

func (c *client) QueryPlugins(ctx context.Context, query Query) (res PluginSearchResults, err error) {
	return c.Search(ctx, query.String())
}
func (c *client) SearchOpts(ctx context.Context, opts ...QueryOption) (res PluginSearchResults, err error) {
	return c.Search(ctx, NewQuery(opts...))
}
func (c *client) Search(ctx context.Context, query string) (res PluginSearchResults, err error) {
	if c.GithubClient == nil {
		c.GithubClient = github.NewClient(nil)
	}
	repos, resp, err := c.GithubClient.Search.Repositories(ctx, query, &github.SearchOptions{
		Sort: "stars",
	})

	if err != nil || resp.StatusCode != http.StatusOK {
		return res, err
	}
	for _, repo := range repos.Repositories {
		p := PluginSearchResult{
			Name:           repo.GetName(),
			Uri:            repo.GetURL(),
			RepositoryName: repo.GetName(),
			Topics:         repo.Topics,
		}
		p.Kind, err = ParseKindFromRepo(*repo)
		res = append(res, p)
	}
	return res, err
}

func getKindFromTopics(topics []string) (kind Kind, err error) {
	for _, topic := range topics {
		regex := regexp.MustCompile("devctl-(?P<kind>.+?)-plugin")
		if !regex.MatchString(topic) {
			continue
		}
		kindStr := regex.ReplaceAllString(topic, "$1")
		kind, err = NewKind(strings.ToUpper(kindStr))
		if err != nil {
			err = errors.Wrapf(ErrInvalidTopic, "%+v", err)
			return kind, err
		}
		return kind, err
	}

	return kind, ErrNoValidTopics
}

func ParseKindFromRepo(repo github.Repository) (k Kind, err error) {
	return getKindFromTopics(repo.Topics)
}
