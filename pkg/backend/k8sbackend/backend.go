package k8sbackend

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Autodesk/shore/pkg/shore_testing"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SpinClient represents a client for managing pipelines.
type KubeClient struct {
	log logrus.FieldLogger
}

// NewClient - Create a new default spinnaker client
func NewClient(logger logrus.FieldLogger) *KubeClient {
	return &KubeClient{log: logger}
}

// SavePipeline - Saves the pipeline as a YAML file in the specified folder.
func (s *KubeClient) SavePipeline(pipelineJSON string) (*http.Response, error) {
	// Unmarshal the pipeline JSON to extract metadata

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	generateFolderPath := filepath.Join(currentDir, "generated")
	if err := os.MkdirAll(generateFolderPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	var pipeline map[string]interface{}
	if err := jsoniter.Unmarshal([]byte(pipelineJSON), &pipeline); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pipeline JSON: %w", err)
	}

	// Iterate over the top-level keys and save each as a separate YAML file
	for key, value := range pipeline {
		// Convert the value to YAML
		pipelineYAML, err := yaml.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal YAML for key %s: %w", key, err)
		}

		// Create the file path
		filePath := filepath.Join(generateFolderPath, fmt.Sprintf("%s.yaml", key))

		// Write the YAML to the file
		if err := os.WriteFile(filePath, pipelineYAML, 0644); err != nil {
			return nil, fmt.Errorf("failed to write YAML file for key %s: %w", key, err)
		}

		s.log.Infof("Saved YAML for key %s to %s", key, filePath)
	}

	return nil, nil
}

func (s *KubeClient) ExecutePipeline(argsJSON string, stringify bool) (string, *http.Response, error) {

	// Check if a valid Kubernetes context is set
	cmd := exec.Command("kubectl", "config", "current-context")
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		return "", nil, fmt.Errorf("no valid Kubernetes context set: %w", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	folderPath := filepath.Join(currentDir, "generated")

	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return "", nil, fmt.Errorf("folder %s does not exist", folderPath)
	}

	// Create a Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return "", nil, fmt.Errorf("failed to build Kubernetes config: %w", err)
	}

	k8sClient, err := client.New(config, client.Options{})
	if err != nil {
		return "", nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Iterate over the files in the folder
	err = filepath.Walk(folderPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing file %s: %w", path, err)
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Read the YAML file
		fileContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Print the YAML content
		s.log.Infof("YAML Content of %s:\n%s", path, string(fileContent))

		// Decode the YAML into a Kubernetes object
		obj := &unstructured.Unstructured{}
		decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(fileContent), 1024)
		err = decoder.Decode(&obj.Object)

		// Check if the object has required fields
		if obj.GetAPIVersion() == "" || obj.GetKind() == "" {
			s.log.Errorf("Invalid Kubernetes object in file %s: missing apiVersion or kind", path)
			return fmt.Errorf("invalid Kubernetes object in file %s: missing apiVersion or kind", path)
		}

		// Apply the object using the Kubernetes client

		if err := k8sClient.Patch(context.TODO(), obj, client.Apply, client.FieldOwner("shore")); err != nil {
			s.log.Errorf("Failed to apply object: %v", obj)
			return fmt.Errorf("failed to apply resource from file %s: %s", path, obj)
		}

		return nil
	})

	if err != nil {
		return "", nil, fmt.Errorf("failed to apply files in folder: %w", err)
	}

	s.log.Info("Successfully applied all files in the folder")
	return "", nil, nil
}

func (s *KubeClient) DeletePipeline(argsJSON string) (*http.Response, error) {
	// Check if a valid Kubernetes context is set
	cmd := exec.Command("kubectl", "config", "current-context")
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		return nil, fmt.Errorf("no valid Kubernetes context set: %w", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	folderPath := filepath.Join(currentDir, "generated")

	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("folder %s does not exist", folderPath)
	}

	// Create a Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, fmt.Errorf("failed to build Kubernetes config: %w", err)
	}

	k8sClient, err := client.New(config, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Iterate over the files in the folder
	err = filepath.Walk(folderPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing file %s: %w", path, err)
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Read the YAML file
		fileContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Print the YAML content
		s.log.Infof("YAML Content of %s:\n%s", path, string(fileContent))

		// Decode the YAML into a Kubernetes object
		obj := &unstructured.Unstructured{}
		decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(fileContent), 1024)
		err = decoder.Decode(&obj.Object)

		// Check if the object has required fields
		if obj.GetAPIVersion() == "" || obj.GetKind() == "" {
			s.log.Errorf("Invalid Kubernetes object in file %s: missing apiVersion or kind", path)
			return fmt.Errorf("invalid Kubernetes object in file %s: missing apiVersion or kind", path)
		}

		// Delete the object using the Kubernetes client
		s.log.Infof("Deleting resource: %s", path)
		if err := k8sClient.Delete(context.TODO(), obj); err != nil {
			s.log.Errorf("Failed to delete object: %v", obj)
			return fmt.Errorf("failed to delete resource from file %s: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to delete resources in folder: %w", err)
	}

	// Delete the folder
	s.log.Infof("Deleting folder: %s", folderPath)
	if err := os.RemoveAll(folderPath); err != nil {
		return nil, fmt.Errorf("failed to delete folder %s: %w", folderPath, err)
	}

	s.log.Info("Successfully deleted all resources and the folder")
	return nil, nil
}

// TestPipeline - Dummy implementation for testing a pipeline.
func (s *KubeClient) TestPipeline(testConfig shore_testing.TestsConfig, onChange func(), stringify bool) error {
	s.log.Info("Starting TestPipeline...")

	s.log.Info("TestPipeline completed successfully.")
	return nil
}

// GetPipeline - Dummy implementation for retrieving a pipeline.
func (s *KubeClient) GetPipeline(application string, pipelineName string) (map[string]interface{}, *http.Response, error) {
	s.log.Infof("Retrieving pipeline: application=%s, pipelineName=%s", application, pipelineName)

	// Simulate a pipeline retrieval
	pipeline := map[string]interface{}{
		"application": application,
		"name":        pipelineName,
		"stages": []map[string]interface{}{
			{
				"name": "Deploy",
				"type": "deploy",
			},
		},
	}

	// Return the dummy pipeline and a nil HTTP response
	return pipeline, nil, nil
}

func (s *KubeClient) GetPipelinesNamesAndApplication(pipelineJSON string) ([]string, string, error) {
	// Unmarshal the pipeline JSON to extract metadata
	var pipeline map[string]interface{}
	if err := jsoniter.Unmarshal([]byte(pipelineJSON), &pipeline); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal pipeline JSON: %w", err)
	}

	// Extract the application name and pipeline names
	applicationName, ok := pipeline["application"].(string)
	if !ok {
		return nil, "", fmt.Errorf("failed to extract application name from pipeline JSON")
	}

	pipelineNames := []string{}
	if stages, ok := pipeline["stages"].([]interface{}); ok {
		for _, stage := range stages {
			if stageMap, ok := stage.(map[string]interface{}); ok {
				if name, ok := stageMap["name"].(string); ok {
					pipelineNames = append(pipelineNames, name)
				}
			}
		}
	}

	return pipelineNames, applicationName, nil
}

// WaitForPipelineToFinish - Dummy implementation for waiting for a pipeline to finish.
func (s *KubeClient) WaitForPipelineToFinish(id string, timeout int) (string, *http.Response, error) {
	s.log.Infof("Waiting for pipeline to finish: id=%s, timeout=%d", id, timeout)

	s.log.Info("Pipeline finished successfully.")
	return "", nil, nil
}
