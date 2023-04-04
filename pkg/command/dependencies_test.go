package command

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/Autodesk/shore/pkg/backend/spinnaker"
	"github.com/Autodesk/shore/pkg/project"
	"github.com/Autodesk/shore/pkg/renderer/jsonnet"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var testPath = "/test"
var shoreConfigTemplate = `{
	"renderer": {
		"type": "%s"
	},
	"executor": {
		"type": "%s",
		"config": {
		"default": "~/.spin/sb-config",
		"prodSpin": "~/.spin/prod-config"
		}
	},
	"profiles": {
		"default": {
		"application": "test1test2test3",
		"pipeline": "simple-pipeline-test",
		"render": "render.yaml",
		"exec": "exec.yaml",
		"e2e": "e2e.yaml"
		}
	}
}`

func SetupTestProject() *project.Project {
	os.Setenv("LOCAL", "true")
	os.Setenv("SHORE_PROJECT_PATH", testPath)

	memFs := afero.NewMemMapFs()
	memFs.Mkdir(testPath, os.ModePerm)

	logger, _ := test.NewNullLogger()

	return project.NewShoreProject(memFs, logger)
}

func TestPassingNewDependencies(t *testing.T) {
	// Given
	proj := SetupTestProject()

	tests := []struct {
		name               string
		configuredBackend  string
		configuredRenderer string
		expectedBackend    interface{}
		expectedRenderer   interface{}
	}{
		{
			name:               "spinnaker/jsonnet",
			configuredBackend:  "spinnaker",
			configuredRenderer: "jsonnet",
			expectedBackend:    &spinnaker.SpinClient{},
			expectedRenderer:   &jsonnet.Jsonnet{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shoreConfig := fmt.Sprintf(shoreConfigTemplate, test.configuredRenderer, test.configuredBackend)
			afero.WriteFile(proj.FS, path.Join(testPath, "shore.json"), []byte(shoreConfig), os.ModePerm)

			// When
			deps, err := NewDependencies(proj)

			// Then
			assert.NoError(t, err)
			assert.IsType(t, test.expectedRenderer, deps.Renderer)
			assert.IsType(t, test.expectedBackend, deps.Backend)
		})
	}
}

func TestFailingNewDependencies(t *testing.T) {
	// Given
	proj := SetupTestProject()

	tests := []struct {
		name               string
		configuredBackend  string
		configuredRenderer string
		expectedError      string
	}{
		{
			name:               "wrong-backend",
			configuredBackend:  "yolo",
			configuredRenderer: "jsonnet",
			expectedError:      "Backend is undefined",
		},
		{
			name:               "wrong-renderer",
			configuredBackend:  "spinnaker",
			configuredRenderer: "yolo",
			expectedError:      "Renderer is undefined",
		},
		{
			name:               "malformed-config",
			configuredBackend:  "\"",
			configuredRenderer: "yolo",
			expectedError:      "object not ended with",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shoreConfig := fmt.Sprintf(shoreConfigTemplate, test.configuredRenderer, test.configuredBackend)
			afero.WriteFile(proj.FS, path.Join(testPath, "shore.json"), []byte(shoreConfig), os.ModePerm)

			// When
			deps, err := NewDependencies(proj)

			// Then
			assert.Empty(t, deps)
			assert.Error(t, err)
			assert.ErrorContains(t, err, test.expectedError)
		})
	}
}
