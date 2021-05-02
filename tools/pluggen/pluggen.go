package main

import (
	"fmt"
	"os"

	"github.com/alex-held/devctl/tools/pluggen/cmd"
)

func main() {
	err := cmd.NewRootCmd().Execute()
	if err != nil {
		fmt.Printf("[tools/pluggen] %v\n", err)
		os.Exit(1)
	}
}
