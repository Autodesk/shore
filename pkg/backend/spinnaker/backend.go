package spinnaker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

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
	// InvokePipelineConfigUsingPOST1(ctx context.Context, application string, pipelineNameOrId string, localVarOptionals *spinGateApi.PipelineControllerApiInvokePipelineConfigUsingPOST1Opts) (*http.Response, error)
}

// SpinClient - Concrete type requiring all the methods of the specified interfaces.
type SpinClient struct {
	// *spinGate.GatewayClient
	initOnce sync.Once
	ApplicationControllerAPI
	PipelineControllerAPI
	context.Context
}

// NewClient - Create a new spinnaker client
func NewClient() *SpinClient {
	return &SpinClient{}
}

// initalizeClient - Lazy initialization of the client, is expected to be called before each method that requires http.
// Concept taken from: https://roberto.selbach.ca/zero-values-in-go-and-lazy-initialization/
func (s *SpinClient) initalizeClient() error {
	var outerErr error
	// If the client is already initialized, not
	if s.ApplicationControllerAPI == nil && s.PipelineControllerAPI == nil && s.Context == nil {
		s.initOnce.Do(func() {
			gateClient, err := spinGate.NewGateClient(&UI{}, "", "", "", true)

			if err != nil {
				outerErr = err
				return
			}

			s.ApplicationControllerAPI = gateClient.ApplicationControllerApi
			s.PipelineControllerAPI = gateClient.PipelineControllerApi
			s.Context = gateClient.Context
		})
	}

	return outerErr
}

// SavePipeline - Create or Update a pipeline.
func (s *SpinClient) SavePipeline(pipelineJSON string) (*http.Response, error) {
	if err := s.initalizeClient(); err != nil {
		return &http.Response{}, err
	}

	var pipeline map[string]interface{}

	err := json.Unmarshal([]byte(pipelineJSON), &pipeline)

	if err != nil {
		log.Fatal(err)
	}

	if err = s.isValidPipeline(pipeline); err != nil {
		return &http.Response{}, err
	}

	if template, exists := pipeline["template"]; exists && len(template.(map[string]interface{})) > 0 {
		if _, exists := pipeline["schema"]; !exists {
			return &http.Response{}, fmt.Errorf("required pipeline key 'schema' missing for templated pipeline")
		}
		pipeline["type"] = "templatedPipeline"
	}

	application := pipeline["application"].(string)
	pipelineName := pipeline["name"].(string)

	foundPipeline, queryResp, _ := s.ApplicationControllerAPI.GetPipelineConfigUsingGET(s.Context, application, pipelineName)
	if queryResp.StatusCode == http.StatusOK {
		// pipeline found, let's use Spinnaker's known Pipeline ID, otherwise we'll get one created for us
		if len(foundPipeline) > 0 {
			pipeline["id"] = foundPipeline["id"].(string)
		}
	} else if queryResp.StatusCode == http.StatusNotFound {
		// pipeline doesn't exists, let's create a new one
	} else {
		b, _ := ioutil.ReadAll(queryResp.Body)
		return nil, fmt.Errorf("unhandled response %d: %s", queryResp.StatusCode, b)
	}

	opt := &spinGateApi.PipelineControllerApiSavePipelineUsingPOSTOpts{}
	return s.PipelineControllerAPI.SavePipelineUsingPOST(s.Context, pipeline, opt)
}

// ExecutePipeline - Execute a spinnaker pipeline.
//
// Patameters are optional.
// func (s *SpinClient) ExecutePipeline(application, pipelineName string, parameters ...interface{}) (*http.Response, error) {
// 	if err := s.initalizeClient(); err != nil {
// 		return &http.Response{}, err
// 	}

// 	var opts *spinGateApi.PipelineControllerApiInvokePipelineConfigUsingPOST1Opts

// 	if len(parameters) > 0 {
// 		opts = &spinGateApi.PipelineControllerApiInvokePipelineConfigUsingPOST1Opts{
// 			Trigger: optional.NewInterface(parameters[0]),
// 		}
// 	} else {
// 		opts = &spinGateApi.PipelineControllerApiInvokePipelineConfigUsingPOST1Opts{}
// 	}

// 	return s.PipelineControllerAPI.InvokePipelineConfigUsingPOST1(s.Context, application, pipelineName, opts)
// }

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
