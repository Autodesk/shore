package spinnaker

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Autodesk/shore/pkg/shore_testing"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	testLog "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

var cli *SpinClient

var logger *logrus.Logger

func init() {
	logger, _ = testLog.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	cli = &SpinClient{
		log:           logger,
		CustomSpinCLI: &MockCustomSpinCli{},
		SpinCLI: &SpinCLI{
			ApplicationControllerAPI: &MockApplicationControllerAPI{},
			PipelineControllerAPI:    &MockPipelineControllerAPI{},
			Context:                  context.Background(),
		},
	}
}

func TestValidv4UUIDOfVariantRFC1422(t *testing.T) {
	_, err := isValidv4UUIDtypeRFC4122("642355a8-eded-4a73-9d49-2a7af7395f4a")
	assert.Nil(t, err)
}

func TestInvalidUUID3(t *testing.T) {
	_, err := isValidv4UUIDtypeRFC4122("a3bb189e-8bf9-3888-9912-ace4e6543002")
	assert.NotNil(t, err)
}

func TestInvalidUUID1(t *testing.T) {
	_, err := isValidv4UUIDtypeRFC4122("85169e50-ceb4-11eb-8382-000000000000")
	assert.NotNil(t, err)
}

func TestInvalidLengthv4UUIDOfVariantRFC1422(t *testing.T) {
	_, err := isValidv4UUIDtypeRFC4122("1234")
	assert.EqualError(t, err, "invalid UUID length: 4")
}

func TestIsSpEL(t *testing.T) {
	err := isSpEL("${ami_application}")
	assert.True(t, err)
}

func TestPipelineNameIsntSpEL(t *testing.T) {
	err := isSpEL("pipeline_name")
	assert.False(t, err)
}

func TestValidv4UUIDOfVariantRFC1422isNotSpEL(t *testing.T) {
	err := isSpEL("642355a8-eded-4a73-9d49-2a7af7395f4a")
	assert.False(t, err)
}

func TestInvalidFindAndReplacePipelineNameWithFoundID(t *testing.T) {
	stageMap := []map[string]interface{}{
		{
			"application": "test-source-app",
			"pipeline":    "test-source-pipeline",
		},
	}
	expectedResult := []map[string]interface{}{
		{
			"application": "test-source-app",
			"pipeline":    "1234",
		},
	}
	_, res := cli.findAndReplacePipelineNameWithFoundID(stageMap[0])
	assert.Equal(t, expectedResult[0], res)
}

func TestValidFindAndReplacePipelineNameWithFoundIDInStage(t *testing.T) {
	stageMap := []map[string]interface{}{
		{
			"application": "test-source-app",
			"pipeline":    "642355a8-eded-4a73-9d49-2a7af7395f4a",
		},
	}
	expectedResult := []map[string]interface{}{
		{
			"application": "test-source-app",
			"pipeline":    "642355a8-eded-4a73-9d49-2a7af7395f4a",
		},
	}

	_, res := cli.findAndReplacePipelineNameWithFoundID(stageMap[0])
	assert.Equal(t, expectedResult[0], res)
}

func TestValidFindAndReplacePipelineNameWithSpelInStage(t *testing.T) {
	stageMap := []map[string]interface{}{
		{
			"application": "test-source-app",
			"pipeline":    "${parameters.pipeline_name}",
		},
	}
	expectedResult := []map[string]interface{}{
		{
			"application": "test-source-app",
			"pipeline":    "${parameters.pipeline_name}",
		},
	}

	_, res := cli.findAndReplacePipelineNameWithFoundID(stageMap[0])
	assert.Equal(t, expectedResult[0], res)
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
	_, res, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "parameters": {"answer": 42}}`, true)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestExecFailedApplication(t *testing.T) {
	// Test
	_, _, err := cli.ExecutePipeline(`{"pipeline": "test"}`, true)

	// Assert
	assert.EqualError(t, err, "required args key 'application' missing")
}

func TestExecFailedName(t *testing.T) {
	// Test
	_, _, err := cli.ExecutePipeline(`{"application": "test"}`, true)

	// Assert
	assert.EqualError(t, err, "required args key 'pipeline' missing")
}

func TestExecArgsFailing(t *testing.T) {
	// Test
	_, _, err := cli.ExecutePipeline(`{}`, true)

	// Assert
	assert.EqualError(t, err, "required args key 'pipeline' missing\nrequired args key 'application' missing")
}

func TestExecParametersMap(t *testing.T) {
	// Test
	_, res, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "parameters": {"answer": {"abc123": "abc123"}}}`, true)
	body, _ := ioutil.ReadAll(res.Request.Body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, `{"parameters":{"answer":"{\"abc123\":\"abc123\"}"}}`, string(body))
}

