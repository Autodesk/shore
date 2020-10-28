package jsonnet

import (
	"github.com/Autodesk/shore/pkg/renderer"
	"github.com/google/go-jsonnet"
)

type Jsonnet struct {
	renderer.Renderer
	VM *jsonnet.VM
}

// NewRenderer - Create new instance of the JSONNET renderer.
func NewRenderer(sharedLibPaths ...string) *Jsonnet {
	fileImporter := jsonnet.FileImporter{}
	fileImporter.JPaths = append(fileImporter.JPaths, sharedLibPaths...)

	jsonnetVM := jsonnet.MakeVM()
	jsonnetVM.Importer(&fileImporter)

	return &Jsonnet{
		VM: jsonnetVM,
	}
}

func (j *Jsonnet) Render(filename string, code string) (string, error) {
	return j.VM.EvaluateSnippet(filename, code)
}
