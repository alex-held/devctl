package util

import (
	"io"
	"os"
	"sync"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
	"github.com/alex-held/devctl-kit/pkg/log"
	"github.com/alex-held/devctl-kit/pkg/system"

	"github.com/alex-held/devctl/pkg/cli/options"
	"github.com/alex-held/devctl/pkg/validation"
)

type factory struct {
	runtimeInfoGetter system.RuntimeInfoGetter
	pather            devctlpath.Pather
	logger            log.Logger

	// Caches OpenAPI document and parsed resources
	//	openAPIParser *openapi.CachedOpenAPIParser
	//	openAPIGetter *openapi.CachedOpenAPIGetter
	//	parser        sync.Once
	getter  sync.Once
	streams options.IOStreams
}

func (f *factory) RuntimeInfo() system.RuntimeInfo {
	return f.runtimeInfoGetter.Get()
}

func (f *factory) NewBuilder() *Builder {
	return &Builder{}
}

func (f *factory) Logger() log.Logger {
	return f.logger
}

func (f *factory) Pather() devctlpath.Pather {
	return f.pather
}

func (f *factory) Streams() options.IOStreams {
	return f.streams
}

func (f *factory) Validator(validate bool) (validation.Schema, error) {
	return validation.NullSchema{}, nil
}

// Factory provides abstractions that allow the Devctl command to be extended across multiple types
// of resources and different API sets.
type Factory interface {
	RuntimeInfo() system.RuntimeInfo

	// NewBuilder returns an object that assists in loading objects from both disk and the server
	// and which implements the common patterns for CLI interactions with generic resources.
	NewBuilder() *Builder

	Logger() log.Logger

	Pather() devctlpath.Pather

	Streams() options.IOStreams

	// Returns a schema that can validate objects stored on disk.
	Validator(validate bool) (validation.Schema, error)

	// OpenAPISchema returns the parsed openapi schema definition
	//	OpenAPISchema() (openapi.Resources, error)
	// OpenAPIGetter returns a getter for the openapi schema document
	//	OpenAPIGetter() discovery.OpenAPISchemaInterface
}

type FactoryConfig struct {
	Pather            devctlpath.Pather
	LoggerConfig      *log.Config
	Streams           *options.IOStreams
	RuntimeInfoGetter system.RuntimeInfoGetter
}

type FactoryOption func(*FactoryConfig) *FactoryConfig

func WithIO(in io.Reader, out, err io.Writer) FactoryOption {
	return func(c *FactoryConfig) *FactoryConfig {
		c.Streams = &options.IOStreams{
			In:     in,
			Out:    out,
			ErrOut: err,
		}
		return c
	}
}

func NewFactory(opts ...FactoryOption) Factory {
	cfg := &FactoryConfig{
		Pather:       devctlpath.DefaultPather(),
		LoggerConfig: &log.DefaultConfig,
	}
	defaults := []FactoryOption{
		WithIO(os.Stdin, os.Stdout, os.Stdout),
	}

	for _, opt := range append(defaults, opts...) {
		opt(cfg)
	}

	return &factory{
		runtimeInfoGetter: cfg.RuntimeInfoGetter,
		pather:            cfg.Pather,
		logger:            log.New(cfg.LoggerConfig),
		getter:            sync.Once{},
		streams:           *cfg.Streams,
	}
}
