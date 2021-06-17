package jsonnet

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Autodeskshore/pkg/renderer"
	"github.com/google/go-jsonnet"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	jbV1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
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
	vm           *jsonnet.VM
	fs           afero.Fs
	log          logrus.FieldLogger
	fileImporter jsonnet.Importer
}

// NewRenderer - Create new instance of the JSONNET renderer.
func NewRenderer(fs afero.Fs, logger logrus.FieldLogger) *Jsonnet {
	return &Jsonnet{vm: jsonnet.MakeVM(), fs: fs, log: logger}
}

// Render - Render the code with the VM.
func (j *Jsonnet) Render(projectPath string, renderArgs string, renderType renderer.RenderType) (string, error) {
	renderFile := filepath.Join(projectPath, RenderFiles[renderType])

	// TODO implement lazy loading
	codeBytes, err := afero.ReadFile(j.fs, renderFile)

	if err != nil {
		return "", err
	}

	jbFile, err := j.loadJsonnetBundlerFile(projectPath)

	// If the file doesn't exist, we can skip the error.
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	j.fileImporter, err = GetFileImporter(projectPath, jbFile)

	if err != nil {
		return "", err
	}

	// Always include params, even if they are empty
	j.vm.TLACode("params", renderArgs)
	j.vm.Importer(j.fileImporter)
	// Currently adds the local `sponnet instance` to be available from `sponnet/*.libsonnet`
	return j.vm.EvaluateSnippet(renderFile, string(codeBytes))
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

// GetFileImporter - Get the Jsonnet File Import customized to the Jsonnet Bundler type.
func GetFileImporter(projectPath string, jbFile jbV1.JsonnetFile) (*jsonnet.FileImporter, error) {
	libsPath := []string{}
	fileImporter := jsonnet.FileImporter{}

	if jbFile.LegacyImports {
		// Jsonnet-Bundler LegacyImports put the imported folders in the top directory with symlinks.
		libPath := filepath.Join(projectPath, ShareLibsPath)
		libsPath = append(libsPath, libPath)
	} else {
		// Jsonnet-Bundler Imports put complies to the GoMod style of artifact management (vendoring)
		// This means we need to take an extra step to find the top level key for each shared folder.
		libsMap := make(map[string][]string)

		for k := range jbFile.Dependencies {
			libPath := filepath.Join(projectPath, ShareLibsPath, k)
			libPathSplit := strings.Split(libPath, "/")
			libsKey := strings.Join(libPathSplit[:len(libPathSplit)-1], "/")

			if len(libsMap[libsKey]) > 0 {
				libsMap[libsKey] = append(libsMap[libsKey], k)
			} else {
				libsMap[libsKey] = []string{k}
			}
		}

		for k := range libsMap {
			libsPath = append(libsPath, k)
		}
	}

	fileImporter.JPaths = append(fileImporter.JPaths, libsPath...)
	return &fileImporter, nil
}
