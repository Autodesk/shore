package jsonnet

import (
	"fmt"
)

// SharedLibErr - Custom error for the custom implementation of JSONNET shared libraries.
type SharedLibErr struct {
	Require string
	Path    string
	Err     error
}

func (s SharedLibErr) Error() string {
	return fmt.Sprintf("name: %s, path: %s, err: %v", s.Require, s.Path, s.Err)
}
