/*
Package spinnaker - a `shore` backend implementation for Spinnaker APIs

An abstraction layer over Spinnaker communications, unifying the experience for `shore` developers when working with a `spinnaker` backend.

The abstraction implements the standard `shore-backend` interface from github.com/Autodesk/shore/pkg/backend/spinnaker
*/
package spinnaker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/Autodesk/shore/internal/retry"
	"github.com/Autodesk/shore/pkg/shore_testing"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	spinGate "github.com/spinnaker/spin/cmd/gateclient"
	spinGateApi "github.com/spinnaker/spin/gateapi"
)

const defaultTestTimeout = 1200 // 20 minutes in seconds

// ApplicationControllerAPI - Interface wrapper for the Application Controller API
type ApplicationControllerAPI interface {
	GetPipelineConfigUsingGET(ctx context.Context, application string, pipelineName string) (map[string]interface{}, *http.Response, error)
}

type TestPipelineResponse struct {
	testName         string
	test             shore_testing.TestConfig
	pipelineID       string
	response         *http.Response
	ExecutionDetails *PipelineExecutionDetailsResponse
	err              error
}

// PipelineControllerAPI - Interface wrapper for the Pipeline Controller API
type PipelineControllerAPI interface {
	SavePipelineUsingPOST(ctx context.Context, pipeline interface{}, localVarOptionals *spinGateApi.PipelineControllerApiSavePipelineUsingPOSTOpts) (*http.Response, error)
	DeletePipelineUsingDELETE(ctx context.Context, application string, pipelineName string) (*http.Response, error)
}

// SpinCLI is a wrapper for the spin-cli gateway client backed by swagger
type SpinCLI struct {
	ApplicationControllerAPI
	PipelineControllerAPI
	context.Context
}

// SpinClient - Concrete type requiring all the methods of the specified interfaces.
type SpinClient struct {
	initOnce sync.Once
	*SpinCLI
	CustomSpinCLI
	log logrus.FieldLogger
}

type DeletePipelineResponse struct {
	StatusCode int
	App        string
	Name       string
}

func (s *SpinClient) GetPipeline(application string, pipelineName string) (map[string]interface{}, *http.Response, error) {
	var outErr error

	err := s.initializeAPI()
	if err != nil {
		outErr = err
		return nil, nil, outErr
	}

	pipeline, res, err := s.ApplicationControllerAPI.GetPipelineConfigUsingGET(s.Context, application, pipelineName)

	if err != nil {
		outErr = err
		return nil, nil, outErr
	}

	return pipeline, res, err
}

// NewClient - Create a new default spinnaker client
func NewClient(logger logrus.FieldLogger) *SpinClient {
	return &SpinClient{log: logger}
}

// initializeAPI - Lazy initialization of the client, is expected to be called before each method that requires http.
// Concept taken from: https://roberto.selbach.ca/zero-values-in-go-and-lazy-initialization/
func (s *SpinClient) initializeAPI() error {
	var outerErr error
	// If the client is already initialized, not
	if s.SpinCLI == nil && s.CustomSpinCLI == nil {
		s.initOnce.Do(func() {
			gateClient, err := spinGate.NewGateClient(&UI{}, "", "", "", false)

			if err != nil {
				outerErr = err
				return
			}

			// `InitializeHTTPClient` is an internal Spinnaker function that takes the auth config from a `.config` file used by spin-cli.
			// The returned `httpClient` is already configured to use `LDAP/OAuth2/Certificates` and the other authentication methods provided by spinnaker.
			httpClient, err := spinGate.InitializeHTTPClient(gateClient.Config.Auth)

			if err != nil {
				outerErr = err
				return
			}

			s.CustomSpinCLI = &CustomSpinClient{Endpoint: gateClient.Config.Gate.Endpoint, HTTPClient: httpClient}
			s.SpinCLI = &SpinCLI{
				ApplicationControllerAPI: gateClient.ApplicationControllerApi,
				PipelineControllerAPI:    gateClient.PipelineControllerApi,
				Context:                  gateClient.Context,
			}
		})
	}

	return outerErr
}

func (s *SpinClient) getOtherPipelineId(application string, pipelineName string) (string, *http.Response, error) {
	pipeline, res, err := s.ApplicationControllerAPI.GetPipelineConfigUsingGET(s.Context, application, pipelineName)
	if err != nil {
		return "", res, NewApplicationControllerError(err, res)
	}

	if !mapContainsKey(pipeline, "id") {
		return "", res, nil
	}

	return pipeline["id"].(string), res, nil
}

