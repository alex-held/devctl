package cli

type DevEnvDir string

func ListDevEnvDirs() []DevEnvDir {
	return []DevEnvDir{SDKS, Downloads, Lib, Bin, Tmp}
}

const (
	SDKS      DevEnvDir = "sdks"
	Downloads DevEnvDir = "downloads"
	Lib       DevEnvDir = "lib"
	Bin       DevEnvDir = "bin"
	Tmp       DevEnvDir = "tmp"
)
