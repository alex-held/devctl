package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

import "github.com/thatisuday/commando"

func main() {

	commando.
		SetExecutableName("java-home").
		SetVersion("v1.0.0").
		SetDescription("This CLI tool helps you create and manage your local java environment.")

	commando.
		Register("current").
		SetDescription("Displays the current configuration for your local java environment.").
		SetShortDescription("prints the current config").
		AddArgument("sdk", "the sdk to print the current configuration from", "").
		SetAction(func(args map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {

			// print arguments
			for k, v := range args {
				fmt.Printf("arg -> %v: %v(%T)\n", k, v.Value, v.Value)
			}

			sdk := args["sdk"].Value

			var directory string
			baseDir := ".dev-env"
			homeDir, homeDirError := os.UserHomeDir()

			if homeDirError != nil {
				fmt.Errorf("Could not resolve the user home directory.\nError=%v", homeDirError.Error())
				os.Exit(1)
			}

			sdkDir := path.Join(homeDir, baseDir, "sdk")
			fmt.Printf("Resolved dev-env sdk root directory %v", sdkDir)

			switch sdk {
			case "java":
				directory = "jdk"
			case "javafx":
				directory = "javafx"
			default:
				fmt.Errorf("%v is not a valid sdk", sdk)
				os.Exit(1)
			}

			fmt.Printf("sdk = %v", sdk)
			fmt.Printf("directory = %v", directory)

		})

	//  examplePtr := flag.String("example", "defaultValue", " Help text.")
	textPtr := flag.String("text", "", "Text to parse (Required)")
	metricPtr := flag.String("metric", "chars", "Metric {chars|words|lines};. (Required)")
	uniquePtr := flag.Bool("unique", false, "Measure unique values of a metric.")

	if *textPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	fmt.Printf("textPtr: %s, metricPtr: %s, uniquePtr: %t\n", *textPtr, *metricPtr, *uniquePtr)
}
