package spinnaker

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	testLog "github.com/sirupsen/logrus/hooks/test"
	spinGateApi "github.com/spinnaker/spin/gateapi"
	"github.com/stretchr/testify/assert"
)

type mockApplicationControllerAPI struct{}

func (a *mockApplicationControllerAPI) GetPipelineConfigUsingGET(ctx context.Context, application string, pipelineName string) (map[string]interface{}, *http.Response, error) {
	var res map[string]interface{}

	if application == "not-exists" {
		res = map[string]interface{}{}
	} else {
		res = map[string]interface{}{
			"id": "1234",
		}
	}

	return res, &http.Response{StatusCode: http.StatusOK}, nil
}

type mockApplicationControllerAPIWithEmptyID struct{}

func (a *mockApplicationControllerAPIWithEmptyID) GetPipelineConfigUsingGET(ctx context.Context, application string, pipelineName string) (map[string]interface{}, *http.Response, error) {
	res := map[string]interface{}{}

	return res, &http.Response{StatusCode: http.StatusOK}, nil
}

type mockPipelineControllerAPI struct{}

func (p *mockPipelineControllerAPI) SavePipelineUsingPOST(ctx context.Context, pipeline interface{}, localVarOptionals *spinGateApi.PipelineControllerApiSavePipelineUsingPOSTOpts) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK}, nil
}

func (p *mockPipelineControllerAPI) InvokePipelineConfigUsingPOST1(ctx context.Context, application string, pipelineNameOrID string, localVarOptionals *spinGateApi.PipelineControllerApiInvokePipelineConfigUsingPOST1Opts) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK}, nil
}

type mockCustomSpinCli struct {
	CustomSpinCLI
}

func (s *mockCustomSpinCli) ExecutePipeline(application, pipelineName string, args *bytes.Buffer) (*ExecutePipelineResponse, *http.Response, error) {
	return &ExecutePipelineResponse{Ref: "/pipeline/1234"}, &http.Response{StatusCode: http.StatusOK}, nil
}

func (s *mockCustomSpinCli) PipelineExecutionDetails(refID string, args *bytes.Buffer) (*PipelineExecutionDetailsResponse, *http.Response, error) {
	return &PipelineExecutionDetailsResponse{
		Application: "application",
		Stages: []map[string]interface{}{
			{
				"name":   "name",
				"status": "SUCCEEDED",
				"outputs": map[string]interface{}{
					"test": "123",
				},
			},
		},
		PipelineName: "pipeline",
	}, &http.Response{StatusCode: http.StatusOK}, nil
}

var cli *SpinClient

var logger *logrus.Logger

func init() {
	logger, _ = testLog.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	cli = &SpinClient{
		log:           logger,
		CustomSpinCLI: &mockCustomSpinCli{},
		SpinCLI: &SpinCLI{
			ApplicationControllerAPI: &mockApplicationControllerAPI{},
			PipelineControllerAPI:    &mockPipelineControllerAPI{},
			Context:                  context.Background(),
		},
	}
}

func TestInternalSaveSuccessForExistingPipeline(t *testing.T) {
	// Test
	pipelineID, res, err := cli.savePipeline(`{"application": "test", "name": "test"}`)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, pipelineID, "1234")
}

func TestInternalSaveSuccessForNonExistingPipeline(t *testing.T) {
	// Test
	pipelineID, res, err := cli.savePipeline(`{"application": "not-exists", "name": "test"}`)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, pipelineID, "")
}

func TestInternalSaveFailedApplication(t *testing.T) {
	// Test
	_, _, err := cli.savePipeline(`{"name": "test"}`)

	// Assert
	assert.EqualError(t, err, "required pipeline key 'application' missing")
}

func TestInternalSaveFailedName(t *testing.T) {

	// Test
	_, _, err := cli.savePipeline(`{"application": "test"}`)

	// Assert
	assert.EqualError(t, err, "required pipeline key 'name' missing")
}

