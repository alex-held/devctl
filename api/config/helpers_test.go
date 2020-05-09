package api

import (
	"fmt"
	semver2 "github.com/blang/semver"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestIsConfigEmpty(t *testing.T) {

	emptyConfig := NewConfig()
	isEmpty := isConfigEmpty(emptyConfig)

	if !isEmpty {
		t.Errorf("The provided Config is not empty. (expected empty)")
		return
	}

	fmt.Printf("\n%v\n", emptyConfig.toConfigFileYaml())
}

func TestParseYaml(t *testing.T) {

	sourceBytes, err := ioutil.ReadFile(filepath.Join("..", "..", "testdata/config.yaml"))
	sourceYaml := string(sourceBytes)

	if err != nil {
		t.Errorf("Unable to read the testdata 'testdata/config.yaml'")
	}

	fmt.Printf("\n[DATA: testdata/config.yaml]\n\n%v\n\n", sourceYaml)

	config, parseError := parseYamlFromConfigFile(sourceYaml)

	if parseError != nil {
		t.Errorf("Unable to parse the testdata 'testdata/config.yaml' into a Config")
	}

	printTestCase(sourceYaml, config.toConfigFileYaml())

	config.assertVersion(t, 0, "openjdk-14.0.1", 14, 0, 1, "openjdk", "openjdk-14.0.1/Contents/Home")
	config.assertSdk(t, 0, "jdk", "sdk/jdk", "/Library/Java/JavaVirtualMachines/")
	config.assertContainsContext(t, "jdk", "openjdk-14.0.1", "/Users/dev/.dev-env/sdk/jdk/openjdk-14.0.1/Contents/Home")
}

func (config Config) assertContainsContext(t *testing.T, key string, version string, path string) {

	context := config.Contexts[key]

	fmt.Printf("\n[ASSERT: SDK]\n\n")
	fmt.Printf("SdkId: Expected=%v {}; Actual=%v\n", key, context.SdkId)
	fmt.Printf("VersionId: Expected=%v {}; Actual=%v\n", version, context.VersionId)
	fmt.Printf("Path: Expected=%v {}; Actual=%v\n", path, context.Path)

	if !(context.SdkId == key && context.VersionId == version && context.Path == path) {
		t.Errorf("\nActual Sdk is not equal to expected parameters!\n")
	}

	fmt.Printf("\n\n")
}

func (config Config) assertSdk(t *testing.T, index int, name string, path string, target string) {
	sdk := config.Sdks[index]

	fmt.Printf("\n[ASSERT: SDK]\n\n")
	fmt.Printf("Name: Expected=%v {}; Actual=%v\n", name, sdk.Name)
	fmt.Printf("Path: Expected=%v {}; Actual=%v\n", path, sdk.Path)
	fmt.Printf("Target: Expected=%v {}; Actual=%v\n", target, sdk.Target)

	if !(sdk.Name == name && sdk.Path == path && sdk.Target == target) {
		t.Errorf("\nActual Sdk is not equal to expected parameters!\n")
	}

	fmt.Printf("\n\n")
	return
}

func (version Version) assertVersion(t *testing.T, id string, major int, minor int, patch int, vendor string, path string) {
	semver := version.Version

	fmt.Printf("\n[ASSERT: SDK]\n\n")
	fmt.Printf("Id: Expected=%v {}; Actual=%v\n", id, version.Id)
	fmt.Printf("Path: Expected=%v {}; Actual=%v\n", path, version.Path)
	fmt.Printf("Vendor: Expected=%v {}; Actual=%v\n", vendor, version.Vendor)
	fmt.Printf("Major: Expected=%v {}; Actual=%v\n", major, version.Version.Major)
	fmt.Printf("Minor: Expected=%v {}; Actual=%v\n", minor, version.Version.Minor)
	fmt.Printf("Patch: Expected=%v {}; Actual=%v\n", patch, version.Version.Patch)

	if !(version.Id == id && version.Vendor == vendor && version.Path == path && semver.Equals(semver2.MustParse(fmt.Sprintf("%v.%v.%v", major, minor, patch)))) {
		t.Errorf("\nActual Version is not equal to expected parameters!\n")
	}

	fmt.Printf("\n\n")
}

func (config Config) assertVersion(t *testing.T, index int, id string, major int, minor int, patch int, vendor string, path string) {
	version := config.Versions[index]
	version.assertVersion(t, id, major, minor, patch, vendor, path)
}

func printTestCase(expected string, actual string) {
	fmt.Printf("\n[EXPECTED]\n\n%v\n\n", expected)
	fmt.Printf("\n[ACTUAL]\n\n%v\n\n", actual)
}
