package repo

import (
	"net/url"

	"github.com/pkg/errors"

	"github.com/alex-held/devctl/internal/devctlpath"
	"github.com/alex-held/devctl/internal/getter"
)

// Entry
type Entry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// SDKRepository represents a sdk repository
// For example sdkman.io
type SDKRepository struct {
	Config     *Entry
	ChartPaths []string
	Client     getter.Getter
	CachePath  string
}

func NewSDKRepo(cfg *Entry, getters getter.Providers) (*SDKRepository, error) {
	u, err := url.Parse(cfg.URL)

	if err != nil {
		return nil, errors.Errorf("invalid chart URL format: %s", cfg.URL)
	}

	client, err := getters.ByScheme(u.Scheme)
	if err != nil {
		return nil, errors.Errorf("could not find protocol handler for: %s", u.Scheme)
	}

	return &SDKRepository{
		Config:    cfg,
		Client:    client,
		CachePath: devctlpath.CachePath("repository"),
	}, nil
}
