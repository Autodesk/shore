package project_test

import (
	"testing"

	"github.com/Autodeskshore/pkg/project"
	log "github.com/sirupsen/logrus"
	testLog "github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var Logger *log.Logger

func init() {
	Logger, _ = testLog.NewNullLogger()
	Logger.SetLevel(log.DebugLevel)
}

func TestIsValidGoVersion(t *testing.T) {
	// Given
	goVersionStr := "go version go1.16.0 darwin/amd64"

	// Test
	res, err := project.IsValidGoVersion(goVersionStr)

	// Assert
	assert.Equal(t, true, res)
	assert.Nil(t, err)
}

func TestIsValidGoVersionNotValid(t *testing.T) {
	// Given
	goVersionStr := "go version go1.10.10 darwin/amd64"

	// Test
	res, err := project.IsValidGoVersion(goVersionStr)

	// Assert
	assert.Equal(t, false, res)
	assert.Error(t, err)
}

func TestIsValidGoVersionBrokenString(t *testing.T) {
	// Given
	goVersionStr := "goversiongo1.0.10amd64"

	// Test
	res, err := project.IsValidGoVersion(goVersionStr)

	// Assert
	assert.Equal(t, false, res)
	assert.Error(t, err)
}

func TestIsValidGoVersionNoVersion(t *testing.T) {
	// Given
	goVersionStr := "goversiongoamd64"

	// Test
	res, err := project.IsValidGoVersion(goVersionStr)

	// Assert
	assert.Equal(t, false, res)
	assert.Error(t, err)
}

func TestIsValidGoVersionReallyBrokenVersion(t *testing.T) {
	// Given
	goVersionStr := "1.11 1.12 1.13 1.1.1.1.1 10.10.10.10 13.13.13.13"

	// Test
	res, err := project.IsValidGoVersion(goVersionStr)

	// Assert
	assert.Equal(t, false, res)
	assert.Error(t, err)
}

func TestInitSuccess(t *testing.T) {
	// Given
	localFs := afero.NewMemMapFs()

	init := project.ShoreProjectInit{
		ProjectName: "my-project",
		Renderer:    "jsonnet",
		Backend:     "spinnaker",
		Libraries:   []string{"https://github.com/Autodeskspin-lib-jsonnet.git", "https://github.com/Autodeskadsk-lib-jsonnet.git"},
	}

	pInit := project.ProjectInitialize{
		Log: Logger,
		Project: project.Project{
			FS:   localFs,
			Log:  Logger,
			Path: "/tmp/test/",
		},
	}

	// Test
	err := pInit.Init(init)

	// Assert
	mainExists, _ := afero.Exists(localFs, "/tmp/test/main.pipeline.jsonnet")
	jsonnetfileExists, _ := afero.Exists(localFs, "/tmp/test/jsonnetfile.json")
	gitIgnoreExists, _ := afero.Exists(localFs, "/tmp/test/.gitignore")
	readmeExists, _ := afero.Exists(localFs, "/tmp/test/README.md")
	renderExists, _ := afero.Exists(localFs, "/tmp/test/render.yml")
	execExists, _ := afero.Exists(localFs, "/tmp/test/exec.yml")
	e2eExists, _ := afero.Exists(localFs, "/tmp/test/E2E.yml")
	testExampleExists, _ := afero.Exists(localFs, "/tmp/test/tests/example_test.libsonnet")

	assert.True(t, mainExists)
	assert.True(t, jsonnetfileExists)
	assert.True(t, gitIgnoreExists)
	assert.True(t, readmeExists)
	assert.True(t, renderExists)
	assert.True(t, execExists)
	assert.True(t, e2eExists)
	assert.True(t, testExampleExists)
	assert.Nil(t, err)
}

func TestShortName(t *testing.T) {
	// Given
	init := project.ShoreProjectInit{
		ProjectName: "my-project",
		Renderer:    "jsonnet",
		Backend:     "spinnaker",
		Libraries:   []string{"https://github.com/Autodeskspin-lib-jsonnet.git", "https://github.com/Autodeskadsk-lib-jsonnet.git"},
	}
	// Tests
	shortName := init.ShortName()
	// Assert
	assert.Equal(t, "myproject", shortName)
}
