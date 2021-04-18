package cmd

// Version is the semver of the devctl binary
var Version = "development"

var (
	// verbose is whether to log debug statements
	verbose bool

	// verbose is the path the the DevCtlConfig file
	config string
)
