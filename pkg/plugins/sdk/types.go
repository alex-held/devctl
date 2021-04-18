package sdk

type SDKPlugin interface {
	Name() string
	ListVersions() []string
	InstallE(version string) (err error)
	Download(version string) (err error)
	Install(version string)
	NewFunc() SDKPlugin
}
