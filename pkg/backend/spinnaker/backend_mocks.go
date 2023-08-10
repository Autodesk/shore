package spinnaker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	spinGateApi "github.com/spinnaker/spin/gateapi"
	"github.com/stretchr/testify/mock"
)

type MockCustomSpinCli struct {
	CustomSpinCLI
}
type MockPipelineControllerAPI struct {
	mock.Mock
}
type MockApplicationControllerAPI struct{}
type MockApplicationControllerAPIWithEmptyID struct{}

func (a *MockApplicationControllerAPI) GetPipelineConfigUsingGET(ctx context.Context, application string, pipelineName string) (map[string]interface{}, *http.Response, error) {
	var res map[string]interface{}

	if application == "not-exists" {
		res = map[string]interface{}{}
	} else {
		res = map[string]interface{}{
			"name": pipelineName,
			"id":   "1234",
		}
	}

	return res, &http.Response{StatusCode: http.StatusOK}, nil
}

func (p *MockPipelineControllerAPI) SavePipelineUsingPOST(ctx context.Context, pipeline interface{}, localVarOptionals *spinGateApi.PipelineControllerApiSavePipelineUsingPOSTOpts) (*http.Response, error) {
	data, err := jsoniter.Marshal(pipeline)

	if err != nil {
		return &http.Response{}, err
	}

	body := ioutil.NopCloser(bytes.NewReader(data))

	return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
}

func (p *MockPipelineControllerAPI) DeletePipelineUsingDELETE(ctx context.Context, application string, pipeline string) (*http.Response, error) {
	data, err := jsoniter.Marshal(pipeline)

	if err != nil {
		return &http.Response{}, err
	}

	body := ioutil.NopCloser(bytes.NewReader(data))

	return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
}

func (p *MockPipelineControllerAPI) InvokePipelineConfigUsingPOST1(ctx context.Context, application string, pipelineNameOrID string, localVarOptionals *spinGateApi.PipelineControllerApiInvokePipelineConfigUsingPOST1Opts) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK}, nil
}

func (s *MockCustomSpinCli) ExecutePipeline(application, pipelineName string, args io.Reader) (*ExecutePipelineResponse, *http.Response, error) {
	req, _ := http.NewRequest("POST", "url", args)
	refID := -1

	switch application {
	case "timeout-app":
		refID = 9999
	default:
		refID = 1234
	}

	return &ExecutePipelineResponse{Ref: fmt.Sprintf("/pipeline/%d", refID)}, &http.Response{StatusCode: http.StatusOK, Request: req}, nil
}

func (s *MockCustomSpinCli) PipelineExecutionDetails(refID string, args io.Reader) (*PipelineExecutionDetailsResponse, *http.Response, error) {
	status := "SUCCEEDED"
	if refID == "9999" { // timeout-app
		status = "NOT_STARTED"
	}
	return &PipelineExecutionDetailsResponse{
		Application: "application",
		Stages: []map[string]interface{}{
			{
				"name":   "testedname",
				"status": status,
				"outputs": map[string]interface{}{
					"test": "123",
				},
			},
		},
		Status:       status,
		PipelineName: "pipeline",
	}, &http.Response{StatusCode: http.StatusOK}, nil
}

func (a *MockApplicationControllerAPIWithEmptyID) GetPipelineConfigUsingGET(ctx context.Context, application string, pipelineName string) (map[string]interface{}, *http.Response, error) {
	res := map[string]interface{}{}

	return res, &http.Response{StatusCode: http.StatusOK}, nil
}
