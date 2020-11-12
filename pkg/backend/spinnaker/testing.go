package spinnaker

import (
	"fmt"
	"reflect"
)

const (
	// PipelineNotStarted - Pipeline is in NOT_STARTED state
	PipelineNotStarted = "NOT_STARTED"
	// PipelineRunning - Pipeline is in RUNNING state
	PipelineRunning = "RUNNING"
	// PipelinePaused - Pipeline is in PAUSED state
	PipelinePaused = "PAUSED"
	// PipelineSuspended - Pipeline is in SUSPENDED state
	PipelineSuspended = "SUSPENDED"
	// PipelineSucceeded - Pipeline is in SUCCEEDED state
	PipelineSucceeded = "SUCCEEDED"
	// PipelineFailedContinue - Pipeline is in FAILED_CONTINUE state
	PipelineFailedContinue = "FAILED_CONTINUE"
	// PipelineTerminal - Pipeline is in TERMINAL state
	PipelineTerminal = "TERMINAL"
	// PipelineCanceled - Pipeline is in CANCELED state
	PipelineCanceled = "CANCELED"
	// PipelineRedirect - Pipeline is in REDIRECT state
	PipelineRedirect = "REDIRECT"
	// PipelineStopped - Pipeline is in STOPPED state
	PipelineStopped = "STOPPED"
	// PipelineSkipped - Pipeline is in SKIPPED state
	PipelineSkipped = "SKIPPED"
	// PipelineBuffered - Pipeline is in BUFFERED state
	PipelineBuffered = "BUFFERED"
)

const statusFailed = `
EXPECTED_STATUS failed assertion for stage '%s'
	expected: '%s'
	got: '%s'
`

const outputFailed = `
EXPECTED_OUTPUT failed assertion for stage '%s'
	expected: '%v'
	got: '%v'
`

// TestsConfig - describes the tests to run agains a pipeline
type TestsConfig struct {
	// RenderArgs
	Application string                `json:"application"`
	Pipeline    string                `json:"pipeline"`
	Tests       map[string]TestConfig `json:"tests"`
}

// TestConfig - describes a high level test config for a pipeline
type TestConfig struct {
	ExecArgs   map[string]interface{} `json:"execution_args"`
	Assertions map[string]Assertion   `json:"assertions"`
}

// Assertion - describes supported stage assertions
type Assertion struct {
	ExpectedStatus string                 `json:"expected_status"`
	ExpectedOutput map[string]interface{} `json:"expected_output"`
}

func isExpectedStatus(expectedStatus, status, stageName string) error {
	switch expectedStatus {
	// Test if this is one of the expected `states` for a `spinnaker stage` to be in.
	case PipelineNotStarted, PipelineRunning, PipelinePaused, PipelineSuspended, PipelineSucceeded, PipelineFailedContinue,
		PipelineTerminal, PipelineCanceled, PipelineRedirect, PipelineStopped, PipelineSkipped, PipelineBuffered:

		if expectedStatus != status {
			return fmt.Errorf(statusFailed, stageName, expectedStatus, status)
		}
	default:
		return fmt.Errorf("wrong status: '%s'", expectedStatus)
	}

	return nil
}

func isExpectedOutput(expectedOutput, outputs map[string]interface{}, stageName string) error {
	var errors []string
	localOutputs := make(map[string]interface{})

	for k := range expectedOutput {
		output, exists := outputs[k]

		if !exists {
			errors = append(errors, fmt.Sprintf("missing expected output key: '%s'", k))
			continue
		}

		localOutputs[k] = output
	}

	// Use deep euqal to validate the struct are 100% identical.
	if !reflect.DeepEqual(expectedOutput, localOutputs) {
		return fmt.Errorf(outputFailed, stageName, expectedOutput, localOutputs)
	}

	return nil
}