func mapContainsKey(mapInput map[string]interface{}, searchString string) bool {
	_, found := mapInput[searchString]
	return found
}

func (s *SpinClient) findAndReplacePipelineNameWithFoundID(spinnakerObject map[string]interface{}) (bool, map[string]interface{}) {
	s.log.WithFields(logrus.Fields{"spinnaker_object": spinnakerObject}).Info("Found spinnaker object with 'application' and 'pipeline' name fields")
	pipelineApp := spinnakerObject["application"].(string)
	pipelineName := spinnakerObject["pipeline"].(string)

	isPipelineUUID, err := isValidv4UUIDtypeRFC4122(pipelineName)
	isPipelineSpEL := isSpEL(pipelineName)
	if !(isPipelineUUID || isPipelineSpEL) {
		s.log.WithFields(logrus.Fields{
			"pipeline_name": pipelineName,
			"uuid_error":    err,
		}).Info("Checking if provided pipeline name is not already a valid pipeline UUID or SpEL expression, looking for existing pipeline.")

		newID, res, err := s.getOtherPipelineId(pipelineApp, pipelineName)
		newSpinnakerObject := spinnakerObject

		if err != nil && res.StatusCode == 404 {
			s.log.WithFields(logrus.Fields{
				"pipeline_name": pipelineName,
				"application":   pipelineApp,
				"status code":   res.StatusCode,
			}).Warn("Failed to find a matching pipeline")
			newSpinnakerObject["pipeline"] = nil
		} else {
			newSpinnakerObject["pipeline"] = newID
		}

		s.log.WithFields(logrus.Fields{
			"pipeline_name": pipelineName,
			"pipeline_id":   newID,
			"application":   pipelineApp,
		}).Info("Replacing pipeline name with valid pipeline UUID from specified application.")

		return true, newSpinnakerObject
	} else {
		s.log.WithFields(logrus.Fields{
			"pipeline_name": pipelineName,
			"application":   pipelineApp,
		}).Info("Provided pipeline name is already a valid pipeline UUID or a SpEL expression")
	}

	return false, spinnakerObject
}

func isValidv4UUIDtypeRFC4122(u string) (bool, error) {
	parsedUUID, err := uuid.Parse(u)

	if err != nil {
		return false, err
	}

	if parsedUUID.Version() != 4 && parsedUUID.Variant() == uuid.RFC4122 {
		return false, errors.New("UUID is not VERSION_4 or of variant RFC4122")
	}

	return true, nil
}

func isSpEL(u string) bool {
	return strings.Contains(u, "${")
}

// TODO: We have to implement transaction based saving everywhere - at the moment if something goes off the state is undefined.
func (s *SpinClient) savePipeline(pipelineJSON string) (string, *http.Response, error) {
	var pipeline map[string]interface{}
	pipelineID := ""

	if err := s.initializeAPI(); err != nil {
		return pipelineID, &http.Response{}, err
	}

	err := jsoniter.Unmarshal([]byte(pipelineJSON), &pipeline)

	if err != nil {
		return pipelineID, &http.Response{}, err
	}

	if err = s.isValidPipeline(pipeline); err != nil {
		return pipelineID, &http.Response{}, err
	}

	application := pipeline["application"].(string)
	pipelineName := pipeline["name"].(string)

	if template, exists := pipeline["template"]; exists && len(template.(map[string]interface{})) > 0 {
		if _, exists := pipeline["schema"]; !exists {
			return pipelineID, &http.Response{}, fmt.Errorf("required pipeline key 'schema' missing for templated pipeline")
		}
		pipeline["type"] = "templatedPipeline"
	}

	foundPipeline, queryResp, err := s.ApplicationControllerAPI.GetPipelineConfigUsingGET(s.Context, application, pipelineName)
	if err != nil {
		wrappedErr := NewApplicationControllerError(err, queryResp)
		if wrappedErr.StatusCode() != http.StatusNotFound {
			return pipelineID, nil, wrappedErr
		}

		s.log.Info("Pipeline %q not found in application %q", pipelineName, application)
	}

	// pipeline found, let's use Spinnaker's known Pipeline ID, otherwise we'll get one created for us
	if len(foundPipeline) > 0 {
		s.log.Info("Pipeline %q found with ID %q", foundPipeline["name"], foundPipeline["id"], application)

		pipeline["id"] = foundPipeline["id"].(string)
		pipelineID = foundPipeline["id"].(string)
	}

	res, err := s.PipelineControllerAPI.SavePipelineUsingPOST(s.Context, pipeline, nil)
	if err != nil {
		return pipelineID, res, NewApplicationControllerError(err, res)
	}

	return pipelineID, res, nil
}

