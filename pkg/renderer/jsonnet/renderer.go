package jsonnet

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Autodesk/shore/pkg/fs"
	"github.com/Autodesk/shore/pkg/renderer"
	"github.com/google/go-jsonnet"
	"github.com/spf13/afero"
	"golang.org/x/mod/modfile"
)

var codeEntryPoint string = "main.pipeline.jsonnet"

type Jsonnet struct {
	renderer.Renderer
	VM          *jsonnet.VM
	projectPath string
	FS          afero.Fs
}

// NewRenderer - Create new instance of the JSONNET renderer.
func NewRenderer(projectPath string) (*Jsonnet, error) {
	fs := fs.GetFs()
	jsonnetVM := jsonnet.MakeVM()
	fileImporter := jsonnet.FileImporter{}
	sharedLibs, err := getSharedLibs(fs, projectPath)

	if _, isPathErr := err.(*os.PathError); err != nil && !isPathErr {
		return &Jsonnet{}, err
	}

	fileImporter.JPaths = append(fileImporter.JPaths, sharedLibs...)
	jsonnetVM.Importer(&fileImporter)

	return &Jsonnet{VM: jsonnetVM, projectPath: projectPath, FS: fs}, nil
}

// Render - Render the code with the VM.
func (j *Jsonnet) Render() (string, error) {
	mainCodeFile := filepath.Join(j.projectPath, codeEntryPoint)
	codeBytes, err := afero.ReadFile(j.FS, mainCodeFile)

	if err != nil {
		return "", err
	}

	// Currently adds the local `sponnet instance` to be available from `sponnet/*.libsonnet`
	return j.VM.EvaluateSnippet(mainCodeFile, string(codeBytes))
}

func getSharedLibs(fs afero.Fs, projectPath string) ([]string, error) {
	modFilePath := filepath.Join(projectPath, "go.mod")
	modFileBytes, err := afero.ReadFile(fs, modFilePath)

	if err != nil {
		return []string{}, err
	}

	modFile, err := modfile.Parse(modFilePath, modFileBytes, nil)

	if err != nil {
		return []string{}, err
	}

	var sharedLibs []string
	commonLibPath := filepath.Join(projectPath, "vendor")

	for _, requirement := range modFile.Require {
		libName := strings.Split(requirement.Mod.Path, "/")
		libNameSlice := libName[0 : len(libName)-1]
		libNameStr := strings.Join(libNameSlice, "/")

		fullLibPath := filepath.Join(commonLibPath, libNameStr)

		if exists, err := afero.DirExists(fs, fullLibPath); err != nil || exists != true {
			if err != nil {
				return []string{}, err
			}

			return []string{}, fmt.Errorf("Cannot find required library %s in library path %s", libName, commonLibPath)
		}

		sharedLibs = append(sharedLibs, fullLibPath)
	}

	return sharedLibs, nil
}
