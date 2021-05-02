package devctlpath_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDevctlpath(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Devctlpath Suite")
}
