package main

import (
	"github.com/alex-held/devctl/pkg/cli"
	"github.com/alex-held/devctl/pkg/cli/util"
)

func main() {
	cmd := cli.NewDefaultKubectlCommand()
	util.CheckErr(cmd.Execute())
}
