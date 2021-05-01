//go:generate gonum -types=KindEnum -output=kind_enum.go
package plugins

type KindEnum struct {
	SDK string `enum:"SDK,installs updates and manages different sdks on your system"`
}

func IsValidKind(k Kind) bool {
	for _, name := range KindNames() {
		if k.Name() == name || k.String() == name {
			return true
		}
	}
	return false
}

func (k Kind) IsValid() bool {
	return IsValidKind(k)
}
