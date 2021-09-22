package main

import (
	"fmt"
	"os"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
	"github.com/alex-held/devctl/pkg/plugins"
)

func main() {
	devctlRoot := os.Args[1]
	fmt.Printf("DEVCTL_ROOT=%s\n", devctlRoot)

	e := plugins.NewEngine(func(c *plugins.Config) *plugins.Config {
		c.Pather = devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
			return devctlRoot
		}))
		return c
	})

	plugins := e.LoadPlugins()
	for _, plugin := range plugins {
		fmt.Printf("loaded plugin %s\n", plugin.Name)
	}

	err := e.Execute("go", os.Args[1:])
	if err != nil {
		fmt.Printf("ERROR=%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