func TestExecParametersArray(t *testing.T) {
	// Test
	_, res, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "parameters": {"answer": [1,2,3,4, "5", "this is something"]}}`, true)
	body, _ := ioutil.ReadAll(res.Request.Body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, `{"parameters":{"answer":"[1,2,3,4,\"5\",\"this is something\"]"}}`, string(body))
}

func TestExecParametersArrayMap(t *testing.T) {
	// Test
	_, res, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "parameters": {"answer": [1, 2, 3, "abc", {"answer": 42}]}}`, true)
	body, _ := ioutil.ReadAll(res.Request.Body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, `{"parameters":{"answer":"[1,2,3,\"abc\",{\"answer\":42}]"}}`, string(body))
}

func TestExecParametersNonMapFails(t *testing.T) {
	// Test
	_, _, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "parameters": [1, 2, 3, "abc", {"answer": 42}]}`, true)

	// Assert
	assert.EqualError(t, err, "`parameters` must be an object")
}

func TestExecExtraFields(t *testing.T) {
	// Test
	_, res, err := cli.ExecutePipeline(`{"application": "test", "github": {"commit": "deadbeef", "branch": "dev"}, "pipeline": "test", "parameters": {}, "artifacts": []}`, true)
	bodyString, _ := ioutil.ReadAll(res.Request.Body)
	var body map[string]interface{}
	jsoniter.Unmarshal([]byte(bodyString), &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 3, len(body))
	assert.Contains(t, body, `parameters`)
	assert.Contains(t, body, `artifacts`)
	assert.Contains(t, body, `github`)
}

func TestExecArtifactsArray(t *testing.T) {
	// Test
	_, res, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "artifacts": [{"type": "custom/object", "name": "test-artifact-one", "reference": "test-value-one"}, {"type": "custom/object", "name": "test-artifact-two", "reference": "test-value-two"}]}`, true)
	bodyString, _ := ioutil.ReadAll(res.Request.Body)
	var body map[string]interface{}
	jsoniter.Unmarshal([]byte(bodyString), &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 1, len(body))
	assert.Contains(t, body, `artifacts`)
	assert.Equal(t, 2, len(body["artifacts"].([]interface{})))
}

func TestExecArtifactsEmptyArray(t *testing.T) {
	// Test
	_, res, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "artifacts": []}`, true)
	bodyString, _ := ioutil.ReadAll(res.Request.Body)
	var body map[string]interface{}
	jsoniter.Unmarshal([]byte(bodyString), &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 1, len(body))
	assert.Contains(t, body, `artifacts`)
	assert.Equal(t, 0, len(body["artifacts"].([]interface{})))
}

func TestExecArtifactsNonArrayFails(t *testing.T) {
	// Test
	_, _, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "artifacts": {"type": "potato"}}`, true)

	// Assert
	assert.EqualError(t, err, "`artifacts` must be an Array")
}

func TestExecArtifactsArrayNonObjectsFails(t *testing.T) {
	// Test
	_, _, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "artifacts": ["potato", 1, ["apple", "orange"]]}`, true)

	// Assert
	assert.EqualError(t, err, "an artifact in `artifacts` must be an object")
}

func TestExecParametersArrayMapStringifyFalse(t *testing.T) {
	// Test
	_, res, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "parameters": {"answer": [1, 2, 3, "abc", {"answer": 42}]}}`, false)
	body, _ := ioutil.ReadAll(res.Request.Body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, `{"parameters":{"answer":[1,2,3,"abc",{"answer":42}]}}`, string(body))
}

