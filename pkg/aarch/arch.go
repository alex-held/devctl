package aarch

// Arch The Processor Architecture the CLI is running at
type Arch string

const (
	// MAC_OSX
	MAC_OSX Arch = "darwinx64"
	
	// LINUX_64
	LINUX_64 Arch = "linuxx64"
	
	// LINUX_ARM32
	LINUX_ARM32 Arch = "linuxarm32"
)
