package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Autodeskshore/pkg/command"
	"github.com/Autodeskshore/pkg/project"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// ShoreConfig - A structure representing the Shore Config.
type ShoreConfig struct {
	Renderer map[string]interface{} `json:"renderer"`
	Executor map[string]interface{} `json:"executor"`
	Profiles map[string]interface{} `json:"profiles"`
}

// ProfileMetaData - Metadata for a given profile.
type ProfileMetaData struct {
	Application string
	Pipeline    string
}

func findIndividualLocalConfig(p *project.Project, configName string) (string, error) {
	projectPath, err := p.GetProjectPath()
	if err != nil {
		return "", err
	}

	files, err := afero.ReadDir(p.FS, projectPath)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		baseFileName := strings.ToLower(strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())))
		if baseFileName == configName {
			return path.Join(projectPath, f.Name()), nil
		}
	}
	return "", nil
}

func writeShoreConfig(p *project.Project, shoreConfig ShoreConfig) error {
	shoreConfigBytes, err := yaml.Marshal(shoreConfig)
	if err != nil {
		return err
	}

	projectPath, err := p.GetProjectPath()
	if err != nil {
		return err
	}

	shoreConfigPath := path.Join(projectPath, "shore.yml")
	afero.WriteFile(p.FS, shoreConfigPath, shoreConfigBytes, os.ModePerm)
	p.Log.Debugf("Created a new shore config located at: %s", shoreConfigPath)
	return nil
}

// LoadShoreConfig - Loads the Shore Config given a Project obj.
func LoadShoreConfig(p *project.Project) (ShoreConfig, error) {
	var shoreConfig ShoreConfig

	configData, err := command.GetConfigFileOrFlag(p, "shore", "")

	if err == nil {
		if err := jsoniter.Unmarshal(configData, &shoreConfig); err != nil {
			// Couldn't parse the file, so error out
			return ShoreConfig{}, err
		}

		p.Log.Debug("Loaded from file Shore Config: ", shoreConfig)
		return shoreConfig, nil
	}

	if strings.Contains(fmt.Sprint(err), "file does not exist") ||
		strings.Contains(fmt.Sprint(err), "no such file") {
		// Couldn't find the file, so provide a default.

		homeDirName, err := os.UserHomeDir()
		if err != nil {
			return ShoreConfig{}, err
		}

		defaultRenderer := map[string]interface{}{
			"type": "jsonnet",
		}

		defaultExecutor := map[string]interface{}{
			"type": "spinnaker",
			"config": map[string]interface{}{
				"default": fmt.Sprintf("%s/.spin/config", homeDirName),
			},
		}

		existingRenderPath, err := findIndividualLocalConfig(p, "render")
		if err != nil {
			return ShoreConfig{}, err
		} else if len(existingRenderPath) == 0 {
			return ShoreConfig{}, fmt.Errorf(`unable to find a render config in the project`)
		}

		existingExecPath, err := findIndividualLocalConfig(p, "exec")
		if err != nil {
			return ShoreConfig{}, err
		} else if len(existingExecPath) == 0 {
			return ShoreConfig{}, fmt.Errorf(`unable to find a exec config in the project`)
		}

		existingE2EPath, err := findIndividualLocalConfig(p, "e2e")
		if err != nil {
			return ShoreConfig{}, err
		} else if len(existingE2EPath) == 0 {
			return ShoreConfig{}, fmt.Errorf(`unable to find a E2E config in the project`)
		}

		defaultProfiles := map[string]interface{}{
			"default": map[string]interface{}{
				"render": existingRenderPath,
				"exec":   existingExecPath,
				"e2e":    existingE2EPath,
			},
		}

		shoreConfig = ShoreConfig{
			Renderer: defaultRenderer,
			Executor: defaultExecutor,
			Profiles: defaultProfiles,
		}

		p.Log.Debug("Loaded default Shore Config: ", shoreConfig)

		writeShoreConfig(p, shoreConfig)
		if err != nil {
			return ShoreConfig{}, err
		}

		return shoreConfig, nil
	}

	// Error out on anything other than missing file.
	return ShoreConfig{}, err
}
