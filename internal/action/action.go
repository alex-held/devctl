package action

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/logging"

	"github.com/alex-held/devctl/internal/devctlpath"
	"github.com/alex-held/devctl/internal/sdkman"
)

// Options contains the configuration options for Actions
type Options struct {
	Fs     afero.Fs
	Pather devctlpath.Pather
	Client *sdkman.Client
	Logger *logging.Logger
}

type Option func(options *Options) *Options

func WithLogger(l *logging.Logger) Option {
	return func(o *Options) *Options {
		o.Logger = l
		return o
	}
}
func WithFs(fs afero.Fs) Option {
	return func(o *Options) *Options {
		o.Fs = fs
		return o
	}
}

func WithSdkmanClient(client *sdkman.Client) Option {
	return func(o *Options) *Options {
		o.Client = client
		return o
	}
}

func WithPather(pather devctlpath.Pather) Option {
	return func(o *Options) *Options {
		o.Pather = pather
		return o
	}
}

var defaults = []Option{
	WithFs(afero.NewOsFs()),
	WithSdkmanClient(sdkman.NewSdkManClient()),
	WithPather(devctlpath.NewPather()),
	WithLogger(logging.NewLogger()),
}

type action struct {
	*Actions
}

type Actions struct {
	*action
	Options  *Options
	Install  *Install
	Download *Download
	Config   *Config
	Symlink  *Symlink
}

func NewActions(opts ...Option) *Actions {
	options := &Options{}

	for _, opt := range defaults {
		options = opt(options)
	}

	for _, opt := range opts {
		options = opt(options)
	}

	actions := &Actions{
		Options: options,
	}

	common := &action{
		Actions: actions,
	}

	actions.action = common
	actions.Download = (*Download)(actions.action)
	actions.Install = (*Install)(actions.action)
	actions.Config = (*Config)(actions.action)
	actions.Symlink = (*Symlink)(actions.action)

	return actions
}

func saveArchive(fs afero.Fs, buf bytes.Buffer, path string) (file afero.File, err error) {
	exist, err := afero.Exists(fs, path)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot check whether downloaded archive already exists; archive=%s\n", path)
	}
	if exist {
		return nil, nil
	}
	archive, err := fs.Create(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create file; path=%s", path)
	}
	n, err := io.Copy(archive, &buf)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to copy http content into archive file; written %d bytes\n", n)
	}
	return archive, nil
}
