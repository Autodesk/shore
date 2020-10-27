package controller

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Autodesk/shore/pkg/renderer/jsonnet"
)

type CodeFile struct {
	Name string
	File []byte
}

// Render - Using a defined renderer, renders a pipeline.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func Render() ([]string, error) {
	path, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	pipelinesPath := filepath.Join(path, "pipelines")
	log.Printf("pipelines path %s \n", pipelinesPath)
	files, err := ioutil.ReadDir(pipelinesPath)

	if err != nil {
		return []string{}, err
	}

	var filesToRender []CodeFile

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if extension := filepath.Ext(f.Name()); extension != ".jsonnet" && extension != ".libsonnet" {
			continue
		}

		filePath := filepath.Join(pipelinesPath, f.Name())
		bytes, err := ioutil.ReadFile(filePath)

		if err != nil {
			return []string{}, err
		}

		filesToRender = append(filesToRender, CodeFile{Name: f.Name(), File: bytes})
	}

	// Currently adds the local `sponnet instance` to be available from `sponnet/*.libsonnet`
	renderer := jsonnet.CreateJsonnetRenderer(
		filepath.Join(path, "../", "../"),
	)

	var renderedOutput []string

	for _, jsonnetFile := range filesToRender {
		jsonRes, err := renderer.Render(jsonnetFile.Name, string(jsonnetFile.File))

		if err != nil {
			return []string{}, err
		}

		renderedOutput = append(renderedOutput, jsonRes)
	}

	return renderedOutput, nil
}
