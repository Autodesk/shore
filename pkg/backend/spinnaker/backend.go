package spinnaker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	spinGate "github.com/spinnaker/spin/cmd/gateclient"
	spinGateApi "github.com/spinnaker/spin/gateapi"
)

// ApplicationControllerAPI - Interface wrapper for the Application Controller API
type ApplicationControllerAPI interface {
	GetPipelineConfigUsingGET(ctx context.Context, application string, pipelineName string) (map[string]interface{}, *http.Response, error)
}

// PipelineControllerAPI - Interface wrapper for the Pipeline Controller API
type PipelineControllerAPI interface {
	SavePipelineUsingPOST(ctx context.Context, pipeline interface{}, localVarOptionals *spinGateApi.PipelineControllerApiSavePipelineUsingPOSTOpts) (*http.Response, error)
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
	log *log.Logger
}

// NewClient - Create a new default spinnaker client
func NewClient(logger *log.Logger) *SpinClient {
	return &SpinClient{log: logger}
}

// initializeAPI - Lazy initialization of the client, is expected to be called before each method that requires http.
// Concept taken from: https://roberto.selbach.ca/zero-values-in-go-and-lazy-initialization/
func (s *SpinClient) initializeAPI() error {
	var outerErr error
	// If the client is already initialized, not
	if s.SpinCLI == nil && s.CustomSpinCLI == nil {
		s.initOnce.Do(func() {
			gateClient, err := spinGate.NewGateClient(&UI{}, "", "", "", true)

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

			s.CustomSpinCLI = &CustomSpinClient{Endpoint: gateClient.Config.Gate.Endpoint, HTTPClient: *httpClient}
			s.SpinCLI = &SpinCLI{
				ApplicationControllerAPI: gateClient.ApplicationControllerApi,
				PipelineControllerAPI:    gateClient.PipelineControllerApi,
				Context:                  gateClient.Context,
			}
		})
	}

	return outerErr
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

	foundPipeline, queryResp, _ := s.ApplicationControllerAPI.GetPipelineConfigUsingGET(s.Context, application, pipelineName)
	if queryResp.StatusCode == http.StatusOK {
		// pipeline found, let's use Spinnaker's known Pipeline ID, otherwise we'll get one created for us
		if len(foundPipeline) > 0 {
			s.log.Info("Pipeline ", foundPipeline["name"], " found with id ", foundPipeline["id"], " in application ", application)

			pipeline["id"] = foundPipeline["id"].(string)
			pipelineID = foundPipeline["id"].(string)
		}

	} else if queryResp.StatusCode == http.StatusNotFound {
		// pipeline doesn't exists, let's create a new one
		s.log.Info("Pipeline ", pipelineName, "not found in application ", application)
	} else {
		b, _ := ioutil.ReadAll(queryResp.Body)
		return pipelineID, nil, fmt.Errorf("unhandled response %d: %s", queryResp.StatusCode, b)
	}
	opt := &spinGateApi.PipelineControllerApiSavePipelineUsingPOSTOpts{}
	res, err := s.PipelineControllerAPI.SavePipelineUsingPOST(s.Context, pipeline, opt)

	return pipelineID, res, err
}

// ExecutePipeline - Execute a spinnaker pipeline.
//
// Patameters are optional.
func (s *SpinClient) ExecutePipeline(argsJSON string) (interface{}, *http.Response, error) {
	// For some crazy reason, spicli invoke doesn't return the ID of the pipeline execution.
	// BTW the crazy reason is that `swagger-code-gen` produces wrong code and Spin-Cli (and shore...) depends on this wrong code.
	// So this request needs to be done 100% manually.
	// For this reason we use the `CustomSpinCLI` interface to implement all the things `SpinCli` either does wrong, or is 100% broken.

	var args map[string]interface{}

	if err := s.initializeAPI(); err != nil {
		return &ExecutePipelineResponse{}, &http.Response{}, err
	}

	err := json.Unmarshal([]byte(argsJSON), &args)

	if err != nil {
		return &ExecutePipelineResponse{}, &http.Response{}, err
	}

	if err = s.isArgsValid(args); err != nil {
		return &ExecutePipelineResponse{}, &http.Response{}, err
	}

	application := args["application"].(string)
	pipelineName := args["pipeline"].(string)

	delete(args, "application")
	delete(args, "pipeline")

	argsBytes, err := json.Marshal(args)

	if err != nil {
		return &ExecutePipelineResponse{}, &http.Response{}, err
	}

	body, res, err := s.CustomSpinCLI.ExecutePipeline(application, pipelineName, bytes.NewBuffer(argsBytes))
	return body, res, err
}

