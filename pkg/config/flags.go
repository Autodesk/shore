package config

import (
	"fmt"
	"path/filepath"

	"github.com/Autodeskshore/pkg/project"
	"github.com/hashicorp/go-multierror"
	jsoniter "github.com/json-iterator/go"
)

// FlagConfigErr is an error that occurs when there's an issue with the configuration flags.
type FlagConfigErr struct {
	Err error
}

func (c *FlagConfigErr) Error() string {
	return c.Err.Error()
}

// GetFlagConfig - Returns a configuration from a flag.
func GetFlagConfig(p *project.Project, flag string) ([]byte, error) {
	var errors FlagConfigErr

	if flag == "" {
		errors.Err = multierror.Append(errors.Err, fmt.Errorf("empty flag"))
		return nil, fmt.Errorf("failed to decode flag as either JSON or filepath with errors:\n%v", errors)
	}

	// First test if this a JSON string.
	// This is a quick check, without touching the FS.
	if data, err := decodeJSON(flag); err != nil {
		errors.Err = multierror.Append(errors.Err, err)
	} else {
		return data, nil
	}

	// This means the flag is a file path - either absolute or relative.

	var manualConfigFile string

	// Don't search the shore project directory path if the file is pointing to an absolute path.
	if filepath.IsAbs(flag) {
		manualConfigFile = flag
	} else {
		dir, err := p.GetProjectPath()
		if err != nil {
			return nil, err
		}
		manualConfigFile = filepath.Join(dir, flag)
	}

	if data, err := ReadConfigFile(p, manualConfigFile); err != nil {
		errors.Err = multierror.Append(errors.Err, fmt.Errorf("could not find file - '%v'", flag))
	} else {
		return data, nil
	}

	return nil, fmt.Errorf("failed to decode flag as either JSON or filepath with errors:\n%v", errors)
}

func decodeJSON(flag string) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data := []byte(flag)

	var values interface{}
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("could not decode JSON - '%v' - %w", flag, err)
	}

	return data, nil
}
