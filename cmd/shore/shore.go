package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/go-jsonnet"
	spingGate "github.com/spinnaker/spin/cmd/gateclient"
)

type CodeFile struct {
	Name string
	File []byte
}

func main() {
	path, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Current path %s\n", path)

	exePath, err := os.Executable()

	if err != nil {
		log.Println(err)
	}

	_, filename, _, ok := runtime.Caller(0)

	if ok == true {
		log.Printf("Caller runtime path %s", filename)
	}

	fmt.Printf("Executable path %s\n", exePath)

	pipelinesPath := filepath.Join(path, "pipelines")
	fmt.Printf("pipelines path %s \n", pipelinesPath)
	files, err := ioutil.ReadDir(pipelinesPath)

	var filesToRender []CodeFile

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if extension := filepath.Ext(f.Name()); extension != ".jsonnet" && extension != ".libsonnet" {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", pipelinesPath, f.Name())
		bytes, err := ioutil.ReadFile(filePath)

		if err != nil {
			log.Println(err)
		}

		filesToRender = append(filesToRender, CodeFile{Name: f.Name(), File: bytes})
	}

	var renderedOutput []string
	jsonnetVM := jsonnet.MakeVM()
	// jsonnetVM.Importer(jso)

	for _, jsonnetFile := range filesToRender {
		jsonRes, err := jsonnetVM.EvaluateSnippet(jsonnetFile.Name, string(jsonnetFile.File))

		if err != nil {
			log.Println(err)
		}

		renderedOutput = append(renderedOutput, jsonRes)
	}

	var pipelineInterface interface{}
	json.Unmarshal([]byte(renderedOutput[0]), &pipelineInterface)

	// homeDir, _ := os.UserHomeDir()
	// spinCliConfigPath, _ := GetUserConfig(filepath.Join(homeDir, ".spin", "config"))
	gateClient, err := spingGate.NewGateClient(nil, "", "", "", true)
	applications, _, err := gateClient.ApplicationControllerApi.GetAllApplicationsUsingGET(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", applications)

	fmt.Printf("%+v\n", pipelineInterface)
	res, err := gateClient.PipelineControllerApi.SavePipelineUsingPOST(context.Background(), pipelineInterface, nil)

	if err != nil {
		log.Printf("%+v\n", res)
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", res)
}
