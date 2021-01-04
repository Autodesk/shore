package jsonnet

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Autodesk/shore/pkg/renderer"
	"github.com/google/go-jsonnet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"golang.org/x/mod/modfile"
)

// MainFileName is the name of the entrypoint file the jsonnet renderer looks for to render a pipeline project
const MainFileName string = "main.pipeline.jsonnet"

// ArgsFileName is the name of the arguments file the jsonnet renderer looks for to pass to the pipeline as TLA veriables.
const ArgsFileName string = "render"

// Jsonnet - A Jsonnet renderer instance.
// The struct  holds the required parameters to render a standard shore pipeline.
type Jsonnet struct {
	renderer.Renderer
	VM  *jsonnet.VM
	FS  afero.Fs
	log logrus.FieldLogger
}

// NewRenderer - Create new instance of the JSONNET renderer.
func NewRenderer(fs afero.Fs, logger logrus.FieldLogger) *Jsonnet {
	return &Jsonnet{VM: jsonnet.MakeVM(), FS: fs, log: logger}
}

// Render - Render the code with the VM.
func (j *Jsonnet) Render(projectPath string, renderArgs string) (string, error) {
	renderFile := filepath.Join(projectPath, MainFileName)
	codeBytes, err := afero.ReadFile(j.FS, renderFile)

	if _, isPathErr := err.(*os.PathError); err != nil && !isPathErr {
		return "", err
	}

	fileImporter, err := j.getFileImporter(projectPath)

	if err != nil {
		return "", err
	}

	// Always include params, even if they are empty
	j.VM.TLACode("params", renderArgs)
	j.VM.Importer(fileImporter)
	// Currently adds the local `sponnet instance` to be available from `sponnet/*.libsonnet`
	return j.VM.EvaluateSnippet(renderFile, string(codeBytes))
}

func (j *Jsonnet) getFileImporter(projectPath string) (*jsonnet.FileImporter, error) {
	fileImporter := jsonnet.FileImporter{}

	sharedLibs, err := j.getSharedLibs(projectPath)

	if _, isPathErr := err.(*os.PathError); err != nil && !isPathErr {
		return &fileImporter, err
	}

	fileImporter.JPaths = append(fileImporter.JPaths, sharedLibs...)
	return &fileImporter, nil
}

func (j *Jsonnet) getSharedLibs(projectPath string) ([]string, error) {
	var sharedLibs []string
	var missingLibs = SharedLibsErr{}

	modFile, err := j.getGomodFile(projectPath)

	if err != nil {
		return []string{}, err
	}

	commonLibPath := filepath.Join(projectPath, "vendor")

	for _, requirement := range modFile.Require {
		libName := strings.Split(requirement.Mod.Path, "/")
		libNameSlice := libName[0 : len(libName)-1]
		libNameStr := strings.Join(libNameSlice, "/")

		libPath := filepath.Join(commonLibPath, libNameStr)
		err := j.libExists(requirement.Mod.Path, libPath)

		if err != nil {
			missingLibs.MissingLibs = append(missingLibs.MissingLibs, err.(SharedLibErr))
			continue
		}

		sharedLibs = append(sharedLibs, libPath)
	}

	if len(missingLibs.MissingLibs) > 0 {
		return sharedLibs, missingLibs
	}

	return sharedLibs, nil
}

func (j *Jsonnet) getGomodFile(projectPath string) (*modfile.File, error) {
	modFilePath := filepath.Join(projectPath, "go.mod")
	modFileBytes, err := afero.ReadFile(j.FS, modFilePath)

	if err != nil {
		return &modfile.File{}, err
	}

	modFile, err := modfile.Parse(modFilePath, modFileBytes, nil)

	if err != nil {
		return &modfile.File{}, err
	}

	return modFile, nil
}

func (j *Jsonnet) libExists(goPath, libPath string) error {
	exists, err := afero.DirExists(j.FS, libPath)

	// In case of a weird error OS error that isn't "NOT FOUND"
	if err != nil {
		return SharedLibErr{
			Require: goPath,
			Path:    libPath,
			Err:     err,
		}
	}

	if exists != true {
		return SharedLibErr{
			Require: goPath,
			Path:    libPath,
			Err:     &os.PathError{Op: "open", Path: libPath, Err: afero.ErrFileNotFound},
		}
	}

	return nil
}
