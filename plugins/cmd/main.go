package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
	"github.com/alex-held/devctl/plugins/zsh"
)

var inputFile = flag.String("f", "", "zsh/config.yaml")

func main() {
	configFile := os.Args[1]
	args := os.Args[2:]
	plugin := os.Args[2]

	switch plugin {
	case "zsh":

		f, err := os.Open(configFile)
		if err != nil {
			fmt.Printf("ERROR=%v\n", err)
			os.Exit(1)
		}

		cfg, err := zsh.ReadConfigFile(f)
		if err != nil {
			fmt.Printf("ERROR=%v\n", err)
			os.Exit(1)
		}

		cfg.Pather = devctlpath.DefaultPather()
		cfg.Out = os.Stdout
		cfg.Context.Context = context.Background()

		// fmt.Printf("CONFIG=%v\n", *cfg)

		err = zsh.Exec(cfg, args[1:])
		if err != nil {
			fmt.Printf("ERROR=%v\n", err)
			os.Exit(1)
		}
	}
}
