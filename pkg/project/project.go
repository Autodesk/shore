package project

import (
	"fmt"
	"os"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type ProjectOrganizer interface {
	GetProjectPath() (string, error)
	GetRenderArgs() (string, error)
	GetExecArgs() (string, error)
	GetTestConfig() (string, error)
}

type Project struct {
	ProjectOrganizer
	FS          afero.Fs
	projectPath string
}

func NewShoreProject(fs afero.Fs, projectPath ...string) *Project {
	if len(projectPath) > 0 {
		return &Project{
			FS:          fs,
			projectPath: projectPath[0],
		}
	}

	return &Project{
		FS: fs,
	}
}

func (p *Project) GetProjectPath() (string, error) {
	// Magic variable to allow working shore without actually beign in the path
	// For testing purposes only!! (probably?)
	if p.projectPath != "" {
		return p.projectPath, nil
	}

	// Magic variable to allow working shore without actually beign in the path
	// For dev purposes only!!
	if isLocal := viper.GetBool("LOCAL"); isLocal == true {
		projectPath := viper.GetString("SHORE_PROJECT_PATH")

		if projectPath == "" {
			return "", fmt.Errorf("env variable `SHORE_PROJECT_PATH` is not set")
		}

		return projectPath, nil
	}

	return os.Getwd()
}

// GetRenderArgs returns the contents of the Render Args file.
func (p *Project) GetRenderArgs() (string, error) {
	// This method causes a Marshal->UnMarshal
	// ReMarshaling here makes the code easier to follow but hurts perofmance.
	return p.readConfigFile("render")
}

// GetExecArgs returns the contents of the Exec Args file.
func (p *Project) GetExecArgs() (string, error) {
	// This method causes a Marshal->UnMarshal
	// ReMarshaling here makes the code easier to follow but hurts perofmance.
	return p.readConfigFile("exec")
}

// GetTestConfig returns the contents of the Test Config file.
func (p *Project) GetTestConfig() (string, error) {
	// This method causes a Marshal->UnMarshal
	// ReMarshaling here makes the code easier to follow but hurts perofmance.
	return p.readConfigFile("E2E")
}

// The same entrypoint for all readConfigFile data.
func (p *Project) readConfigFile(filename string) (string, error) {
	argsData := make(map[interface{}]interface{})

	projectPath, err := p.GetProjectPath()

	if err != nil {
		return "", err
	}

	for _, extension := range []string{"json", "yaml", "yml"} {
		filepath := filepath.Join(projectPath, fmt.Sprintf("%s.%s", filename, extension))
		exists, err := afero.Exists(p.FS, filepath)

		if err != nil || !exists {
			continue
		}

		argsBytes, err := afero.ReadFile(p.FS, filepath)

		if err != nil {
			return "", err
		}
		// Validate that the contents are valid JSON/YAML
		err = yaml.Unmarshal(argsBytes, &argsData)

		if err != nil {
			return "", err
		}
		// Turn the data back to a []byte to pass back as a string.
		args, err := jsoniter.Marshal(argsData)

		if err != nil {
			return "", err
		}

		return string(args), nil
	}

	// No file was found.
	return "", &os.PathError{Op: "open", Path: fmt.Sprintf("%s.[json|yaml|yml]", filename), Err: os.ErrNotExist}
}
