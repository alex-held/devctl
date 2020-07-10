package kernel

import (
    "bytes"
    "fmt"
    "reflect"
    "strings"
    "testing"
    
    "github.com/stretchr/testify/require"
)

func TestOptionsInterface(t *testing.T) {
    
    var option = func() (string, func(interface{})) {
        return "Dry-Run", func(i interface{}) {
            switch c := i.(type) {
            case BrewConfig:
                c.DryRun = true
            case *BrewConfig:
                c.DryRun = true
            default:
                panic(fmt.Errorf("Cannot cast config of type %s\n", reflect.TypeOf(i).Name()))
            }
        }
    }
    
    OptionTest(option)
}

func TestProgress(t *testing.T) {
  
    out := &bytes.Buffer{}
    
    manifest := BrewManifest{
        Taps: []string{
            "homebrew/bundle",
            "homebrew/cask-versions",
        },
        Brews: []string{
            "curl",
            "wget",
            "zsh",
        },
    }
    i := NewBrewInstaller(
        DryRun(),
        WithManifest(manifest),
        UseWriterOptions(
            WithTestWriter(out),
            UseStdout(),
        ),
    )
    
    i.Install()
}

func TestNewBrewInstaller(t *testing.T) {
    
    manifest := BrewManifest{
        Taps: []string{
            "homebrew/bundle",
            "homebrew/cask-versions",
        },
        Brews: []string{
            "curl",
            "wget",
            "zsh",
        },
    }
    i := NewBrewInstaller(
        DryRun(),
        WithManifest(manifest),
        UseWriterOptions(
            WithTestWriter(&strings.Builder{}),
            UseStdout(),
        ),
    ).(brewInstaller)
    
    require.Equal(t, manifest, i.config.Manifest)
    require.Equal(t, true, i.config.DryRun)
    require.Equal(t, 2, len(i.config.WriterOptions))
    require.NotNil(t, i.writer)
}

func OptionTest(opts ...Option) {
    var config = &BrewConfig{}
    
    for _, opt := range opts {
        name, option := opt()
        println("Applying Option '" + name + "'")
        option(config)
    }
    
    println("done")
}

func TestInstall(t *testing.T) {
    
    i := NewBrewInstaller(
        func(config *BrewConfig) error {
            config.DryRun = true
            return nil
        },
        func(c *BrewConfig) error {
            c.Manifest = BrewManifest{
                Taps: []string{
                    "homebrew/bundle",
                    "homebrew/cask-versions",
                },
                Brews: []string{
                    "curl",
                    "wget",
                    "zsh",
                },
            }
            return nil
        },
    )
    
    i.Install()
}
