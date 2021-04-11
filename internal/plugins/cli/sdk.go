package cli

type SdkPlugin interface {
	Download(version string) (err error)
}
