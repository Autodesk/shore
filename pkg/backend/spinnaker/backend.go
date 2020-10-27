package spinnaker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	spingGate "github.com/spinnaker/spin/cmd/gateclient"
	spingGateApi "github.com/spinnaker/spin/gateapi"
)

type SpinClient struct {
	spingGate.GatewayClient
}

// NewClient - Create a new spinnaker client
func NewClient() (*SpinClient, error) {
	gateClient, err := spingGate.NewGateClient(nil, "", "", "", true)

	if err != nil {
		return &SpinClient{}, err
	}

	return &SpinClient{GatewayClient: *gateClient}, nil
}

// SavePipeline - Create or Update a pipeline.
func (s *SpinClient) SavePipeline(pipelineJSON string) (*http.Response, error) {
	var pipeline map[string]interface{}
	var errorsList []string
	valid := true

	err := json.Unmarshal([]byte(pipelineJSON), &pipeline)

	if err != nil {
		log.Fatal(err)
	}

	if _, exists := pipeline["name"]; !exists {
		errorsList = append(errorsList, "required pipeline key 'name' missing")
		valid = false
	}

	if _, exists := pipeline["application"]; !exists {
		errorsList = append(errorsList, "Required pipeline key 'application' missing")
		valid = false
	}

	if template, exists := pipeline["template"]; exists && len(template.(map[string]interface{})) > 0 {
		if _, exists := pipeline["schema"]; !exists {
			errorsList = append(errorsList, "Required pipeline key 'schema' missing for templated pipeline")
			valid = false
		}
		pipeline["type"] = "templatedPipeline"
	}

	if !valid {
		return nil, fmt.Errorf(strings.Join(errorsList, "\n"))
	}

	application := pipeline["application"].(string)
	pipelineName := pipeline["name"].(string)

	foundPipeline, queryResp, _ := s.ApplicationControllerApi.GetPipelineConfigUsingGET(s.Context, application, pipelineName)
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

	opt := &spingGateApi.PipelineControllerApiSavePipelineUsingPOSTOpts{}
	return s.PipelineControllerApi.SavePipelineUsingPOST(s.Context, pipelineJSON, opt)
}
