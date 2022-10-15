package config

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/Autodesk/shore/pkg/project"
	"github.com/hashicorp/go-multierror"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// FileConfErr - an error that occurs when there's an issue with the configuration files.
type FileConfErr struct {
	Err error
}

func (c *FileConfErr) Error() string {
	return c.Err.Error()
}

// Returns the list of possible yaml extensions.
// Example taken from - https://qvault.io/golang/golang-constant-maps-slices/
func getExtensions() []string {
	return []string{"json", "yaml", "yml"}
}

// GetFileConfig - Returns a configuration from a flag. Looks for .json/.yml/.yaml formats
func GetFileConfig(p *project.Project, fileName string) ([]byte, error) {
	dir, err := p.GetProjectPath()
	if err != nil {
		return nil, err
	}

	var pathErrors FileConfErr
	for _, ext := range getExtensions() {
		data, err := ReadConfigFile(p, filepath.Join(dir, fmt.Sprintf("%s.%s", fileName, ext)))

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

// ReadConfigFile - Reads in a config. Supports json/yaml/yml
func ReadConfigFile(p *project.Project, filePath string) ([]byte, error) {
	data, err := afero.ReadFile(p.FS, filePath)
	extension := filepath.Ext(filePath)

	if err != nil {
		return nil, err
	}

	var config interface{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
