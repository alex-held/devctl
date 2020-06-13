package meta

import (
	. "github.com/alex-held/dev-env/meta"
)


const MetaGoyaml string = `
{{ $package := "dotnet" }}
{{ $version := "3.1.202" }}
{{ $InstallRoot := .GetPkgPath $package $version }}
{{ $LinkRoot := "/usr/local/share/dotnet" }}
package:
	name: {{ $package }}
	version: {{ $version }}
	install_root:  {{  $InstallRoot }}
	link_root: {{  $LinkRoot }}

sources:{{ range .Sources }}
- sha256: {{ .Sha256 }}{{if .URL}}
  url: {{.URL}}{{end}}{{if .Folder}}
  folder: {{.Folder}}{{end}}{{end}}

install: {{range .Install}}
- {{.}}{{end}}

link: {{range .Link}}
- {{.}}{{end}}

about:
	homepage: {{.Homepage}}
	summary: {{.Summary}}
`


func NewDotnetMeta() Meta {
	meta := Meta{
		Name:        "dotnet",
		Version:     "3.1.202",
		InstallRoot: "/User/dev/.devenv/pkg/dotnet/3.1.202",
		LinkRoot:    "/usr/local/share/dotnet",
		Sources: []Source{
			NewRemoteArchiveSource(
				"e67b13b4d6aaf6198188efc2f2c09531555ddbe1",
				"https://download.visualstudio.microsoft.com/download/pr/08088821-e58b-4bf3-9e4a-2c04448eee4b/e6e50aff8769ad382ed279730405ee3e/{{$package}}-sdk-{{$version}}-osx-x64.tar.gz", //nolint:lll
			),
		},
		Install: []string{
			"curl https://download.visualstudio.microsoft.com | tar -x -C ~/.devenv/pkg/dotnet/3.1.202",
		},
		Link: []string{
			"ln -s [[install-root]]/host/fxr [[link-root]]/host/fxr'",
		},
		Homepage: "https://github.com/microsoft/dotnet",
		Summary:  "Getting image size from png/jpeg/jpeg2000/gif file",
	}
	return meta
}
