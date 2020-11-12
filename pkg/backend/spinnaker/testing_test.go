package spinnaker

// Funniest filename to ever exists!
// Testing_test validates that the remote-test function is working properly.
// Mainly, the validation logic.

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	// Give
	// Test
	// Assert
	assert.EqualError(t, errors.New(""), "")
}

func TestIsExpectedStatusSuccess(t *testing.T) {
	// Test
	err := isExpectedStatus("TERMINAL", "TERMINAL", "myStage")

	assert.Nil(t, err)
}

func TestIsExpectedStatusFailure(t *testing.T) {
	// Test
	err := isExpectedStatus("FATAILITY", "FATAILITY", "myStage")

	assert.EqualError(t, err, "wrong status: 'FATAILITY'")
}

func TestIsExpectedOutputSuccess(t *testing.T) {
	expectedOutput := map[string]interface{}{
		"This": "test",
		"is":   []string{"Ok!"},
	}

	output := expectedOutput

	err := isExpectedOutput(expectedOutput, output, "myStage")
	assert.Nil(t, err)
}

func TestIsExpectedOutputFailure(t *testing.T) {
	expectedOutput := map[string]interface{}{
		"This": "test",
		"is":   []string{"Ok!"},
	}

	output := map[string]interface{}{
		"This": "test",
		"is":   []string{"Not Ok!"},
	}

	expectedError := fmt.Sprintf(outputFailed, "myStage", expectedOutput, output)

	err := isExpectedOutput(expectedOutput, output, "myStage")
	assert.EqualError(t, err, expectedError)
}
