package command

import (
	"path"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func GetConfigFileOrFlag(d *Dependencies, fileName string, flagName string) ([]byte, error) {
	dir, err := d.Project.GetProjectPath()

	if err != nil {
		return nil, err
	}
	// If a flag is set, no need to read the config.
	// The config will be provided by the flag.
	if viper.IsSet(flagName) {
		var manualConfigFile string

		// Check if the flag is pointing to a file.
		possiblePath := viper.GetString(flagName)
		// Don't merge paths if the file is pointing to an absolute path.
		if path.IsAbs(possiblePath) {
			manualConfigFile = possiblePath
		} else {
			manualConfigFile = path.Join(dir, possiblePath)
		}

		exists, err := afero.Exists(d.Project.FS, manualConfigFile)

		viper.SetConfigType("yaml")

		if err != nil {
			return nil, err
		}

		// If the string represents a file and it exists in the project path,
		// we want to read the config file as is.
		if exists {
			dir, manualConfigFileName := filepath.Split(manualConfigFile)
			fileName := strings.TrimSuffix(manualConfigFileName, filepath.Ext(manualConfigFileName))
			d.Logger.Debug("Looking for file: ", fileName)
			return readConfigFile(d, fileName, dir)
		}

		// If the string wasn't a file path, go
		values := viper.GetStringMap(flagName)
		return jsoniter.Marshal(values)
	}

	// If the flag isn't set we want to read the config file.
	return readConfigFile(d, fileName)
}

func readConfigFile(d *Dependencies, fileName string, searchPaths ...string) ([]byte, error) {
	// To allow for a clean config setup, we will create a new `Viper` instance.
	// This allows us to retrieve a clean `map[string]interface{}` without the accompanying flag configurations.
	v := viper.New()
	path, err := d.Project.GetProjectPath()

	if err != nil {
		return nil, err
	}
	// Reset the lookup path.
	v.SetFs(d.Project.FS)
	v.AddConfigPath(path)

	// Add additional paths if necessary.
	for _, s := range searchPaths {
		v.AddConfigPath(s)
	}

	v.SetConfigName(fileName)

	valuesErr := v.ReadInConfig()

	if valuesErr != nil {
		d.Logger.Error("Failed to load values")
		return nil, valuesErr
	}

	values := v.AllSettings()
	d.Logger.Debug(values)
	return jsoniter.Marshal(values)
}
