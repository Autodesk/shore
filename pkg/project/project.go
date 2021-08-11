package project

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Project struct {
	FS   afero.Fs
	Log  *log.Logger
	Path string
}

func NewShoreProject(fs afero.Fs, logger *log.Logger) *Project {
	return &Project{
		FS:  fs,
		Log: logger,
	}
}

func (p *Project) GetProjectPath() (string, error) {
	// Magic variable to allow working shore without actually being in the path
	// For testing purposes only!! (probably?)
	if p.Path != "" {
		return p.Path, nil
	}

	// Magic variable to allow working shore without actually being in the path
	// For dev purposes only!!
	p.Log.Debug("Lazy Loading project path")

	if isLocal := os.Getenv("LOCAL"); isLocal == "true" {
		p.Log.Debug("Found local development setup, using `SHORE_PROJECT_PATH` env variable")
		projectPath := os.Getenv("SHORE_PROJECT_PATH")
		p.Log.Debug("`SHORE_PROJECT_PATH` set too ", projectPath)
		if projectPath == "" {
			return "", fmt.Errorf("env variable `SHORE_PROJECT_PATH` is not set")
		}

		return projectPath, nil
	}
	p.Log.Debug("Use `Getwd` (pwd) local path")

	projectPath, err := os.Getwd()

	if err != nil {
		p.Path = projectPath
	}

	return projectPath, err
}

// WriteFile write a file to the project path
func (p *Project) WriteFile(fileName, data string) error {
	// This method causes a Marshal->UnMarshal
	// ReMarshaling here makes the code easier to follow but hurts performance.
	p.Log.Debug("`WriteFile` was called with argument ", fileName)
	projectPath, err := p.GetProjectPath()

	if err != nil {
		return err
	}

	fulllFilepath := filepath.Join(projectPath, fileName)
	parentDir := filepath.Dir(fulllFilepath)
	// user (r/w/e), group (r/w), other (r/w)
	p.FS.MkdirAll(parentDir, 0766)

	return afero.WriteFile(p.FS, fulllFilepath, []byte(data), 0644)
}
