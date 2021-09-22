module github.com/alex-held/devctl/plugins/golang

go 1.16

require (
	github.com/alex-held/devctl v0.9.7-beta
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/muesli/termenv v0.9.0 // indirect
	github.com/spf13/afero v1.6.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20210921065528-437939a70204 // indirect
	golang.org/x/text v0.3.7 // indirect
)


replace (
	github.com/alex-held/devctl v0.9.7-beta => ../..
)