// ExecutePipeline - Execute a spinnaker pipeline.
//
// `patameters` are optional.
func (s *SpinClient) ExecutePipeline(argsJSON string, stringify bool) (string, *http.Response, error) {
	// For some crazy reason, spincli invoke doesn't return the ID of the pipeline execution.
	// BTW the crazy reason is that `swagger-code-gen` produces wrong code and Spin-Cli (and shore...) depends on this wrong code.
	// So this request needs to be done 100% manually.
	// For this reason we use the `CustomSpinCLI` interface to implement all the things `SpinCli` either does wrong, or is 100% broken.

	var args map[string]interface{}

	if err := s.initializeAPI(); err != nil {
		return "", &http.Response{}, err
	}

	err := jsoniter.Unmarshal([]byte(argsJSON), &args)

	if err != nil {
		return "", &http.Response{}, err
	}

	if err = s.isArgsValid(args); err != nil {
		return "", &http.Response{}, err
	}

	application := args["application"].(string)
	pipelineName := args["pipeline"].(string)

	// The Spinnaker API is very weird in terms of JSON
	// Sending JSON as is just kills the request (400) status code.
	// However, if the JSON is stringified (example {"a": "a"} -> "{\"a\": \"a\"}")
	// the request will work fine.
	// This is due to the fact that the parameters API can only handle scalar values (string, int, bool)
	// This logic checks if a pipeline parameter looks like: {"key": "value"}, ["key"], [{"key": "value"}]
	// If the value is of one of the example types, the algorithm will stringify the property before the request is sent.
	if params, exists := args["parameters"]; exists {
		if reflect.TypeOf(params).Kind() != reflect.Map {
			return "", &http.Response{}, fmt.Errorf("`parameters` must be an object")
		}
		if stringify {
			parameters := args["parameters"].(map[string]interface{})

			for key, val := range parameters {
				switch v := val.(type) {
				case map[string]interface{}, []interface{}:
					{
						semiMarshal, _ := jsoniter.Marshal(v)
						parameters[key] = string(semiMarshal)
					}
				}
			}
		}
	}

	// Check if artifacts are present, and if it's an array.
	if artifacts, exists := args["artifacts"]; exists {
		if reflect.TypeOf(artifacts).Kind() != reflect.Slice {
			return "", &http.Response{}, fmt.Errorf("`artifacts` must be an Array")
		}

		artifactsSlice, _ := artifacts.([]interface{})
		s.log.Debug("Number of artifacts: ", len(artifactsSlice))
		// Ideally would check the structure of each artifact so that it's correct - beyond just checking that it's an object/map
		for _, artifact := range artifactsSlice {
			if reflect.TypeOf(artifact).Kind() != reflect.Map {
				return "", &http.Response{}, fmt.Errorf("an artifact in `artifacts` must be an object")
			}
		}
	}

	delete(args, "application")
	delete(args, "pipeline")

	argsBytes, err := jsoniter.Marshal(args)

	if err != nil {
		return "", &http.Response{}, err
	}

	body, res, err := s.CustomSpinCLI.ExecutePipeline(application, pipelineName, bytes.NewBuffer(argsBytes))

	if len(body.Ref) == 0 {
		return "", res, err
	}

	refID := strings.Split(body.Ref, "/")[2]

	return refID, res, err
}

