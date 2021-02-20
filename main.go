package main

import (
	"github.com/alex-held/devctl/cmd"
)

func main() {
	/*
		header := &doc.GenManHeader{
			Title:   "MINE",
			Section: "3",
		}
		err := doc.GenManTree(cmd.rootCmd, header, "/Users/dev/.devctl/tmp")
		if err != nil {
			log.Default().Fatal(err)
		}*/
	cmd.Execute()
}
