package aarch

// Arch The Processor Architecture the CLI is running at
type Arch string

const (
	// MacOsx
	MacOsx Arch = "darwinx64"

	// Linux64
	Linux64 Arch = "linuxx64"

	// LinuxArm32
	LinuxArm32 Arch = "linuxarm32"
)
