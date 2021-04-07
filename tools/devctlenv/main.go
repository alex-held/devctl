package main

import (
	"fmt"
	"os"
	"devctlenv/cli"
)

func main() {
	app := cli.NewCLIApp("0.0.1")
	fmt.Println(app .Version())

	shim, err := app .GetShimForTool(os.Args[1])
	if err != nil {
		_ = fmt.Sprintf("No shim has been found: arguments=%v", os.Args)
		os.Exit(1)
	}

	err = shim.Exec(os.Args[2:]...)
	if err != nil {
		_ = fmt.Sprintf("The underlieying tool exited with failure: arguments=%v\n%+v\n", os.Args, err)
		os.Exit(1)
	}
}
