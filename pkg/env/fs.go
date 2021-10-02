package env

import (
	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()

// SetFs sets the default afero.Fs
func SetFs(_fs afero.Fs) {
	fs = _fs
}

// IsOsFs returns true if the afero.Fs is of type *afero.OsFs
func IsOsFs(f afero.Fs) bool {
	_, ok := f.(*afero.OsFs)
	return ok
}

// IsMemMapFs returns true if the afero.Fs is of type *afero.MemMapFs
func IsMemMapFs(f afero.Fs) bool {
	_, ok := f.(*afero.MemMapFs)
	return ok
}

// GetFs gets the default afero.Fs
func GetFs() afero.Fs {
	return fs
}
