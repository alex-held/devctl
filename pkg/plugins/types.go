package plugins

type SDKPlugin interface {
	Name() string
	ListVersions() []string
	Download(version string)
	Install(version string)
	NewFunc() SDKPlugin
}
