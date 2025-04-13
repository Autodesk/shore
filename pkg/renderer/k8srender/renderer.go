package k8srender

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Autodesk/shore/pkg/renderer"
	"github.com/google/go-jsonnet"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	jbV1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

var RenderFiles = map[renderer.RenderType]string{
	// MainFileName is the name of the entrypoint file the jsonnet renderer looks for to render a pipeline project
	renderer.MainFileName:    "main.pipeline.jsonnet",
	renderer.CleanUpFileName: "cleanup/cleanup.pipeline.jsonnet",
}

// ArgsFileName is the name of the arguments file the jsonnet renderer looks for to pass to the pipeline as TLA veriables.
const ArgsFileName string = "render"

// ShareLibsPath - the path that jsonnet should look into when looking for shared libraries.
//
// Designed to work with - https://github.com/jsonnet-bundler/jsonnet-bundler/
const ShareLibsPath string = "vendor"

const JsonnetFileName string = "jsonnetfile.json"

// Jsonnet - A Jsonnet renderer instance.
// The struct  holds the required parameters to render a standard shore pipeline.
type Jsonnet struct {
	renderer.Renderer
	vm  *jsonnet.VM
	fs  afero.Fs
	log logrus.FieldLogger
}

// NewRenderer - Create new instance of the JSONNET renderer.
func NewRenderer(fs afero.Fs, logger logrus.FieldLogger) *Jsonnet {
	return &Jsonnet{vm: jsonnet.MakeVM(), fs: fs, log: logger}
}

// Render - Render the code with the VM.
func (j *Jsonnet) Render(projectPath string, renderArgs string, renderType renderer.RenderType) (string, error) {
	renderFile := filepath.Join(projectPath, RenderFiles[renderType])

	jbFile, err := j.loadJsonnetBundlerFile(projectPath)

	// If the file doesn't exist, we can skip the error.
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	// Always include params, even if they are empty
	j.vm.TLACode("params", renderArgs)
	j.vm.Importer(NewImporter(j.fs, projectPath, jbFile))

	// Evaluate the Jsonnet file to get JSON output
	jsonOutput, err := j.vm.EvaluateFile(renderFile)
	if err != nil {
		return "", err
	}

	// Convert JSON output to YAML
	var jsonData map[string]interface{}
	// Unmarshal the JSON string into a Go map
	if err := json.Unmarshal([]byte(jsonOutput), &jsonData); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Marshal the Go map into YAML
	yamlOutput, err := yaml.Marshal(jsonData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return string(yamlOutput), nil
}

// A compliant wrapper implementing jsonnetfile.Load but using `Afero` instrad of `ioutil`.
func (j *Jsonnet) loadJsonnetBundlerFile(path string) (jbV1.JsonnetFile, error) {
	jsonnetFilePath := filepath.Join(path, JsonnetFileName)
	bytes, err := afero.ReadFile(j.fs, jsonnetFilePath)
	if err != nil {
		return jbV1.New(), err
	}

	return jsonnetfile.Unmarshal(bytes)
}
