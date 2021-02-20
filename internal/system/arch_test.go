package system

import (
	"fmt"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestGetCurrent(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("GetCurrent", func() {
		g.AfterEach(func() {
			getDefaultGoosRuntimeGoos = nil
		})

		g.Describe("GIVEN GOOS=linux", func() {
			g.JustBeforeEach(func() {
				getDefaultGoosRuntimeGoos = func() string {
					return OsLinux
				}
			})

			g.It("WHEN GetCurrent() => THEN returns arch.Linux 'linux' ", func() {
				expected := Linux
				actual := GetCurrent()
				Expect(actual).To(Equal(expected))
			})
		})

		g.Describe("GIVEN GOOS=darwin", func() {
			g.JustBeforeEach(func() {
				getDefaultGoosRuntimeGoos = func() string {
					return OsDarwin
				}
			})

			g.It("WHEN GetCurrent() => THEN returns arch.IsDarwin 'darwin' ", func() {
				expected := Darwin
				actual := GetCurrent()
				Expect(actual).To(Equal(expected))
			})
		})

		g.Describe("GIVEN GOOS=windows", func() {
			g.JustBeforeEach(func() {
				getDefaultGoosRuntimeGoos = func() string {
					return OsWindows
				}
			})

			g.It("WHEN GetCurrent() => THEN returns arch.Windows 'windows' ", func() {
				expected := Windows
				actual := GetCurrent()
				Expect(actual).To(Equal(expected))
			})
		})
	})
}

func TestArch_String(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	tests := map[string][]struct {
		Input    Arch
		Expected string
	}{
		"Darwin": {
			{
				Input:    Darwin,
				Expected: "darwin",
			},
			{
				Input:    DarwinX64,
				Expected: "darwinx64",
			},
		},
		"Linux": {
			{
				Input:    Linux,
				Expected: "linux",
			},
			{
				Input:    LinuxX64,
				Expected: "linuxx64",
			},
			{
				Input:    LinuxArm32,
				Expected: "linuxarm32",
			},
		},
		"Windows": {
			{
				Input:    Windows,
				Expected: "windows",
			},
		},
	}

	for key, tests := range tests {
		g.Describe(fmt.Sprintf("GIVEN %s", key), func() {
			for _, tc := range tests {
				g.It(fmt.Sprintf("WHEN Input=%s", tc.Input), func() {
					actual := tc.Input.String()
					Expect(actual).To(Equal(tc.Expected))
				})
			}
		})
	}
}

func TestArch_Is(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	var a Arch
	var family Arch

	g.Describe("GIVEN Darwin", func() {
		g.Describe("WHEN IsOfFamily", func() {
			g.JustBeforeEach(func() {
				family = Darwin
			})

			g.It("WHEN Arch=darwin => THEN return true", func() {
				a = Darwin
				actual := a.IsOfFamily(family)
				Expect(actual).To(BeTrue())
			})

			g.It("WHEN Arch=darwinx64 => THEN return true", func() {
				a = DarwinX64
				actual := a.IsOfFamily(family)
				Expect(actual).To(BeTrue())
			})
		})

		g.Describe("IsDarwin", func() {
			g.It("WHEN Arch=darwinx64 => THEN return true", func() {
				a = DarwinX64
				actual := a.IsDarwin()
				Expect(actual).To(BeTrue())
			})

			g.It("WHEN Arch=darwin => THEN return true", func() {
				a = Darwin
				actual := a.IsDarwin()
				Expect(actual).To(BeTrue())
			})

			g.It("WHEN Arch=linux => THEN return false", func() {
				a = Linux
				actual := a.IsDarwin()
				Expect(actual).To(BeFalse())
			})

			g.It("WHEN Arch=linuxx64 => THEN return false", func() {
				a = LinuxX64
				actual := a.IsDarwin()
				Expect(actual).To(BeFalse())
			})

			g.It("WHEN Arch=windows => THEN return false", func() {
				a = Windows
				actual := a.IsDarwin()
				Expect(actual).To(BeFalse())
			})
		})
	})

	g.Describe("GIVEN Linux", func() {
		g.Describe("WHEN IsOfFamily", func() {
			g.JustBeforeEach(func() {
				family = Linux
			})

			g.It("WHEN Arch=linux => THEN return true", func() {
				a = Linux
				actual := a.IsOfFamily(family)
				Expect(actual).To(BeTrue())
			})

			g.It("WHEN Arch=linuxx64 => THEN return true", func() {
				a = LinuxX64
				actual := a.IsOfFamily(family)
				Expect(actual).To(BeTrue())
			})

			g.It("WHEN Arch=linuxarm32 => THEN return true", func() {
				a = LinuxArm32
				actual := a.IsOfFamily(family)
				Expect(actual).To(BeTrue())
			})
		})

		g.Describe("IsLinux", func() {
			g.It("WHEN Arch=linuxx64  => THEN return true", func() {
				a = LinuxX64
				actual := a.IsLinux()
				Expect(actual).To(BeTrue())
			})

			g.It("WHEN Arch=linux => THEN return true", func() {
				a = Linux
				actual := a.IsLinux()
				Expect(actual).To(BeTrue())
			})

			g.It("WHEN Arch=linuxarm32 => THEN return true", func() {
				a = LinuxArm32
				actual := a.IsLinux()
				Expect(actual).To(BeTrue())
			})

			g.It("WHEN Arch=darwin => THEN return false", func() {
				a = Darwin
				actual := a.IsLinux()
				Expect(actual).To(BeFalse())
			})

			g.It("WHEN Arch=darwinx64 => THEN return false", func() {
				a = DarwinX64
				actual := a.IsLinux()
				Expect(actual).To(BeFalse())
			})

			g.It("WHEN Arch=windows => THEN return false", func() {
				a = Windows
				actual := a.IsLinux()
				Expect(actual).To(BeFalse())
			})
		})
	})
}
