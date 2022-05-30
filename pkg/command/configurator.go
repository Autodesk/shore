package command

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/Autodeskshore/pkg/project"

	"github.com/hashicorp/go-multierror"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type DefaultConfErr struct {
	Err error
}

func (c *DefaultConfErr) Error() string {
	return c.Err.Error()
}

type FlagErr struct {
	Err error
}

func (c *FlagErr) Error() string {
	return c.Err.Error()
}

// Returns the list of possible yaml extensions.
// Example taken from - https://qvault.io/golang/golang-constant-maps-slices/
func getExtensions() []string {
	return []string{"json", "yaml", "yml"}
}

func GetConfigFileOrFlag(p *project.Project, fileName string, flag string) ([]byte, error) {
	var errors FlagErr

	dir, err := p.GetProjectPath()

	if err != nil {
		return nil, err
	}
	// If a flag is set, no need to read the config.
	// The config will be provided by the flag.
	if flag != "" {
		// first test if this a JSON string.
		// This is a quick check, without touching the FS.
		if data, err := decodeJson(flag); err != nil {
			errors.Err = multierror.Append(errors.Err, err)
		} else {
			return data, nil
		}

		var manualConfigFile string

		// Don't search the shore project directory path if the file is pointing to an absolute path.
		if filepath.IsAbs(flag) {
			manualConfigFile = flag
		} else {
			manualConfigFile = filepath.Join(dir, flag)
		}

		if data, err := readConfigFile(p, manualConfigFile); err != nil {
			errors.Err = multierror.Append(errors.Err, fmt.Errorf("could not find file - '%v'", flag))
		} else {
			return data, nil
		}

		return nil, fmt.Errorf("failed to decode flag as either JSON or filepath with errors:\n%v", errors)
	}

	// If the flag isn't set we want to read the config file.
	var pathErrors DefaultConfErr
	for _, ext := range getExtensions() {
		data, err := readConfigFile(p, filepath.Join(dir, fmt.Sprintf("%s.%s", fileName, ext)))

		if err, ok := err.(*fs.PathError); err != nil && ok {
			pathErrors.Err = multierror.Append(pathErrors.Err, err.Unwrap())
			continue
		}

		if err != nil {
			return nil, err
		}

		return data, err
	}

	return nil, fmt.Errorf("%w", &pathErrors)
}

func readConfigFile(p *project.Project, filePath string) ([]byte, error) {
	data, err := afero.ReadFile(p.FS, filePath)
	extension := filepath.Ext(filePath)

	if err != nil {
		return nil, err
	}

	var config interface{}

	switch extension {
	// Test the JSON case separately because this `{"a": "a",}` is valid YAML-v1.2 but not valid JSON
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, err
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, err
		}
	}

	data, err = json.Marshal(config)
	return data, err
}

func decodeJson(flag string) ([]byte, error) {
	data := []byte(flag)

	var values interface{}
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("could not decode JSON - '%v' - %w", flag, err)
	}

	return data, nil
}
