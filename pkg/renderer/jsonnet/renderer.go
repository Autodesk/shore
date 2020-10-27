package jsonnet

import (
	"github.com/Autodesk/shore/pkg/renderer"
	"github.com/google/go-jsonnet"
)

type Jsonnet struct {
	renderer.Renderer
	VM *jsonnet.VM
}

func CreateJsonnetRenderer(sharedLibPaths ...string) *Jsonnet {
	jsonnetVM := jsonnet.MakeVM()
	fileImporter := jsonnet.FileImporter{}
	fileImporter.JPaths = append(fileImporter.JPaths, sharedLibPaths...)
	jsonnetVM.Importer(&fileImporter)

	return &Jsonnet{
		VM: jsonnetVM,
	}
}

func (j *Jsonnet) Render(filename string, code string) (string, error) {
	return j.VM.EvaluateSnippet(filename, code)
}
