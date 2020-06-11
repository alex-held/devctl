module github.com/alex-held/dev-env

go 1.14

require (
	github.com/Workiva/go-datastructures v1.0.52 // indirect
	github.com/disiqueira/gotree v1.0.0
	github.com/ganbarodigital/go_pipe/v5 v5.2.0
	github.com/ganbarodigital/go_scriptish v1.4.0
	github.com/ghodss/yaml v1.0.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/songgao/colorgo v0.0.0-20161028043718-1e1a5b5cef5c // indirect
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.5.1
	github.com/tinylib/msgp v1.1.2 // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/alex-held/dev-env/test => ./test
	github.com/alex-held/dev-env/utils => ./utils
)
