package api

import (
    "fmt"
    "strings"
    "testing"
)

func TestNewEmptyConfig(t *testing.T) {

    sourceConfig := NewConfig()
    expected := `
sdks: []
contexts: {}
versions: []`
    test := formatTest{inputConfig: sourceConfig, expected: expected}
    test.run(t)
}

func TestConfigWithSdk(t *testing.T) {

    sourceSdk := NewSdk()
    sourceSdk.Name =  "jdk"
    sourceSdk.Path = "sdk/jdk"
    sourceSdk.Target = "/Library/Java/JavaVirtualMachines/"

    sourceConfig := NewConfig()
    sourceConfig.Sdks = append(sourceConfig.Sdks, sourceSdk)

    expected := `
sdks:
- name: jdk
  path: sdk/jdk
  target: /Library/Java/JavaVirtualMachines/
contexts: {}
versions: []`

    test := formatTest{inputConfig: sourceConfig, expected: expected}
    test.run(t)
}



func TestFullConfig(t *testing.T) {

    sourceVersion, err := NewVersion("14.0.1", "openjdk-14.0.1/Contents/Home", "openjdk")
    if err != nil {
        t.Errorf("Unable to create type Version")
    }

    sourceContext := NewContext()
    sourceContext.SdkId =  "jdk"
    sourceContext.VersionId = "openjdk-14.0.1"
    sourceContext.Path = "/Users/dev/.dev-env/sdk/jdk/openjdk-14.0.1/Contents/Home"

    sourceSdk := NewSdk()
    sourceSdk.Name =  "jdk"
    sourceSdk.Path = "sdk/jdk"
    sourceSdk.Target = "/Library/Java/JavaVirtualMachines/"

    sourceConfig := NewConfig()
    sourceConfig.Sdks = append(sourceConfig.Sdks, sourceSdk)
    sourceConfig.Versions = [] *Version{&sourceVersion}
    sourceConfig.Contexts[sourceContext.SdkId] = sourceContext

    expected := `
sdks:
- name: jdk
  path: sdk/jdk
  target: /Library/Java/JavaVirtualMachines/
contexts:
  jdk:
    version: openjdk-14.0.1
    path: /Users/dev/.dev-env/sdk/jdk/openjdk-14.0.1/Contents/Home
versions:
- id: openjdk-14.0.1
  version:
    major: 14
    minor: 0
    patch: 1
  vendor: openjdk
  path: openjdk-14.0.1/Contents/Home`

    test := formatTest{inputConfig: sourceConfig, expected: expected}
    test.run(t)
}



func TestConfigWithVersion(t *testing.T) {

    sourceVersion, err := NewVersion("14.0.1", "openjdk-14.0.1/Contents/Home", "openjdk")
    if err != nil {
        t.Errorf("Unable to create type Version")
    }

    sourceConfig := NewConfig()
    sourceConfig.Versions = [] *Version{&sourceVersion}

    expected := `
sdks: []
contexts: {}
versions:
- id: openjdk-14.0.1
  version:
    major: 14
    minor: 0
    patch: 1
  vendor: openjdk
  path: openjdk-14.0.1/Contents/Home`

    test := formatTest{inputConfig: sourceConfig, expected: expected}
    test.run(t)
}


func TestConfigWithContext(t *testing.T) {

    sourceContext := NewContext()
    sourceContext.SdkId =  "jdk"
    sourceContext.VersionId = "openjdk-14.0.1"
    sourceContext.Path = "/Users/dev/.dev-env/sdk/jdk/openjdk-14.0.1/Contents/Home"

    sourceConfig := NewConfig()
    sourceConfig.Contexts[sourceContext.SdkId] = sourceContext

    expected := `
sdks: []
contexts:
  jdk:
    version: openjdk-14.0.1
    path: /Users/dev/.dev-env/sdk/jdk/openjdk-14.0.1/Contents/Home
versions: []`

    test := formatTest{inputConfig: sourceConfig, expected: expected}
    test.run(t)
}



func TestParseSemVer(t *testing.T) {

    version := "14.0.1"
    semver, err := parseSemVer(version)

    if err != nil {
        t.Error(err)
    }

    if !(semver.Major == 14 && semver.Minor == 0 && semver.Patch == 1) {
        t.Errorf("SemVer has been parsed incorrectly!Major=%v\nMinor=%v\nPatch=%v",semver.Major,semver.Minor,semver.Patch )
    }

    fmt.Printf("Major=%v\nMinor=%v\nPatch=%v\n",semver.Major,semver.Minor,semver.Patch)
}


func TestParseSemVer2(t *testing.T) {

    version := "1.034.12341234"
    semver, err := parseSemVer(version)

    if err != nil {
        t.Error(err)
    }

    if !(semver.Major == 1 && semver.Minor == 34 && semver.Patch == 12341234) {
        t.Errorf("SemVer has been parsed incorrectly!Major=%v\nMinor=%v\nPatch=%v",semver.Major,semver.Minor,semver.Patch )
    }

    fmt.Printf("Major=%v\nMinor=%v\nPatch=%v\n",semver.Major,semver.Minor,semver.Patch)
}

func TestParseFail(t *testing.T) {

    version := "1.wd.12"
    _, err := parseSemVer(version)

    if err == nil {
        t.Error(err)
    }

    t.Logf(err.Error())
}

func TestNewVersion(t *testing.T) {

    version, err := NewVersion("14.0.1", "openjdk-14.0.1", "openjdk")
    if err != nil {
        t.Error(err)
    }

    semver := version.Version
    if !(version.Vendor == "openjdk" && version.Path == "openjdk-14.0.1" && semver.Major == 14 && semver.Minor == 0 && semver.Patch == 1) {
        t.Errorf("\nPath=%v\nVendor=%v\nId=%v\nMajor=%v\nMinor=%v\nPatch=%v",
            version.Path,
            version.Vendor,
            version.Id,
            semver.Major,
            semver.Minor,
            semver.Patch)
    }

    t.Logf("\nPath=%v\nVendor=%v\nId=%v\nMajor=%v\nMinor=%v\nPatch=%v",
        version.Path,
        version.Vendor,
        version.Id,
        semver.Major,
        semver.Minor,
        semver.Patch)
}











type formatTest struct {
    inputConfig *Config
    expected    string
}

type configTest struct {
    inputConfig     *Config
    expectedConfig  *Config
    expectedOutputs [] string
}



func (test configTest) checkOutput(out string, expectedOutputs []string, t *testing.T) {
    for _, expectedOutput := range expectedOutputs {
        if !strings.Contains(out, expectedOutput) {
            t.Errorf("expected '%s' in output, got '%s'", expectedOutput, out)
        }
    }
}



func (config *formatTest) run(t *testing.T) string {

    expected := strings.TrimSpace(config.expected)
    output := config.inputConfig.toYaml()

    fmt.Printf("\n[EXPECTED]\n\n%v\n\n", expected)
    fmt.Printf("\n\n[ACTUAL]\n\n%v\n\n", output)

    if output != expected { t.Fail()}
    return output
}

