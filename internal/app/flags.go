package app

type GlobalFlags struct {
	// Verbose is whether to log debug statements.
	Verbose bool

	// Config is the path the the DevCtlConfig file
	Config string

	// AppVersion is the semver of the devctl binary
	AppVersion string
}
