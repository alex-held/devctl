package system

import (
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestRuntimeInfo_Format(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	const os = "darwin"
	const arch = "amd64"
	var sut = RuntimeInfo{
		OS:   os,
		Arch: arch,
	}

	g.Describe("GIVEN pattern contains just runtime info templates ", func() {
		g.It("WHEN pattern contains [os] => THEN [os] gets replaced", func() {
			actual := sut.Format("/some/filename.1.32.3[os]")
			Expect(actual).Should(Equal("/some/filename.1.32.3" + os))
		})
	})
}
