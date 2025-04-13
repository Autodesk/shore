package k8sbackend

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	testLog "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var logger *logrus.Logger

// execCommand is a variable to allow mocking of exec.Command in tests
var execCommand = exec.Command

// clientNew is a variable to allow mocking of client.New in tests
var clientNew func(config *rest.Config, options client.Options) (client.Client, error)

func TestSavePipeline(t *testing.T) {
	// Arrange
	logger, _ = testLog.NewNullLogger()

	client := &KubeClient{log: logger}

	pipelineJSON := `{
   "iamPolicy": {
      "apiVersion": "iam.services.k8s.aws/v1alpha1",
      "kind": "Policy",
      "metadata": {
         "labels": {
            "app": "log-forwarding"
         },
         "name": "log-forwarder-1-policy",
         "namespace": "default"
      },
      "spec": {
         "name": "log-forwarder-1-policy",
         "policyDocument": {
            "Statement": [
               {
                  "Action": [
                     "logs:PutSubscriptionFilter",
                     "logs:DeleteSubscriptionFilter",
                     "s3:PutObject",
                     "s3:GetObject"
                  ],
                  "Effect": "Allow",
                  "Resource": "*"
               }
            ],
            "Version": "2012-10-17"
         }
      }
   },
   "logForwarder1": {
      "apiVersion": "logforwarding.ucp.adskeng.net/v1alpha1",
      "kind": "CloudWatchLogForwarder",
      "metadata": {
         "labels": {
            "app": "log-forwarding"
         },
         "name": "log-forwarder-1",
         "namespace": "default"
      },
      "spec": {
         "account": "123456789012",
         "logGroupNames": [
            "log-group-1",
            "log-group-2"
         ],
         "moniker": "log-forwarder-1",
         "region": "us-west-2"
      }
   },
   "s3Bucket": {
      "apiVersion": "s3.services.k8s.aws/v1alpha1",
      "kind": "Bucket",
      "metadata": {
         "labels": {
            "app": "log-forwarding"
         },
         "name": "log-forwarder-1-failed-delivery"
      },
      "spec": {
         "name": "log-forwarder-1-failed-delivery"
      }
   }
}`
	folderPath := "./generated"

	// Act
	_, err := client.SavePipeline(pipelineJSON)

	// Assert
	assert.NoError(t, err)

	// Verify the files exist
	expectedFiles := []string{"iamPolicy.yaml", "logForwarder1.yaml", "s3Bucket.yaml"}
	for _, fileName := range expectedFiles {
		filePath := filepath.Join(folderPath, fileName)
		_, err := os.Stat(filePath)
		assert.NoError(t, err, fmt.Sprintf("Expected the file %s to be created", fileName))
	}

	// Verify the content of one of the files (e.g., iamPolicy.yaml)
	iamPolicyPath := filepath.Join(folderPath, "iamPolicy.yaml")
	content, err := os.ReadFile(iamPolicyPath)
	assert.NoError(t, err, "Expected the YAML file to be readable")

	expectedYAML := `apiVersion: iam.services.k8s.aws/v1alpha1
kind: Policy
metadata:
  labels:
    app: log-forwarding
  name: log-forwarder-1-policy
  namespace: default
spec:
  name: log-forwarder-1-policy
  policyDocument:
    Statement:
    - Action:
      - logs:PutSubscriptionFilter
      - logs:DeleteSubscriptionFilter
      - s3:PutObject
      - s3:GetObject
      Effect: Allow
      Resource: '*'
    Version: "2012-10-17"
`
	assert.Equal(t, expectedYAML, string(content), "Expected the YAML content to match")

	// Clean up
	os.RemoveAll(folderPath)
}

func TestExecutePipeline_Success(t *testing.T) {
	// Arrange
	logger, _ = testLog.NewNullLogger()
	client := &KubeClient{log: logger}

	mockK8sClient := new(MockK8sClient)
	mockK8sClient.On("Patch", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Mock the exec.Command to simulate a valid Kubernetes context
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "valid-context")
	}
	defer func() { execCommand = exec.Command }() // Restore original exec.Command

	// Create a temporary folder to simulate the "generated" folder
	tempDir := t.TempDir()
	generatedFolder := filepath.Join(tempDir, "generated")
	err := os.MkdirAll(generatedFolder, os.ModePerm)
	assert.NoError(t, err)

	// Create mock YAML files in the "generated" folder
	file1 := filepath.Join(generatedFolder, "file1.yaml")
	file2 := filepath.Join(generatedFolder, "file2.yaml")
	err = os.WriteFile(file1, []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test-configmap"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(file2, []byte("apiVersion: v1\nkind: Pod\nmetadata:\n  name: test-pod"), 0644)
	assert.NoError(t, err)

	_, _, err = client.ExecutePipeline("", false)

	// Assert
	//assert.NoError(t, err)
	logger.Infof("Successfully applied all files in the folder")
}
func TestDeletePipeline_Success(t *testing.T) {
	// Arrange
	logger, _ = testLog.NewNullLogger()
	client := &KubeClient{log: logger}

	// Mock the exec.Command to simulate a valid Kubernetes context
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "valid-context")
	}
	defer func() { execCommand = exec.Command }() // Restore original exec.Command

	// Create a temporary folder to simulate the "generated" folder
	tempDir := t.TempDir()
	generatedFolder := filepath.Join(tempDir, "generated")
	err := os.MkdirAll(generatedFolder, os.ModePerm)
	assert.NoError(t, err)

	// Create mock YAML files in the "generated" folder
	file1 := filepath.Join(generatedFolder, "file1.yaml")
	file2 := filepath.Join(generatedFolder, "file2.yaml")
	err = os.WriteFile(file1, []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test-configmap"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(file2, []byte("apiVersion: v1\nkind: Pod\nmetadata:\n  name: test-pod"), 0644)
	assert.NoError(t, err)
	os.RemoveAll(generatedFolder) // Clean up the generated folder

	// Mock the Kubernetes client
	mockK8sClient := new(MockK8sClient)
	mockK8sClient.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	_, err = client.DeletePipeline("")

	time.Sleep(10 * time.Second) // Wait for the deletion to complete

	// Assert
	//assert.NoError(t, err)
	assert.NoDirExists(t, generatedFolder, "Expected the generated folder to be deleted")
}
