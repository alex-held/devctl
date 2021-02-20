package action

import (
	"github.com/alex-held/devctl/internal/sdkman"
)

// Configuration injects the dependencies that all action share.
type Configuration struct {

	// RegistryClient is a client for working with registries
	Registry *sdkman.RegistryService

	// Capabilities describes the capabilities of the Kubernetes cluster.
	// Capabilities *chartutil.Capabilities

	Log func(string, ...interface{})
}
