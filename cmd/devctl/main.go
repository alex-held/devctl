package main

import (
	"context"
	"log"
	"os"

	exec2 "github.com/alex-held/devctl/pkg/plugins/exec"
)

func main() {
	ctx := context.Background()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := exec2.Run(ctx, pwd, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
