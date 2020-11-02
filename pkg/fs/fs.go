/*
Package fs - A FileSystem abstraction layer for easy TDD testing with FileSystem needs.
The package is expected to be initialized with the `initFs` function before it is used in the project.
Currently supported FileSystem abstractions, OS (Currently running operating system), MEM (in-memeory file system abstraction)
Singleton concept from: https://gist.github.com/jeffotoni/da52ac6a1eaee7aae3c4fd17a8f8826a#file-golang-singleton-new-3-go
*/
package fs

import (
	"sync"

	"github.com/spf13/afero"
)

var (
	fs   afero.Fs
	once sync.Once
)

const (
	// OS File System type
	OS = iota
	// MEM File System type
	MEM
)

// InitFs - initialize the FileSystem abstraction
func InitFs(fsType int) afero.Fs {
	once.Do(func() {
		switch fsType {
		case OS:
			fs = afero.NewOsFs()
		case MEM:
			fs = afero.NewMemMapFs()
		}
	})

	return fs
}

// GetFs - Get the currently running filesystem singleton.
func GetFs() afero.Fs {
	return fs
}
