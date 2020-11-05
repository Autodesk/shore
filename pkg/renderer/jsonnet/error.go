package jsonnet

import (
	"fmt"
	"strings"
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

// SharedLibsErr - Custom error for the custom implementation of JSONNET shared libraries.
type SharedLibsErr struct {
	MissingLibs []SharedLibErr
}

func (s SharedLibsErr) Error() string {
	var errors []string

	for _, err := range s.MissingLibs {
		errors = append(errors, err.Error())
	}

	return fmt.Sprintf("missing shared libraries errors:\n%s", strings.Join(errors, ","))
}
