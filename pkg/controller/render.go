package controller

import (
	"github.com/Autodesk/shore/pkg/renderer/jsonnet"
)

// Render - Using a defined renderer, renders a pipeline.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func Render(projectPath string) (string, error) {
	// TODO: The renderer should either be DI'ed or imported from a "global" context.
	// We don't know which renderer the customer may choose to use in the future.
	renderer, err := jsonnet.NewRenderer(projectPath)

	if err != nil {
		return "", err
	}

	pipelineJSON, err := renderer.Render()

	if err != nil {
		return "", err
	}

	return pipelineJSON, nil
}