// TestPipeline - Run a Spinnaker testing
// Currently returns a not-so-well formatted error.
// The intended solution is to create a shared API between `shore-cli` & the `backend` to expect well formatted struct for the CLI to render correctly.
// TODO: Design a struct to pass data back to `shore-cli` so the UI layer could render the test-results correctly.
func (s *SpinClient) TestPipeline(testConfig shore_testing.TestsConfig, onChange func(), stringify bool) error {
	s.log.Info("Starting test suite")
	var err error

	// Validating test parameters
	if err := s.initializeAPI(); err != nil {
		return err
	}

	if testConfig.Application == "" {
		return fmt.Errorf("test config missing required property `application`")
	}

	if testConfig.Pipeline == "" {
		return fmt.Errorf("test config missing required property `pipeline`")
	}

	if testConfig.Timeout == 0 {
		s.log.Info(fmt.Sprintf("Detected a timeout of 0 sec for testing, defaulting to %d seconds.", defaultTestTimeout))
		testConfig.Timeout = defaultTestTimeout
	} else if testConfig.Timeout < 0 {
		return fmt.Errorf("test config specifies the property `timeout` as %d seconds, but it must be greater than 0", testConfig.Timeout)
	}

	configuredTestNames := make([]string, 0)
	for testName := range testConfig.Tests {
		configuredTestNames = append(configuredTestNames, testName)
	}

	var testsToRun []string
	if testConfig.Ordering != nil && len(testConfig.Ordering) > 0 {
		testsToRun = testConfig.Ordering
	} else {
		testsToRun = configuredTestNames
	}

	// TODO: Rethink the channel size (https://github.com/Autodesk/shore/pull/200#discussion_r2847971)
	var ch = make(chan *TestPipelineResponse, len(testsToRun))
	var wg = sync.WaitGroup{}
	testErrors := make(map[string][]string)

	for _, testName := range testsToRun {
		if testConfig.Parallel {
			// We add 1 to the wait group - each worker will decrease it back
			wg.Add(1)
			go func(testName string, test shore_testing.TestConfig, stringify bool) {
				defer wg.Done()
				ch <- s.RunTest(testName, testConfig, testErrors, stringify)
			}(testName, testConfig.Tests[testName], stringify)
		} else {
			ch <- s.RunTest(testName, testConfig, testErrors, stringify)
		}
	}

	wg.Wait()
	close(ch)

	for testResponse := range ch {
		err = s.ValidateTestResponses(*testResponse, testConfig.Timeout, testErrors)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SpinClient) RunTest(testName string, testConfig shore_testing.TestsConfig, testErrors map[string][]string, stringify bool) *TestPipelineResponse {

	s.log.Info(fmt.Sprintf("Running test %s", testName))
	test := testConfig.Tests[testName]
	if test.ExecArgs == nil {
		test.ExecArgs = make(map[string]interface{})
	}

	test.ExecArgs["application"] = testConfig.Application
	test.ExecArgs["pipeline"] = testConfig.Pipeline

	execArgs, err := jsoniter.Marshal(test.ExecArgs)

	if err != nil {
		return &TestPipelineResponse{testName, test, "", nil, nil, err}
	}

	s.log.Info("Executing pipeline for test: ", testName)
	refID, res, err := s.ExecutePipeline(string(execArgs), stringify)
	if err != nil {
		return &TestPipelineResponse{testName, test, "", nil, nil, err}
	}

	var execDetails *PipelineExecutionDetailsResponse
	execDetails, _, err = s.CustomSpinCLI.PipelineExecutionDetails(refID, bytes.NewBuffer(make([]byte, 0)))
	if execDetails.Status == PipelineRunning || execDetails.Status == PipelineNotStarted {
		s.log.Info("Waiting for pipeline to finish for test: ", testName)

		execDetails, _, err = s.waitForPipelineToFinish(refID, testConfig.Timeout)
	}
	if err != nil {
		testErrors[testName] = append(testErrors[testName], err.Error())
	}
	return &TestPipelineResponse{testName, test, refID, res, execDetails, err}
}

func (s *SpinClient) ValidateTestResponses(testResponse TestPipelineResponse, timeout int, testErrors map[string][]string) error {
	if testResponse.err != nil {
		s.log.Debug("Pipeline execution failed for test: ", testResponse.testName)
		testErrors[testResponse.testName] = append(testErrors[testResponse.testName], testResponse.err.Error())
	}

	// TODO: Need to check what happens in a 404 case and format an error for it.
	if len(testResponse.pipelineID) == 0 {
		return nil
	}

	for _, stage := range testResponse.ExecutionDetails.Stages {
		stageName := stage["name"].(string)

		assertion, exists := testResponse.test.Assertions[stageName]

		if !exists {
			testErrors[testResponse.testName] = append(testErrors[testResponse.testName], fmt.Sprintf("missing assertion for stage %s", stageName))
			continue
		}

		expectedStatus := strings.ToUpper(assertion.ExpectedStatus)

		if err := isExpectedStatus(expectedStatus, stage["status"].(string), stageName); err != nil {
			testErrors[testResponse.testName] = append(testErrors[testResponse.testName], err.Error())
		}

		if err := isExpectedOutput(assertion.ExpectedOutput, stage["outputs"].(map[string]interface{}), stageName); err != nil {
			testErrors[testResponse.testName] = append(testErrors[testResponse.testName], err.Error())
		}
	}

	// This should really be handled by a Golang template.
	// TODO: move this logic to the CLI layer.
	if len(testErrors) > 0 {
		formattedErrors := ""

		for testName, testErrors := range testErrors {
			formattedErrors += fmt.Sprintf("`%s` failure:\n", testName)
			for _, testError := range testErrors {
				formattedErrors += fmt.Sprintf("%s\n", testError)
			}
			formattedErrors += "\n"
		}

		// TODO: The backend shouldn't concern itself with rendering, this should be replaced with the correct struct response.
		return errors.New(formattedErrors)
	}
	return nil
}

// WaitForPipelineToFinish - Wait for the pipeline to finish running.
// This call uses sleeps and is a blocking call.
func (s *SpinClient) WaitForPipelineToFinish(refID string, timeout int) (string, *http.Response, error) {
	execDetails, res, err := s.waitForPipelineToFinish(refID, timeout)

	data, marshalErr := jsoniter.Marshal(execDetails)

	if err != nil {
		return "", res, multierror.Append(err, marshalErr)
	}

	return string(data), res, err
}

// The actual implementation for WaitForPipelineToFinish.
// This implementation is hidden to allow internal package code to use *PipelineExecutionDetailsResponse, without exposing internal package logic.
func (s *SpinClient) waitForPipelineToFinish(refID string, timeout int) (*PipelineExecutionDetailsResponse, *http.Response, error) {
	var errors error

	// Simple reverse
	tries := int(math.Floor(math.Sqrt(float64(timeout * 2))))

	retryConfig := retry.Config{
		Tries: tries,
		// Linear regresion, wait for the the amount of seconds that this "try" matches.
		// I.E. wait 1 second, 2 seconds, 3 seconds.... etc...
		DelayFunc: func(try int) time.Duration { return time.Duration(time.Second * time.Duration(try)) },
	}

	var execDetails *PipelineExecutionDetailsResponse
	var res *http.Response
	var err error

	retryFunc := func() error {
		execDetails, res, err = s.CustomSpinCLI.PipelineExecutionDetails(refID, bytes.NewBuffer(make([]byte, 0)))

		// Other statuses to consider - PipelinePaused / PipelineSuspended
		if execDetails.Status == PipelineRunning || execDetails.Status == PipelineNotStarted {
			return retry.ErrRetry
		}

		return nil
	}

	if retryErr := retry.Retry(retryFunc, retryConfig); retryErr != nil {
		return execDetails, res, multierror.Append(errors, retryErr, err)
	}

	return execDetails, res, err
}

func (s *SpinClient) isValidPipeline(pipeline map[string]interface{}) error {
	var errorsList []string

	if _, exists := pipeline["name"]; !exists {
		errorsList = append(errorsList, "required pipeline key 'name' missing")
	}

	if _, exists := pipeline["application"]; !exists {
		errorsList = append(errorsList, "required pipeline key 'application' missing")
	}

	if len(errorsList) > 0 {
		return fmt.Errorf(strings.Join(errorsList, "\n"))
	}

	return nil
}

func (s *SpinClient) isArgsValid(args map[string]interface{}) error {
	var errorsList []string

	if _, exists := args["pipeline"]; !exists {
		errorsList = append(errorsList, "required args key 'pipeline' missing")
	}

	if _, exists := args["application"]; !exists {
		errorsList = append(errorsList, "required args key 'application' missing")
	}

	if len(errorsList) > 0 {
		return fmt.Errorf(strings.Join(errorsList, "\n"))
	}

	return nil
}

func isValidPipelineApplication(pipeline map[string]interface{}, parentPipelineApp []string) error {
	if len(parentPipelineApp) > 0 {
		if parentApplication := parentPipelineApp[0]; pipeline["application"].(string) != parentApplication {
			return fmt.Errorf("pipeline 'application' key value should match the value of parent pipeline 'application' key")
		}
	}

	return nil
}

func hasValidChildPipelineStages(stages []interface{}, parentPipelineApp []string) (bool, error) {
	var errorsList []string
	hasPipelineStages := false

	for _, stage := range stages {
		stageMap := stage.(map[string]interface{})

		// if it is a pipeline stage it has to have "name" key.
		if !mapContainsKey(stageMap, "name") {
			errorsList = append(errorsList, "required stage key 'name' missing for stage")
		}

		// if it is a pipeline stage, it must contain "type" of pipeline
		if !mapContainsKey(stageMap, "type") {
			errorsList = append(errorsList, "required stage key 'type' missing for stage")
		} else {
			if stageMap["type"].(string) != "pipeline" {
				continue
			}
		}
		// Application value should match the one with the parent pipeline's one
		if pipeline, exists := stage.(map[string]interface{})["pipeline"]; exists {

			if reflect.TypeOf(pipeline).Kind() != reflect.Map {
				continue
			}

			hasPipelineStages = true
			stageApplication, applicationExists := stage.(map[string]interface{})["application"]
			if !applicationExists {
				errorsList = append(errorsList, "required stage key 'application' missing for stage")
			} else {
				if len(parentPipelineApp) > 0 {
					if parentApplication := parentPipelineApp[0]; stageApplication.(string) != parentApplication {
						errorsList = append(errorsList, "'application' key value of stage of type 'pipeline' should match the one of parent pipeline 'application' value")
					}
				}
			}

		}

		if len(errorsList) > 0 {
			return hasPipelineStages, fmt.Errorf(strings.Join(errorsList, "\n"))
		}
	}

	return hasPipelineStages, nil
}

// Add ctx support to configure polling parameters
func (s *SpinClient) pollSpinnakerGetPipelineConfigUsingGET(application string, pipelineName string) (string, *http.Response, error) {
	var pollingStep int = 10
	var pollingTimeout int = 61

	for pollingSleep := 1; pollingSleep <= pollingTimeout; pollingSleep += pollingStep {
		foundPipeline, queryResp, err := s.ApplicationControllerAPI.GetPipelineConfigUsingGET(s.Context, application, pipelineName)
		if err != nil || queryResp.StatusCode != http.StatusOK {
			return "", queryResp, fmt.Errorf(
				"nested pipeline %q wasn't created in application %q, pipeline Stage will be unbound: %w",
				pipelineName,
				application,
				NewApplicationControllerError(err, queryResp))
		}

		if len(foundPipeline) == 0 {
			s.log.Info("get pipeline request didn't return a payload, sleeping for", pollingSleep)
			time.Sleep(time.Duration(pollingSleep))
			continue
		}

		s.log.Info("Pipeline", foundPipeline["name"], "found with id", foundPipeline["id"], "in application", application)

		return foundPipeline["id"].(string), queryResp, nil
	}

	return "", &http.Response{}, fmt.Errorf("couldn't get pipeline until hitting timeout")
}

// A DFS implementation that runs through the pipeline/stages tree.
// Finds child pipelines and saves them.
// Each iteration of the stages loop will look for "pipeline" key in each element
// If it finds one and it is another pipeline it will start another iteration on that pipeline's stages
// Once a pipeline with no "pipeline" stages is met - it is saved
// It's pipeline UUID is assigned to the parent pipeline relevant stage's "pipeline" key
// The parent loop continues until all its stages of type pipeline are updated with pipeline UUIDs and it is saved.
// Once all loops are closed the most top level pipeline has all all stage's pipelines replaced with UUIDs and it is saved
func (s *SpinClient) saveNestedPipeline(stages interface{}, pipeline map[string]interface{}) error {
	for _, stage := range stages.([]interface{}) {
		stage := stage.(map[string]interface{})
		stagePipelineField, exists := stage["pipeline"]
		if !exists {
			continue
		}

		switch stagePipelineField.(type) {
		case string:
			continue
		}

		childPipeline := stage["pipeline"].(map[string]interface{})

		if err := s.isValidPipeline(childPipeline); err != nil {
			return err
		}

		parentPipelineApplication := []string{pipeline["application"].(string)}

		if err := isValidPipelineApplication(childPipeline, parentPipelineApplication); err != nil {
			return err
		}

		childPipelineStages, exists := childPipeline["stages"]
		if exists {

			hasChildPipelines, err := hasValidChildPipelineStages(childPipelineStages.([]interface{}), parentPipelineApplication)
			if err != nil {
				return err
			}

			// If any of stages is of type pipeline create those pipelines recursively
			if hasChildPipelines {
				if err := s.saveNestedPipeline(childPipelineStages, childPipeline); err != nil {
					return err
				}
			}

			//  Replace upstream pipeline name string with real pipeline ID
			s.log.Info("Searching for Stages with PipelineIDs needing replacement")
			for _, stage := range childPipelineStages.([]interface{}) {
				innerStage := stage.(map[string]interface{})

				if !hasChildPipelines && mapContainsKey(innerStage, "application") && mapContainsKey(innerStage, "pipeline") && reflect.TypeOf(innerStage["pipeline"]).Kind() == reflect.String {
					if response, result := s.findAndReplacePipelineNameWithFoundID(innerStage); response {
						stage = result
					}
				}
			}
		}

		// After we return from recursion we save "this layer" child pipeline
		childPipelineBytes, err := jsoniter.Marshal(childPipeline)
		if err != nil {
			return err
		}

		pipelineID, res, err := s.savePipeline(string(childPipelineBytes))
		if err != nil {
			return err
		}
		s.log.Info(res)

		// Do not try to poll for pipeline ID again if exists already
		if pipelineID == "" {
			pipelineID, res, err = s.pollSpinnakerGetPipelineConfigUsingGET(childPipeline["application"].(string), childPipeline["name"].(string))
			if err != nil {
				return err
			}
			// pipelineID = pipelineID

			s.log.Info("Created new pipeline with id:", pipelineID)
			s.log.Info(res)
		}

		// And override stage pipeline value with an the id (a UUID String) received from spinnaker pipeline.
		// maps are updated by reference so the save of the parent pipeline will be updated as well
		stage["pipeline"] = pipelineID
	}

	return nil
}

// SavePipeline - Creates or Update nested pipelines recursively
func (s *SpinClient) SavePipeline(pipelineJSON string) (*http.Response, error) {

	if err := s.initializeAPI(); err != nil {
		return &http.Response{}, err
	}

	var pipeline map[string]interface{}

	err := jsoniter.Unmarshal([]byte(pipelineJSON), &pipeline)

	if err != nil {
		log.Fatal(err)
	}

	if err := s.isValidPipeline(pipeline); err != nil {
		return &http.Response{}, err
	}

	s.log.Info("Searching for Triggers with PipelineID needing replacement")

	if triggers, exists := pipeline["triggers"]; exists {
		//  Replace upstream pipeline name string with real pipeline ID
		for _, trigger := range triggers.([]interface{}) {
			triggerObj := trigger.(map[string]interface{})
			if mapContainsKey(triggerObj, "application") && mapContainsKey(triggerObj, "pipeline") {
				if response, result := s.findAndReplacePipelineNameWithFoundID(triggerObj); response {
					trigger = result
				}
			}
		}
	}

	// Check whether a pipeline has stages list
	if stages, exists := pipeline["stages"]; exists {

		//  Replace upstream pipeline name string with real pipeline ID
		s.log.Info("Searching for Stages with PipelineIDs needing replacement")
		for _, stage := range stages.([]interface{}) {
			stage := stage.(map[string]interface{})

			if mapContainsKey(stage, "application") && mapContainsKey(stage, "pipeline") {
				pipelineDataType := reflect.TypeOf(stage["pipeline"])
				if pipelineDataType.Kind() == reflect.String {
					if response, result := s.findAndReplacePipelineNameWithFoundID(stage); response {
						stage = result
					}
				}
			}
		}

		hasChildPipelines, err := hasValidChildPipelineStages(stages.([]interface{}), []string{pipeline["application"].(string)})
		if err != nil {
			return &http.Response{}, err
		}

		// If any of stages is of type pipeline create those pipelines recursively
		if hasChildPipelines {
			if err := s.saveNestedPipeline(stages, pipeline); err != nil {
				return &http.Response{}, err
			}
		}
	}

	pipelineBytes, err := jsoniter.Marshal(pipeline)
	if err != nil {
		return &http.Response{}, err
	}

	pipelineID, res, err := s.savePipeline(string(pipelineBytes))
	if err != nil {
		return &http.Response{}, err
	}

	if pipelineID != "" {
		s.log.Info("Saved already existing pipeline with ID", pipelineID)
	}

	return res, nil
}

// A DFS implementation that runs through the pipeline/stages tree.
// Finds child pipelines and deletes them.
// Each iteration of the stages loop will look for "pipeline" key in each element
// If it finds one and it is another pipeline it will start another iteration on that pipeline's stages
// Once a pipeline with no "pipeline" stages is met - it is deleted
// The parent loop continues until all its stages of type pipeline are deleted.
// Once all loops are closed the most top level pipeline gets deleted.
func (s *SpinClient) getNestedPipelinesNames(stages interface{}, pipeline map[string]interface{}) ([]string, error) {
	pipelineNames := []string{}
	for _, stage := range stages.([]interface{}) {
		stage := stage.(map[string]interface{})
		stagePipelineField, exists := stage["pipeline"]
		if !exists {
			continue
		}

		switch stagePipelineField.(type) {
		case string:
			continue
		}

		childPipeline := stage["pipeline"].(map[string]interface{})

		if err := s.isValidPipeline(childPipeline); err != nil {
			return pipelineNames, err
		}

		parentPipelineApplication := []string{pipeline["application"].(string)}

		if err := isValidPipelineApplication(childPipeline, parentPipelineApplication); err != nil {
			return pipelineNames, err
		}

		childPipelineStages, exists := childPipeline["stages"]
		if exists {
			hasChildPipelines, err := hasValidChildPipelineStages(childPipelineStages.([]interface{}), parentPipelineApplication)
			if err != nil {
				return pipelineNames, err
			}

			// If any of stages is of type pipeline get those pipelines names recursively
			if hasChildPipelines {
				nestedPipelineNames, err := s.getNestedPipelinesNames(childPipelineStages, childPipeline)
				if err != nil {
					return pipelineNames, err
				}
				pipelineNames = append(pipelineNames, nestedPipelineNames...)
			}
		}

		// After we return from recursion we add "this layer" child pipeline
		childPipelineName := childPipeline["name"].(string)
		pipelineNames = append(pipelineNames, childPipelineName)
	}

	return pipelineNames, nil
}

func (s *SpinClient) getPipelinesNames(pipeline map[string]interface{}) ([]string, error) {
	pipelineNames := []string{}

	if stages, exists := pipeline["stages"]; exists {

		hasChildPipelines, err := hasValidChildPipelineStages(stages.([]interface{}), []string{pipeline["application"].(string)})
		if err != nil {
			return pipelineNames, err
		}

		// If any of stages is of type pipeline collect those pipelines names recursively
		if hasChildPipelines {
			pipelineNames, err = s.getNestedPipelinesNames(stages, pipeline)
			if err != nil {
				return pipelineNames, err
			}
		}
	}

	pipelineName := pipeline["name"].(string)
	pipelineNames = append(pipelineNames, pipelineName)

	return pipelineNames, nil
}

// GetPipelinesNamesAndApplication - gets list of names of all pipelines and application name configured
func (s *SpinClient) GetPipelinesNamesAndApplication(pipelineJSON string) ([]string, string, error) {

	var pipeline map[string]interface{}

	err := jsoniter.Unmarshal([]byte(pipelineJSON), &pipeline)

	if err != nil {
		log.Fatal(err)
	}

	if err := s.isValidPipeline(pipeline); err != nil {
		return []string{}, "", err
	}

	application := pipeline["application"].(string)

	pipelineNames, err := s.getPipelinesNames(pipeline)

	if err != nil {
		return []string{}, "", err
	}

	return pipelineNames, application, nil
}

// DeletePipeline - deletes rendered pipeline (recursively, if there are nested pipelines)
func (s *SpinClient) DeletePipeline(pipelineJSON string) (*http.Response, error) {

	if err := s.initializeAPI(); err != nil {
		return &http.Response{}, err
	}

	pipelineNames, application, err := s.GetPipelinesNamesAndApplication(pipelineJSON)

	if err != nil {
		return &http.Response{}, err
	}
	// here and bellow %-40s hack is to replace possible spinner suffix interferring with the output
	color.Yellow(fmt.Sprintf("\rApplication: %-40s ", application))
	color.Yellow(fmt.Sprintf("Pipelines to delete: %s%-20s", pipelineNames, ""))

	ch := make(chan DeletePipelineResponse, len(pipelineNames))
	errCh := make(chan error)

	go s.DeletePipelines(application, pipelineNames, ch, errCh)
	err = <-errCh

	deletions := []DeletePipelineResponse{}

	for deletion := range ch {
		deletions = append(deletions, deletion)
	}
	fmt.Printf("\r")
	for _, d := range deletions {
		color.Red(fmt.Sprintf("DELETED: %s - %-40s", d.App, d.Name))
	}

	if err != nil {
		return &http.Response{StatusCode: http.StatusBadRequest}, err
	}

	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
	}, nil
}

// DeletePipelines - asynchronously deletes pipelines list from the application provided as input.
// Writes results to ch.
// Error encountered is written to errCh and waits till all running in parallel deletion goroutines finish.
func (s *SpinClient) DeletePipelines(application string, pipelineNames []string, ch chan DeletePipelineResponse, errCh chan error) {
	if err := s.initializeAPI(); err != nil {
		errCh <- err
		return
	}

	wg := sync.WaitGroup{}
	for _, pipelineName := range pipelineNames {
		wg.Add(1)
		go func(appName string, pipelineName string, ch chan DeletePipelineResponse, errCh chan error) {
			defer wg.Done()
			res, err := s.PipelineControllerAPI.DeletePipelineUsingDELETE(s.Context, appName, pipelineName)

			ch <- DeletePipelineResponse{
				StatusCode: res.StatusCode,
				App:        appName,
				Name:       pipelineName,
			}
			if err != nil {
				errCh <- err
			}
		}(application, pipelineName, ch, errCh)
	}

	wg.Wait()
	close(errCh)
	close(ch)
}
