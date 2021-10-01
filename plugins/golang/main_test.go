package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"

	"github.com/alex-held/devctl/pkg/plugins"
)

func TestExecuteGoPlugin(t *testing.T) {
	out := &bytes.Buffer{}
	sut := plugins.NewEngine(func(c *plugins.Config) *plugins.Config {
		c.Out = out
		c.Fs = afero.NewOsFs()
		c.Pather = devctlpath.DefaultPather()
		return c
	})

	p, err := sut.LoadPlugin("plugin.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, p)

	err = sut.Execute("go", []string{"current"})
	assert.NoError(t, err)
	assert.Equal(t, "v1.16.8\n", out.String())
}

func TestHandleCurrent(t *testing.T) {
	var err error

	out := captureStdout(func() {
		err = handleCurrent(&Config{InstallPath: "/Users/dev/.devctl/sdks/go"})
	})

	assert.NoError(t, err)
	assert.Equal(t, "v1.16.8\n", out)
}

func TestHandleList(t *testing.T) {
	var err error
	expected := "v1.13.5\nv1.16\nv1.16.3\nv1.16.4\nv1.16.8\nv1.17\nv1.17.1\n"

	out := captureStdout(func() {
		err = handleList(&Config{InstallPath: "/Users/dev/.devctl/sdks/go"})
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestHandleInstall(t *testing.T) {
	var err error
	expected := ""
	const installPath = "/Users/dev/.devctl/sdks/go"
	const version = "1.17.1"
	defer func() {
		_ = os.RemoveAll(path.Join(installPath, version))
	}()

	out := captureStdout(func() {
		err = handleInstall(version, &Config{InstallPath: installPath, Fs: afero.NewOsFs()})
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}


func TestHandleUse(t *testing.T) {
	var err error
	expected := ""
	const installPath = "/Users/dev/.devctl/sdks/go"
	const version = "1.17.1"

	defer func() {
		_ = os.RemoveAll(path.Join(installPath, version))
	}()

	out := captureStdout(func() {
		err = handleInstall(version, &Config{InstallPath: installPath, Fs: afero.NewOsFs()})
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}

func captureStdout(action func()) (out string) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	action()

	w.Close()
	o, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	return string(o)
}