// TestPipeline - Run a Spinnaker testing
// Currently returns a not-so-well formatted error.
// The indended solution is to create a shared API between `shore-cli` & the `backend` to expect well formatted struct for the CLI to render correctly.
// TODO: Design a struct to pass data back to `shore-cli` so the UI layer could render the test-results correctly.
func (s *SpinClient) TestPipeline(config string, onChange func()) error {
	var testConfig TestsConfig

	if err := jsoniter.Unmarshal([]byte(config), &testConfig); err != nil {
		return err
	}

	if err := s.initializeAPI(); err != nil {
		return err
	}

	if testConfig.Application == "" {
		return fmt.Errorf("test config missing required property `application`")
	}

	if testConfig.Pipeline == "" {
		return fmt.Errorf("test config missing required property `pipeline`")
	}

	testErrors := make(map[string][]string)

	for testName, test := range testConfig.Tests {
		execArgs, err := json.Marshal(test.ExecArgs)

		if err != nil {
			return err
		}

		body, _, err := s.CustomSpinCLI.ExecutePipeline(testConfig.Application, testConfig.Pipeline, bytes.NewBuffer(execArgs))

		if err != nil {
			testErrors[testName] = append(testErrors[testName], err.Error())
		}

		refID := strings.Split(body.Ref, "/")[2]

		var execDetails *PipelineExecutionDetailsResponse

		execDetails, _, err = s.CustomSpinCLI.PipelineExecutionDetails(refID, bytes.NewBuffer(make([]byte, 0)))

		// Currently waiting for roughly 20 minutes before exiting the loop
		// TODO: Make this value configurable
		maxTries := 50
		tries := 0

		for execDetails.Status == PipelineRunning && tries < maxTries {
			execDetails, _, err = s.CustomSpinCLI.PipelineExecutionDetails(refID, bytes.NewBuffer(make([]byte, 0)))
			tries++
			time.Sleep(time.Second * time.Duration(tries))
		}

		if tries == maxTries {
			testErrors[testName] = append(testErrors[testName], fmt.Sprintf("max timed out reached for test: '%s'", testName))
			continue
		}

		if err != nil {
			testErrors[testName] = append(testErrors[testName], err.Error())
		}

		for _, stage := range execDetails.Stages {
			stageName := stage["name"].(string)
			assetion, exists := test.Assertions[stageName]

			if !exists {
				testErrors[testName] = append(testErrors[testName], fmt.Sprintf("missing assertion for stage %s", stageName))
				continue
			}

			expectedStatus := strings.ToUpper(assetion.ExpectedStatus)

			if err := isExpectedStatus(expectedStatus, stage["status"].(string), stageName); err != nil {
				testErrors[testName] = append(testErrors[testName], err.Error())
			}

			if err := isExpectedOutput(assetion.ExpectedOutput, stage["outputs"].(map[string]interface{}), stageName); err != nil {
				testErrors[testName] = append(testErrors[testName], err.Error())
			}
		}
	}

	// This should really be handled by a Golang template.
	// TODO: move this logic to the CLI layer.
	if len(testErrors) > 0 {
		formmatedErrors := ""

		for testName, testErrors := range testErrors {
			formmatedErrors += fmt.Sprintf("`%s` failure:\n", testName)
			for _, testError := range testErrors {
				formmatedErrors += fmt.Sprintf("%s\n", testError)
			}
			formmatedErrors += "\n"
		}

		// TODO: The backend shouldn't concern itself with rendering, this should be replaced with the correct struct response.
		return errors.New(formmatedErrors)
	}

	return nil
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
		if _, exists := stage.(map[string]interface{})["name"]; !exists {
			errorsList = append(errorsList, "required stage key 'name' missing for stage")
		}

		// if it is a pipeline stage it has to have "application" key.
		// Application value should match the one with the parent pipeline's one
		if _, exists := stage.(map[string]interface{})["pipeline"]; exists {
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
		foundPipeline, queryResp, _ := s.ApplicationControllerAPI.GetPipelineConfigUsingGET(s.Context, application, pipelineName)

		if queryResp.StatusCode != http.StatusOK {
			b, _ := ioutil.ReadAll(queryResp.Body)
			err := fmt.Errorf("response %d: %s. nested pipeline '%s' wasn't created in application '%s', pipeline Stage will be unbound",
				queryResp.StatusCode,
				b,
				pipelineName,
				application)
			return "", queryResp, err
		}

		if len(foundPipeline) == 0 {
			log.Println("get pipeline request didn't return a payload, sleeping for", pollingSleep)
			time.Sleep(time.Duration(pollingSleep))
			continue
		}

		log.Println("Pipeline", foundPipeline["name"], "found with id", foundPipeline["id"], "in application", application)

		return foundPipeline["id"].(string), queryResp, nil
	}

	return "", &http.Response{}, fmt.Errorf("Couldn't get pipeline until hitting timeout")
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
		stagePipelineField, exists := stage.(map[string]interface{})["pipeline"]
		if !exists {
			continue
		}

		switch stagePipelineField.(type) {
		case string:
			continue
		}

		childPipeline := stage.(map[string]interface{})["pipeline"].(map[string]interface{})

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
		log.Println(res)

		// Do not try to poll for pipeline ID again if exists already
		if pipelineID == "" {
			pipelineID, res, err = s.pollSpinnakerGetPipelineConfigUsingGET(childPipeline["application"].(string), childPipeline["name"].(string))
			if err != nil {
				return err
			}
			// pipelineID = pipelineID

			log.Println("Created new pipeline with id:", pipelineID)
			log.Println(res)
		}

		// And override stage pipeline value with an the id (a UUID String) received from spinnaker pipeline.
		// maps are updated by reference so the save of the parent pipeline will be updated as well
		stage.(map[string]interface{})["pipeline"] = pipelineID
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

	// Check whether a pipeline has stages list
	stages, exists := pipeline["stages"]

	if exists {
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
