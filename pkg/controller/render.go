package controller

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/Autodesk/shore/pkg/renderer/jsonnet"
)

type CodeFile struct {
	Name string
	File string
}

// Render - Using a defined renderer, renders a pipeline.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func Render(projectPath string) (string, error) {
	log.Printf("projectPath path %s \n", projectPath)
	mainCodeFile := filepath.Join(projectPath, "main.pipeline.jsonnet")
	codeBytes, err := ioutil.ReadFile(mainCodeFile)

	if err != nil {
		return "", err
	}

	jsonnetFile := CodeFile{Name: mainCodeFile, File: string(codeBytes)}

	// Currently adds the local `sponnet instance` to be available from `sponnet/*.libsonnet`
	sharedLibPath := filepath.Join(projectPath, "../", "../")
	renderer := jsonnet.NewRenderer(sharedLibPath)
	jsonRes, err := renderer.Render(jsonnetFile.Name, jsonnetFile.File)

	if err != nil {
		return "", err
	}

	return jsonRes, nil
}
