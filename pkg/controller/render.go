package controller

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

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
	modFilePath := filepath.Join(projectPath, "go.mod")
	modFileBytes, err := ioutil.ReadFile(modFilePath)

	if err != nil {
		return "", err
	}

	modFile, err := modfile.Parse(modFilePath, modFileBytes, nil)

	if err != nil {
		return "", err
	}

	var sharedLibs []string

	for _, requirement := range modFile.Require {
		libPath := strings.Split(requirement.Mod.Path, "/")
		libPathSlice := libPath[0 : len(libPath)-1]
		libPathStr := strings.Join(libPathSlice, "/")

		fullLibPath := filepath.Join(projectPath, "vendor", libPathStr)
		sharedLibs = append(sharedLibs, fullLibPath)
	}

	renderer := jsonnet.NewRenderer(sharedLibs...)

	mainCodeFile := filepath.Join(projectPath, "main.pipeline.jsonnet")
	codeBytes, err := ioutil.ReadFile(mainCodeFile)

	if err != nil {
		return "", err
	}

	jsonnetFile := CodeFile{Name: mainCodeFile, File: string(codeBytes)}
	// Currently adds the local `sponnet instance` to be available from `sponnet/*.libsonnet`
	jsonRes, err := renderer.Render(jsonnetFile.Name, jsonnetFile.File)

	if err != nil {
		return "", err
	}

	return jsonRes, nil
}
