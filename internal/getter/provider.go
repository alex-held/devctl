package getter

import (
	"bytes"
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/alex-held/devctl/internal/cli"
	"github.com/alex-held/devctl/internal/sdkman"
)

// Getter is an interface to support GET to the specified URL.
type Getter interface {
	// Get file content by url string
	Get(url string, options ...Option) (*bytes.Buffer, error)
}

// options are generic parameters to be provided to the getter during instantiation.
//
// Getters may or may not ignore these parameters as they are passed in.
type options struct {
	url                   string
	unTar                 bool
	insecureSkipVerifyTLS bool
	userAgent             string
	version               string
	timeout               time.Duration
	ctx                   context.Context
	registry              *sdkman.RegistryService
}

// Option allows specifying various settings configurable by the user for overriding the defaults
// used when performing Get operations with the Getter.
type Option func(*options)

// WithURL informs the getter the server name that will be used when fetching objects. Used in conjunction with
// WithTLSClientConfig to set the TLSClientConfig's server name.
func WithURL(url string) Option {
	return func(opts *options) {
		opts.url = url
	}
}

func WithTagName(tagname string) Option {
	return func(opts *options) {
		opts.version = tagname
	}
}

func WithRegistryClient(client *sdkman.RegistryService) Option {
	return func(opts *options) {
		opts.registry = client
	}
}

func WithUntar() Option {
	return func(opts *options) {
		opts.unTar = true
	}
}

// Constructor is the function for every getter which creates a specific instance
// according to the configuration
type Constructor func(options ...Option) (Getter, error)

type Providers []Provider

type Provider struct {
	Schemes []string
	New     Constructor
}

// Provides returns true if the given scheme is supported by this Provider.
func (p Provider) Provides(scheme string) bool {
	for _, i := range p.Schemes {
		if i == scheme {
			return true
		}
	}
	return false
}

// ByScheme returns a Provider that handles the given scheme.
//
// If no provider handles this scheme, this will return an error.
func (p Providers) ByScheme(scheme string) (Getter, error) {
	for _, pp := range p {
		if pp.Provides(scheme) {
			return pp.New()
		}
	}
	return nil, errors.Errorf("scheme %q not supported", scheme)
}

var httpProvider = Provider{
	Schemes: []string{"http", "https"},
	New:     NewHTTPGetter,
}

// All finds all of the registered getters as a list of Provider instances.
// Currently, the built-in getters and the discovered plugins with downloader
// notations are collected.
func All(settings *cli.Env) Providers {
	result := Providers{httpProvider}
	return result
}
