package api

import (
    "fmt"
    yaml2 "gopkg.in/yaml.v2"
    "regexp"
    "strconv"
    "strings"
)

func main() {
}

type Config struct {
    Sdks     [] *Sdk     `yaml:"sdks"`
    Contexts map[string] *Context `yaml:"contexts"`
    Versions [] *Version `yaml:"versions"`
}

type Sdk struct {
    Name     string     `yaml:"name"`
    Path     string     `yaml:"path"`
    Target string `yaml:"target"`
 //   Versions []*Version `yaml:"versions"`
}

type Context struct {
    SdkId string `yaml:"-"`
    VersionId string `yaml:"version"`
    Sdk     *Sdk     `yaml:"-"`
    Version *Version `yaml:"-"`
    Path    string   `yaml:"path"`
}

type Version struct {
    Id      string `yaml:"id"`
    Version SemVer `yaml:"version"`
    Vendor  string `yaml:"vendor"`
    Path    string `yaml:"path"`
}

func NewVersion(version string, path string, vendor string) (Version, error) {

    semver, err := parseSemVer(version)
    if err != nil {
        return Version{}, err
    }

    id := fmt.Sprintf("%v-%v", vendor, semver.String())

    v := Version{
        Vendor: vendor,
        Path:   path,
        Id:     id,
        Version: semver,
    }

    return v, nil
}

func parseSemVer(version string) (SemVer, error) {

    pattern := `\d*\.\d*\.\d`
    semver := &SemVer{}

    reg, err := regexp.Compile(pattern)
    if err != nil {
        return *semver, err
    }
    reg.Match([] byte(version))

    parts := strings.Split(version, ".")

    major, err := strconv.Atoi(parts[0])
    minor, err := strconv.Atoi(parts[1])
    patch, err := strconv.Atoi(parts[2])

    if err != nil {
        return *semver, err
    }
    semver.Major = major
    semver.Minor = minor
    semver.Patch = patch

    return *semver, nil
}

type SemVer struct {
    Major int `yaml:"major"`
    Minor int `yaml:"minor"`
    Patch int `yaml:"patch"`
}

func (semVer SemVer) GoString() string {
    return semVer.String()
}

func (semVer SemVer) String() string {
    return fmt.Sprintf("%v.%v.%v", semVer.Major, semVer.Minor, semVer.Patch)
}

var _ fmt.Stringer = new(SemVer)
var _ fmt.GoStringer = new(SemVer)

// NewConfig is a convenience function that returns a new Config object with non-nil maps
func NewConfig() *Config {
    return &Config{
        Sdks:     [] *Sdk{},
        Contexts: map[string] *Context{},
        Versions: [] *Version{},
    }
}

func NewSdk() *Sdk {
    return &Sdk{}
}

func NewContext() *Context {
    return &Context{}
}



func (config *Config) toYaml() string {
    bytes, err := yaml2.Marshal(config)
    if err != nil {
        fmt.Printf("Unexpected error: %v", err)
    }
    output := string(bytes)
    output = strings.TrimSpace(output)
    return output
}
