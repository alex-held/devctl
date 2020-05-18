package main

import (
	. "fmt"
	. "github.com/alex-held/dev-env/manifest"
)

var imported = Manifest{
	Version: "3.2.202",
	SDK:     "dotnet",
	Variables: map[string]interface{}{
		"url":          "https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/dotnet-sdk-3.1.202-osx-x64.tar.gz", //nolint:lll
		"install-root": "[[_sdks]]/[[sdk]]/[[version]]",
		"link-root":    "/usr/local/share/dotnet",
	},
	Install: Install{
		Commands: []string{
			"mkdir -p [[install-root]]",
			"curl [[url]] | tar -C [[install-root]] -x",
		},
	},
	Links: []Link{
		{Source: "[[install-root]]/host/fxr", Target: "[[link-root]]/host/fxr"},
		{Source: "[[install-root]]/sdk/[[version]]", Target: "[[link-root]]/sdk/[[version]]"},
		{Source: "[[install-root]]/shared/Microsoft.NETCore.App", Target: "[[link-root]]/shared/Microsoft.NETCore.App/[[version]]"},
		{Source: "[[install-root]]/shared/Microsoft.AspNetCore.All", Target: "[[link-root]]/shared/Microsoft.AspNetCore.All/[[version]]"},
		{Source: "[[install-root]]/shared/Microsoft.AspNetCore.App", Target: "[[link-root]]/shared/Microsoft.AspNetCore.App/[[version]]"},
	},
}

func main() {
	Println(imported.Format())

	PrintJson(imported)
	PrintYaml(imported)
}