func TestExecSuccess(t *testing.T) {
	// Test
	_, res, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "parameters": {"answer": 42}}`)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestExecFailedApplication(t *testing.T) {
	// Test
	_, _, err := cli.ExecutePipeline(`{"pipeline": "test"}`)

	// Assert
	assert.EqualError(t, err, "required args key 'application' missing")
}

func TestExecFailedName(t *testing.T) {
	// Test
	_, _, err := cli.ExecutePipeline(`{"application": "test"}`)

	// Assert
	assert.EqualError(t, err, "required args key 'pipeline' missing")
}

func TestExecArgsFailing(t *testing.T) {
	// Test
	_, _, err := cli.ExecutePipeline(`{}`)

	// Assert
	assert.EqualError(t, err, "required args key 'pipeline' missing\nrequired args key 'application' missing")
}

func TestSaveSuccess(t *testing.T) {
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"application": "appname",
							"name":  "Child pipeline stage",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2",
								"stages": [
									{
										"name": "Child pipeline 2 stage",
										"application": "appname"
									}
								]
							}
						},
						{
							"application": "appname",
							"name":  "Child pipeline stage 2",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2.2",
								"stages": [
									{
										"name": "Child pipeline 2.2 stage",
										"application": "appname"
									}
								]
							}
						}
					]
				}
			}
		]
	}
	`

	// Test
	res, err := cli.SavePipeline(nestedPipelineString)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestSaveSimplePipelineSuccess(t *testing.T) {
	// Given
	nestedPipelineString := `
	{
		"application": "test123",
		"name": "simple pipeline",
		"stages": [
			 {
					"name": "Wait",
					"waitTime": 1
			 }
		]
 }
	`

	// Test
	res, err := cli.SavePipeline(nestedPipelineString)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestMissingApplicationFailedSave(t *testing.T) {
	// Given
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"pipeline": {
					"name": "Child pipeline 1",
					"stages": [
						{
							"name": "Child pipeline 2 stage",
							"application": "appname"
						}
					]
				}
			}
		]
	}
	`

	// Test
	_, err := cli.SavePipeline(nestedPipelineString)

	// Assert
	assert.EqualError(t, err, "required pipeline key 'application' missing")
}

func TestMissingNameFailedSave(t *testing.T) {
	// Given
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"pipeline": {
					"application": "appname",
					"stages": [
						{
							"name": "Child pipeline 2 stage",
							"application": "appname"
						}
					]
				}
			}
		]
	}
	`

	// Test
	_, err := cli.SavePipeline(nestedPipelineString)

	// Assert
	assert.EqualError(t, err, "required pipeline key 'name' missing")
}

func TestPipelineChildPipelineWrongApplicationFailedSave(t *testing.T) {
	// Given
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"pipeline": {
					"application": "another appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"name": "child pipeline 1 stage"
						}
					]
				}
			}
		]
	}
	`

	// Test
	_, err := cli.SavePipeline(nestedPipelineString)

	// Assert
	assert.EqualError(t, err, "pipeline 'application' key value should match the value of parent pipeline 'application' key")
}

func TestPipelineStageMissingApplicationFailedSave(t *testing.T) {
	// Given
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"name":  "Nested pipeline stage",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"name": "child pipeline 1 stage"
						}
					]
				}
			}
		]
	}
	`

	// Test
	_, err := cli.SavePipeline(nestedPipelineString)

	// Assert
	assert.EqualError(t, err, "required stage key 'application' missing for stage")
}

func TestPipelineStageWrongApplicationFailedSave(t *testing.T) {
	// Given
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "another appname",
				"name":  "Nested pipeline stage",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"name": "child pipeline 1 stage"
						}
					]
				}
			}
		]
	}
	`

	// Test
	_, err := cli.SavePipeline(nestedPipelineString)

	// Assert
	assert.EqualError(t, err, "'application' key value of stage of type 'pipeline' should match the one of parent pipeline 'application' value")
}

func TestTestingRemoteSuccess(t *testing.T) {
	config := `{
	"application": "test1test2test3",
	"pipeline": "abc",
	"tests": {
		"Test Success": {
			"execution_args": {
				"parameters": {
				"a": "a"
				}
			},
			"assertions": {
				"name": {
					"expected_status": "succeeded",
					"expected_output": {
						"test": "123"
					}
				}
			}
		}
	}
}
`

	err := cli.TestPipeline(config, func() {})

	assert.Nil(t, err)
}

func TestTestingRemoteNoAssertionFailed(t *testing.T) {
	config := `{
	"application": "test1test2test3",
	"pipeline": "abc",
	"tests": {
		"Test Success": {
			"execution_args": {
				"parameters": {
				"a": "a"
				}
			},
			"assertions": {
				"name": {
					"expected_status": "succeeded",
					"expected_output": {
						"test": "1234"
					}
				}
			}
		}
	}
}
`

	err := cli.TestPipeline(config, func() {})

	assert.Error(t, err)
}

func TestTestingRemoteNoAssertionForStageError(t *testing.T) {
	config := `{
	"application": "test1test2test3",
	"pipeline": "abc",
	"tests": {
		"Test Success": {
			"execution_args": {},
			"assertions": {}
		}
	}
}
`

	err := cli.TestPipeline(config, func() {})

	assert.Error(t, err)
}

func TestTestingRemoteMissingExecArgs(t *testing.T) {
	config := `{
	"application": "test1test2test3",
	"pipeline": "abc",
	"tests": {
		"Test Success": {
			"assertions": {
				"name": {
					"expected_status": "succeeded",
					"expected_output": {
						"test": "123"
					}
				}
			}
		}
	}
}
`

	err := cli.TestPipeline(config, func() {})

	assert.Nil(t, err)
}

func TestTestingNoApplicationFailed(t *testing.T) {
	config := `{
	"application": "test1test2test3",
	"tests": {
		"Test Success": {
			"assertions": {
				"name": {
					"expected_status": "succeeded",
					"expected_output": {
						"test": "123"
					}
				}
			}
		}
	}
}
`

	err := cli.TestPipeline(config, func() {})

	assert.Error(t, err)
}

func TestTestingNoPipelineFailed(t *testing.T) {
	config := `{
	"pipeline": "abc",
	"tests": {
		"Test Success": {
			"assertions": {
				"name": {
					"expected_status": "succeeded",
					"expected_output": {
						"test": "123"
					}
				}
			}
		}
	}
}
`

	err := cli.TestPipeline(config, func() {})

	assert.Error(t, err)
}
