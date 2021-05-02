package xdg_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestXdg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Xdg Suite")
}
