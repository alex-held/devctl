package api

import (
	"fmt"
	"github.com/blang/semver"
)

type Config struct {
	Sdks     []*Sdk              `yaml:"sdks"`
	Contexts map[string]*Context `yaml:"contexts"`
	Versions []*Version          `yaml:"versions"`
}

func (config *Config) String() string {
	return config.GoString()
}

func (config *Config) GoString() string {
	return fmt.Sprintf("%#v", config.toYaml())
}

type Sdk struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Target string `yaml:"target"`
}

type Context struct {
	SdkId     string   `yaml:"-"`
	VersionId string   `yaml:"version"`
	Sdk       *Sdk     `yaml:"-"`
	Version   *Version `yaml:"-"`
	Path      string   `yaml:"path"`
}

type Version struct {
	Id      string         `yaml:"id"`
	Version semver.Version `yaml:",inline"`
	Vendor  string         `yaml:"vendor"`
	Path    string         `yaml:"path"`
}

func NewVersion(version string, path string, vendor string) Version {

	semVer := semver.MustParse(version)

	id := fmt.Sprintf("%v-%v", vendor, semVer.String())

	v := Version{
		Vendor:  vendor,
		Path:    path,
		Id:      id,
		Version: semVer,
	}

	return v
}

/*

func (version Version) String() string {

   bytes, err := yaml.Marshal(version.Version)
    if err != nil {
        return "<nil>"
    }

   v := string(bytes)
   return v
}

func (v Version) GoString() string {
   return v.String()
}

var _ fmt.Stringer = new(Version)
var _ fmt.GoStringer = new(Version)
*/
/*
func parseSemVer(version string) (SemVer, error) {

	pattern := `\d*\.\d*\.\d*`
	semver := &SemVer{}

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return *semver, err
	}
	reg.Match([]byte(version))

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
*/
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
		Sdks:     []*Sdk{},
		Contexts: map[string]*Context{},
		Versions: []*Version{},
	}
}

// NewSdk is a convenience function that returns a new Config object with non-nil maps
func NewSdk() *Sdk {
	return &Sdk{}
}

// NewContext is a convenience function that returns a new Config object with non-nil maps
func NewContext() *Context {
	return &Context{}
}