func TestExecParametersArrayMapStringifyTrue(t *testing.T) {
	// Test
	_, res, err := cli.ExecutePipeline(`{"application": "test", "pipeline": "test", "parameters": {"answer": [1, 2, 3, "abc", {"answer": 42}]}}`, true)
	body, _ := ioutil.ReadAll(res.Request.Body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, `{"parameters":{"answer":"[1,2,3,\"abc\",{\"answer\":42}]"}}`, string(body))
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
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"application": "appname",
							"name":  "Child pipeline stage",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2",
								"stages": [
									{
										"name": "Child pipeline 2 stage",
										"application": "appname",
										"type": "pipeline"
									}
								]
							}
						},
						{
							"application": "appname",
							"name":  "Child pipeline stage 2",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2.2",
								"type": "pipeline",
								"stages": [
									{
										"name": "Child pipeline 2.2 stage",
										"application": "appname",
										"type": "pipeline"
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
					"waitTime": 1,
					"type": "wait"
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
				"type": "pipeline",
				"pipeline": {
					"name": "Child pipeline 1",
					"stages": [
						{
							"name": "Child pipeline 2 stage",
							"application": "appname",
							"type": "pipeline"
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
				"name": "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"stages": [
						{
							"name": "Child pipeline 2 stage",
							"application": "appname",
							"type": "pipeline"
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
				"type": "pipeline",
				"pipeline": {
					"application": "another appname",
					"name": "Child pipeline 1",
					"type": "pipeline",
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
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"type": "pipeline",
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
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"name": "child pipeline 1 stage",
							"type": "pipeline"
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

func TestPipelineSaveEmptyTriggersAndStages(t *testing.T) {
	// Given
	pipelineString := `
	{
		"application": "appname",
		"name":  "test-app",
		"triggers": [],
		"stages": []
	}
	`

	// Test
	_, err := cli.SavePipeline(pipelineString)

	// Assert
	assert.NoError(t, err)

}

func TestPipelineSaveTriggersAndStagesWithValidPipelineID(t *testing.T) {
	// Given
	pipelineString := `
	{
		"application": "appname",
		"name":  "test-app",
		"triggers": [
			{
				"application": "test-source-app",
				"type": "pipeline",
				"pipeline": "642355a8-eded-4a73-9d49-2a7af7395f4a"
			}
		],
		"stages": [
			{
				"application": "test-source-app",
				"type": "findArtifactFromExecution",
				"pipeline": "642355a8-eded-4a73-9d49-2a7af7395f4a"
			}
		]
	}
	`

	// Test
	_, err := cli.SavePipeline(pipelineString)

	// Assert
	assert.NoError(t, err)

}

func TestPipelineSaveTriggersAndStagesWithValidPipelineIDPlusNested(t *testing.T) {
	// Given
	pipelineString := `
	{
		"application": "appname",
		"name":  "test-app",
		"triggers": [
			{
				"application": "test-source-app",
				"type": "pipeline",
				"pipeline": "642355a8-eded-4a73-9d49-2a7af7395f4a"
			}
		],
		"stages": [
			{
				"application": "test-source-app",
				"type": "findArtifactFromExecution",
				"pipeline": "642355a8-eded-4a73-9d49-2a7af7395f4a",
				"name": "resolve them artifacts"
			},
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"application": "appname",
							"name":  "Child pipeline stage",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2",
								"stages": [
									{
										"name": "Child pipeline 2 stage",
										"application": "appname",
										"type": "pipeline"
									}
								]
							}
						},
						{
							"application": "appname",
							"name":  "Child pipeline stage 2",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2.2",
								"type": "pipeline",
								"stages": [
									{
										"name": "Child pipeline 2.2 stage",
										"application": "appname",
										"type": "pipeline"
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
	_, err := cli.SavePipeline(pipelineString)

	// Assert
	assert.NoError(t, err)

}

func TestPipelineSaveTriggersAndStagesAndLookupPipelineIDPlusNested(t *testing.T) {
	// Given
	pipelineString := `
	{
		"application": "appname",
		"name":  "test-app",
		"triggers": [
			{
				"application": "test-source-app",
				"type": "pipeline",
				"pipeline": "test-source-pipeline-name"
			}
		],
		"stages": [
			{
				"application": "test-source-app",
				"type": "findArtifactFromExecution",
				"pipeline": "test-source-pipeline-name",
				"name": "resolve them artifacts"
			},
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"application": "appname",
							"name":  "Child pipeline stage",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2",
								"stages": [
									{
										"name": "Child pipeline 2 stage",
										"application": "appname",
										"type": "pipeline"
									}
								]
							}
						},
						{
							"application": "appname",
							"name":  "Child pipeline stage 2",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2.2",
								"type": "pipeline",
								"stages": [
									{
										"name": "Child pipeline 2.2 stage",
										"application": "appname",
										"type": "pipeline"
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
	_, err := cli.SavePipeline(pipelineString)

	// Assert
	assert.NoError(t, err)

}

func TestPipelineSaveTriggersAndStagesAndLookupPipelineID(t *testing.T) {
	// Given
	pipelineString := `
	{
		"application": "appname",
		"name":  "test-app",
		"triggers": [
			{
				"application": "test-source-app",
				"type": "pipeline",
				"pipeline": "test-source-pipeline-name"
			}
		],
		"stages": [
			{
				"application": "test-source-app",
				"type": "findArtifactFromExecution",
				"pipeline": "test-source-pipeline-name"
			}
		]
	}
	`

	// Test
	_, err := cli.SavePipeline(pipelineString)

	// Assert
	assert.NoError(t, err)

}

func TestPipelineSaveWithNestedLookUp(t *testing.T) {
	// Given
	pipelineString := `
	{
		"application": "appname",
		"name":  "test-app",
		"triggers": [
			{
				"application": "test-source-app",
				"type": "pipeline",
				"pipeline": "test-source-pipeline-name"
			}
		],
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"application": "appname",
							"name":  "Child pipeline stage",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2",
								"stages": [
									{
										"name": "Waiting for a better tomorrow",
										"type": "wait"
									}
								]
							}
						},
						{
							"application": "appname",
							"name":  "Child pipeline stage 2",
							"type": "pipeline",
							"pipeline": "Finding Pipeline nested"
						}
					]
				}
			}
		]
	}
	`

	// Test
	res, saveErr := cli.SavePipeline(pipelineString)
	defer res.Body.Close()

	body, bodyErr := ioutil.ReadAll(res.Body)

	var pipeline map[string]interface{}
	jsoniter.Unmarshal(body, &pipeline)

	// Assert
	assert.NoError(t, saveErr)
	assert.NoError(t, bodyErr)
	assert.NotNil(t, pipeline)
}

func TestPipelineSaveTriggersAndStagesWithReplacablePipelineID(t *testing.T) {
	// Given
	pipelineString := `
	{
		"application": "appname",
		"name":  "test-app",
		"triggers": [
			{
				"application": "test-source-app",
				"pipeline": "642355a8-eded-4a73-9d49-2a7af7395f4a",
				"type": "pipeline"
			}
		],
		"stages": [
			{
				"name": "test-stage-name",
				"type": "findArtifactFromExecution",
				"application": "test-source-app",
				"pipeline": "642355a8-eded-4a73-9d49-2a7af7395f4a"
			}
		]
	}
	`

	// Test
	_, err := cli.SavePipeline(pipelineString)

	// Assert
	assert.NoError(t, err)

}

func TestTestingRemoteSuccess(t *testing.T) {
	config := shore_testing.TestsConfig{
		Application: "test1test2test3",
		Pipeline:    "abc",
		Tests: map[string]shore_testing.TestConfig{
			"test success": {
				ExecArgs: map[string]interface{}{
					"parameters": map[string]string{
						"a": "a",
					},
				},
				Assertions: map[string]shore_testing.Assertion{
					"testedname": {
						ExpectedStatus: "succeeded",
						ExpectedOutput: map[string]interface{}{
							"test": "123",
						},
					},
				},
			},
		},
	}

	err := cli.TestPipeline(config, func() {}, true)

	assert.Nil(t, err)
}

func TestTestingRemoteNoTestsFound(t *testing.T) {
	config := shore_testing.TestsConfig{
		Application: "test1test2test3",
		Pipeline:    "abc",
		Ordering:    []string{"non-existant test"},
		Tests: map[string]shore_testing.TestConfig{
			"test success": {
				ExecArgs: map[string]interface{}{
					"parameters": map[string]string{
						"a": "a",
					},
				},
				Assertions: map[string]shore_testing.Assertion{
					"testedname": {
						ExpectedStatus: "succeeded",
						ExpectedOutput: map[string]interface{}{
							"test": "123",
						},
					},
				},
			},
		},
	}

	execError := "`non-existant test` failure:\nmissing assertion for stage testedname\n\n"

	err := cli.TestPipeline(config, func() {}, true)

	assert.NotNil(t, err)
	assert.Equal(t, execError, err.Error())
}

func TestTestingRemoteNoAssertionFailed(t *testing.T) {
	config := shore_testing.TestsConfig{
		Application: "test1test2test3",
		Pipeline:    "abc",
		Tests: map[string]shore_testing.TestConfig{
			"test success": {
				ExecArgs: map[string]interface{}{
					"parameters": map[string]string{
						"a": "a",
					},
				},
				Assertions: map[string]shore_testing.Assertion{
					"testedname": {
						ExpectedStatus: "succeeded",
						ExpectedOutput: map[string]interface{}{
							"test": "1234",
						},
					},
				},
			},
		},
	}

	err := cli.TestPipeline(config, func() {}, true)

	assert.Error(t, err)
}

func TestTestingRemoteNoAssertionForStageError(t *testing.T) {
	config := shore_testing.TestsConfig{
		Application: "test1test2test3",
		Pipeline:    "abc",
		Tests: map[string]shore_testing.TestConfig{
			"test success": {
				ExecArgs:   map[string]interface{}{},
				Assertions: map[string]shore_testing.Assertion{},
			},
		},
	}

	err := cli.TestPipeline(config, func() {}, true)

	assert.Error(t, err)
}

func TestTestingRemoteMissingExecArgs(t *testing.T) {
	config := shore_testing.TestsConfig{
		Application: "test1test2test3",
		Pipeline:    "abc",
		Tests: map[string]shore_testing.TestConfig{
			"test success": {
				Assertions: map[string]shore_testing.Assertion{
					"testedname": {
						ExpectedStatus: "succeeded",
						ExpectedOutput: map[string]interface{}{
							"test": "123",
						},
					},
				},
			},
		},
	}

	err := cli.TestPipeline(config, func() {}, true)

	assert.Nil(t, err)
}

func TestTestingNoApplicationFailed(t *testing.T) {
	config := shore_testing.TestsConfig{
		Pipeline: "abc",
		Tests: map[string]shore_testing.TestConfig{
			"test success": {
				Assertions: map[string]shore_testing.Assertion{
					"testedname": {
						ExpectedStatus: "succeeded",
						ExpectedOutput: map[string]interface{}{
							"test": "123",
						},
					},
				},
			},
		},
	}

	err := cli.TestPipeline(config, func() {}, true)

	assert.Error(t, err)
}

func TestTestingNoPipelineFailed(t *testing.T) {
	config := shore_testing.TestsConfig{
		Application: "test1test2test3",
		Tests: map[string]shore_testing.TestConfig{
			"test success": {
				Assertions: map[string]shore_testing.Assertion{
					"testedname": {
						ExpectedStatus: "succeeded",
						ExpectedOutput: map[string]interface{}{
							"test": "123",
						},
					},
				},
			},
		},
	}

	err := cli.TestPipeline(config, func() {}, true)

	assert.Error(t, err)
}

func TestTestingBadTimeout(t *testing.T) {
	config := shore_testing.TestsConfig{
		Application: "test1test2test3",
		Pipeline:    "abc",
		Timeout:     -1,
		Tests: map[string]shore_testing.TestConfig{
			"test success": {
				ExecArgs: map[string]interface{}{
					"parameters": map[string]string{
						"a": "a",
					},
				},
				Assertions: map[string]shore_testing.Assertion{
					"testedname": {
						ExpectedStatus: "succeeded",
						ExpectedOutput: map[string]interface{}{
							"test": "123",
						},
					},
				},
			},
		},
	}

	err := cli.TestPipeline(config, func() {}, true)

	assert.Error(t, err)
}

func TestTestingTimeout(t *testing.T) {
	config := shore_testing.TestsConfig{
		Application: "timeout-app",
		Pipeline:    "abc",
		Timeout:     1,
		Tests: map[string]shore_testing.TestConfig{
			"test success": {
				ExecArgs: map[string]interface{}{
					"parameters": map[string]string{
						"a": "a",
					},
				},
				Assertions: map[string]shore_testing.Assertion{
					"testedname": {
						ExpectedStatus: "succeeded",
						ExpectedOutput: map[string]interface{}{
							"test": "123",
						},
					},
				},
			},
		},
	}

	err := cli.TestPipeline(config, func() {}, true)

	assert.Error(t, err)
}

func TestTestingRemoteStringifyTrue(t *testing.T) {
	config := shore_testing.TestsConfig{
		Application: "test1test2test3",
		Pipeline:    "abc",
		Tests: map[string]shore_testing.TestConfig{
			"test success": {
				ExecArgs: map[string]interface{}{
					"parameters": map[string]interface{}{
						"answer": map[string]string{
							"abc123": "abc123",
						},
					},
				},
				Assertions: map[string]shore_testing.Assertion{
					"testedname": {
						ExpectedStatus: "succeeded",
						ExpectedOutput: map[string]interface{}{
							"test": "123",
						},
					},
				},
			},
		},
	}

	err := cli.TestPipeline(config, func() {}, true)

	assert.Nil(t, err)
}

func TestTestingRemoteStringifyFalse(t *testing.T) {
	config := shore_testing.TestsConfig{
		Application: "test1test2test3",
		Pipeline:    "abc",
		Tests: map[string]shore_testing.TestConfig{
			"test success": {
				ExecArgs: map[string]interface{}{
					"parameters": map[string]interface{}{
						"answer": map[string]string{
							"abc123": "abc123",
						},
					},
				},
				Assertions: map[string]shore_testing.Assertion{
					"testedname": {
						ExpectedStatus: "succeeded",
						ExpectedOutput: map[string]interface{}{
							"test": "123",
						},
					},
				},
			},
		},
	}

	err := cli.TestPipeline(config, func() {}, false)

	assert.Nil(t, err)
}

func TestDeleteSuccess(t *testing.T) {
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"application": "appname",
							"name":  "Child pipeline stage",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2",
								"stages": [
									{
										"name": "Child pipeline 2 stage",
										"application": "appname",
										"type": "pipeline"
									}
								]
							}
						},
						{
							"application": "appname",
							"name":  "Child pipeline stage 2",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2.2",
								"type": "pipeline",
								"stages": [
									{
										"name": "Child pipeline 2.2 stage",
										"application": "appname",
										"type": "pipeline"
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
	res, err := cli.DeletePipeline(nestedPipelineString)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestDeleteSimplePipelineSuccess(t *testing.T) {
	// Given
	nestedPipelineString := `
	{
		"application": "test123",
		"name": "simple pipeline",
		"stages": [
			 {
					"name": "Wait",
					"waitTime": 1,
					"type": "wait"
			 }
		]
 }
	`

	// Test
	res, err := cli.DeletePipeline(nestedPipelineString)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestMissingApplicationFailedDelete(t *testing.T) {
	// Given
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"name": "Child pipeline 1",
					"stages": [
						{
							"name": "Child pipeline 2 stage",
							"application": "appname",
							"type": "pipeline"
						}
					]
				}
			}
		]
	}
	`

	// Test
	_, err := cli.DeletePipeline(nestedPipelineString)

	// Assert
	assert.EqualError(t, err, "required pipeline key 'application' missing")
}

func TestMissingNameFailedDelete(t *testing.T) {
	// Given
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name": "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"stages": [
						{
							"name": "Child pipeline 2 stage",
							"application": "appname",
							"type": "pipeline"
						}
					]
				}
			}
		]
	}
	`

	// Test
	_, err := cli.DeletePipeline(nestedPipelineString)

	// Assert
	assert.EqualError(t, err, "required pipeline key 'name' missing")
}

func TestPipelineChildPipelineWrongApplicationFailedDelete(t *testing.T) {
	// Given
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"application": "another appname",
					"name": "Child pipeline 1",
					"type": "pipeline",
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
	_, err := cli.DeletePipeline(nestedPipelineString)

	// Assert
	assert.EqualError(t, err, "pipeline 'application' key value should match the value of parent pipeline 'application' key")
}

func TestDeleteDryRunSuccess(t *testing.T) {
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"application": "appname",
							"name":  "Child pipeline stage",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2",
								"stages": [
									{
										"name": "Child pipeline 2 stage",
										"application": "appname",
										"type": "pipeline"
									}
								]
							}
						},
						{
							"application": "appname",
							"name":  "Child pipeline stage 2",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2.2",
								"type": "pipeline",
								"stages": [
									{
										"name": "Child pipeline 2.2 stage",
										"application": "appname",
										"type": "pipeline"
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
	res, err := cli.DeletePipeline(nestedPipelineString)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestGetNestedPipelinesNames(t *testing.T) {
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"application": "appname",
							"name":  "Child pipeline stage",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2",
								"stages": [
									{
										"name": "Child pipeline 2 stage",
										"application": "appname",
										"type": "pipeline"
									}
								]
							}
						},
						{
							"application": "appname",
							"name":  "Child pipeline stage 2",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2.2",
								"type": "pipeline",
								"stages": [
									{
										"name": "Child pipeline 2.2 stage",
										"application": "appname",
										"type": "pipeline"
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
	expectedPipelineNames := []string{"Child pipeline 2", "Child pipeline 2.2", "Child pipeline 1"}

	var pipeline map[string]interface{}

	err := jsoniter.Unmarshal([]byte(nestedPipelineString), &pipeline)

	stages := pipeline["stages"]
	// Test
	pipelineNames, err2 := cli.getNestedPipelinesNames(stages, pipeline)

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, err2)
	assert.Equal(t, expectedPipelineNames, pipelineNames)
}

func TestGetPipelinesNames(t *testing.T) {
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"application": "appname",
							"name":  "Child pipeline stage",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2",
								"stages": [
									{
										"name": "Child pipeline 2 stage",
										"application": "appname",
										"type": "pipeline"
									}
								]
							}
						},
						{
							"application": "appname",
							"name":  "Child pipeline stage 2",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2.2",
								"type": "pipeline",
								"stages": [
									{
										"name": "Child pipeline 2.2 stage",
										"application": "appname",
										"type": "pipeline"
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
	expectedPipelineNames := []string{"Child pipeline 2", "Child pipeline 2.2", "Child pipeline 1", "Nested pipeline"}

	var pipeline map[string]interface{}

	err := jsoniter.Unmarshal([]byte(nestedPipelineString), &pipeline)

	// Test
	pipelineNames, err2 := cli.getPipelinesNames(pipeline)

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, err2)
	assert.Equal(t, expectedPipelineNames, pipelineNames)
}

func TestGetPipelinesNamesAndApplication(t *testing.T) {
	nestedPipelineString := `
	{
		"application": "appname",
		"name":  "Nested pipeline",
		"stages": [
			{
				"application": "appname",
				"name":  "Nested pipeline stage",
				"type": "pipeline",
				"pipeline": {
					"application": "appname",
					"name": "Child pipeline 1",
					"stages": [
						{
							"application": "appname",
							"name":  "Child pipeline stage",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2",
								"stages": [
									{
										"name": "Child pipeline 2 stage",
										"application": "appname",
										"type": "pipeline"
									}
								]
							}
						},
						{
							"application": "appname",
							"name":  "Child pipeline stage 2",
							"type": "pipeline",
							"pipeline": {
								"application": "appname",
								"name": "Child pipeline 2.2",
								"type": "pipeline",
								"stages": [
									{
										"name": "Child pipeline 2.2 stage",
										"application": "appname",
										"type": "pipeline"
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
	expectedPipelineNames := []string{"Child pipeline 2", "Child pipeline 2.2", "Child pipeline 1", "Nested pipeline"}
	expectedApplication := "appname"

	// Test
	pipelineNames, application, err := cli.GetPipelinesNamesAndApplication(nestedPipelineString)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedPipelineNames, pipelineNames)
	assert.Equal(t, expectedApplication, application)
}
