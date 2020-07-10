package kernel

import (
    "fmt"
    . "io"
    "os"
    "time"
    
    "github.com/gosuri/uilive"
)

type NamedWriter interface {
    Writer
    Name() string
}

type namedWriter struct {
    name   string
    Writer Writer
}

func (w namedWriter) Write(bytes []byte) (n int, err error) { return w.Writer.Write(bytes) }
func (w namedWriter) Name() string                          { return w.name }

func NewNamedWriter(name string, writer Writer) NamedWriter {
    return namedWriter{
        name:   name,
        Writer: writer,
    }
}

type Installer interface {
    Install()
}

type BrewManifest struct {
    Taps  []string `yaml:"taps"`
    Brews []string `yaml:"brews"`
    Casks []string `yaml:"casks"`
}

type Option func() (string, func(interface{}))

type WriterOption struct {
    name string
    opt  func(config *WriterConfig)
}

func (c *WriterConfig) ConfigName() string { return "WriterConfig" }

func (b BrewConfig) ConfigName() string { return "BrewConfig" }

type brewInstaller struct {
    config BrewConfig
    writer Writer
}

func CreateWriter(opts ...WriterOptions) Writer {
    config := WriterConfig{
        UseStdOut: true,
        Writers:   []Writer{},
    }
    
    for _, opt := range opts {
        _ = opt(&config)
    }
    
    if config.UseStdOut {
        config.Writers = append(config.Writers, NewNamedWriter("std:out", os.Stdout))
    }
    
    return MultiWriter(config.Writers...)
}

func NewBrewInstaller(options ...BrewOptions) Installer {
    config := BrewConfig{}
    config.Manifest = BrewManifest{
        Taps:  []string{},
        Brews: []string{},
        Casks: []string{},
    }
    
    for _, option := range options {
        _ = option(&config)
    }
    
    installer := brewInstaller{config: config}
    installer.writer = CreateWriter(installer.config.WriterOptions...)
    
    return installer
}

func (installer brewInstaller) Install() {
    writer := uilive.New()
    writer.Out = installer.writer
    
    writer.Start()
    
    m := installer.config.Manifest
    
    tapCount := len(m.Taps)
    brewCount := len(m.Brews)
    caskCount := len(m.Casks)
    
    for i, tap := range m.Taps {
        _, _ = fmt.Fprintf(writer, "Adding Tap '%s' (%d/%d)\n", tap, i, tapCount)
        time.Sleep(time.Millisecond * 200)
    }
    
    for i, brew := range m.Brews {
        _, _ = fmt.Fprintf(writer, "Installing Brew '%s' (%d/%d)\n", brew, i, brewCount)
        time.Sleep(time.Millisecond * 200)
    }
    
    for i, cask := range m.Casks {
        _, _ = fmt.Fprintf(writer, "Adding Cask '%s' (%d/%d)\n", cask, i, caskCount)
        time.Sleep(time.Millisecond * 500)
    }
    
    _, _ = fmt.Fprintln(writer, "Finished: Downloaded 100GB")
    
    writer.Stop()
}
