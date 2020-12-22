package project

import (
	"fmt"
	"os"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
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
	log         logrus.FieldLogger
	projectPath string
}

// NewShoreProject creates a new shore project with default values
// `projectPath` may be used to override the default project path with is the execution directory (os.Getwd)
func NewShoreProject(fs afero.Fs, logger logrus.FieldLogger, projectPath ...string) *Project {
	if len(projectPath) > 0 {
		return &Project{
			FS:          fs,
			projectPath: projectPath[0],
			log:         logger,
		}
	}

	return &Project{
		FS:  fs,
		log: logger,
	}
}

func (p *Project) GetProjectPath() (string, error) {
	// Magic variable to allow working shore without actually being in the path
	// For testing purposes only!! (probably?)
	if p.projectPath != "" {
		return p.projectPath, nil
	}

	// Magic variable to allow working shore without actually being in the path
	// For dev purposes only!!
	p.log.Debug("Lazy Loading project path")
	if isLocal := viper.GetBool("LOCAL"); isLocal == true {
		p.log.Debug("Found local development setup, using `SHORE_PROJECT_PATH` env variable")
		projectPath := viper.GetString("SHORE_PROJECT_PATH")
		p.log.Debug("`SHORE_PROJECT_PATH` set too ", projectPath)
		if projectPath == "" {
			return "", fmt.Errorf("env variable `SHORE_PROJECT_PATH` is not set")
		}

		return projectPath, nil
	}
	p.log.Debug("Use `Getwd` (pwd) local path")
	return os.Getwd()
}

// GetRenderArgs returns the contents of the Render Args file.
func (p *Project) GetRenderArgs() (string, error) {
	// This method causes a Marshal->UnMarshal
	// ReMarshaling here makes the code easier to follow but hurts performance.
	p.log.Debug("`GetRenderArgs` was called")
	return p.readConfigFile("render")
}

// GetExecArgs returns the contents of the Exec Args file.
func (p *Project) GetExecArgs() (string, error) {
	// This method causes a Marshal->UnMarshal
	// ReMarshaling here makes the code easier to follow but hurts performance.
	p.log.Debug("`GetExecArgs` was called")
	return p.readConfigFile("exec")
}

// GetTestConfig returns the contents of the Test Config file.
func (p *Project) GetTestConfig() (string, error) {
	// This method causes a Marshal->UnMarshal
	// ReMarshaling here makes the code easier to follow but hurts performance.
	p.log.Debug("`GetTestConfig` was called")
	return p.readConfigFile("E2E")
}

// The same entrypoint for all readConfigFile data.
func (p *Project) readConfigFile(filename string) (string, error) {
	p.log.Debug("`readConfigFile` was called")
	argsData := make(map[interface{}]interface{})

	projectPath, err := p.GetProjectPath()

	if err != nil {
		p.log.Error("`readConfigFile` failed with error", err)
		return "", err
	}

	for _, extension := range []string{"json", "yaml", "yml"} {
		p.log.Debug("Looking for file extension for filename ", filename, extension)
		filepath := filepath.Join(projectPath, fmt.Sprintf("%s.%s", filename, extension))
		exists, err := afero.Exists(p.FS, filepath)

		if err != nil || !exists {
			continue
		}
		p.log.Debug("reading file", filepath)
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
