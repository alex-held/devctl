package sdk

import (
	"plugin"
	"testing"

	"github.com/stretchr/testify/assert"
)

// //go:generate go build -buildmode=plugin -o "./testdata/sdk-01.so" ../../../plugins/testdata/sdk-01/main.go
//go:generate pluggen gen -o "./testdata/sdk-01.so" -p "plugins/testdata/sdk-01" --pkg devctl
func TestLoadSDKPlugin(t *testing.T) {

	plug, err := plugin.Open("testdata/sdk-01.so")
	if err != nil {
		t.Fatal(err)
	}

	sdkPlugin, errs := LookupSDKPlugin(plug)
	assert.Empty(t, errs)
	assert.NotNil(t, sdkPlugin)
	err = sdkPlugin.Install(nil, []string{"1.16.3"})
	err = sdkPlugin.Download(nil, []string{"1.16.3"})
	err = sdkPlugin.List(nil, []string{})
	assert.NoError(t, err)
}
