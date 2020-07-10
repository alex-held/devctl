package kernel

import (
    . "io"
)


// BrewOptions
type BrewOptions func(*BrewConfig) (err error)
type BrewConfig struct {
    Manifest   BrewManifest
    DryRun     bool
    WriterOptions []WriterOptions
}

func DryRun() BrewOptions {
    return func(options *BrewConfig) (err error) {
        options.DryRun = true
        return nil
    }
}

func WithManifest(m BrewManifest) BrewOptions {
    return func(c *BrewConfig) (err error) {
        c.Manifest = m
        return nil
    }
}

func UseWriterOptions(opts ...WriterOptions) BrewOptions {
    return func(c *BrewConfig) (err error) {
        c.WriterOptions = append(c.WriterOptions, opts...)
        return nil
    }
}



// WriterOptions
type WriterOptions func(*WriterConfig) (err error)
type WriterConfig struct {
    UseStdOut bool
    Writers   []Writer
}

func WithTestWriter(writer Writer) WriterOptions {
    return WithNamedWriter("test-output", writer)
}

func UseStdout() WriterOptions {
    return func(c *WriterConfig) (err error) {
        c.UseStdOut = true
        return nil
    }
}

func Quiet() WriterOptions {
    return func(c *WriterConfig) (err error) {
        c.UseStdOut = false
        return nil
    }
}


func WithNamedWriter(name string, writer Writer) WriterOptions {
    return func(c *WriterConfig) (err error) {
        c.Writers = append(c.Writers, NewNamedWriter(name, writer))
        return nil
    }
}
