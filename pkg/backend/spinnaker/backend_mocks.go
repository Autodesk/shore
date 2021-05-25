package spinnaker

import (
	"context"
	"io"
	"net/http"

	spinGateApi "github.com/spinnaker/spin/gateapi"
)

type MockCustomSpinCli struct {
	CustomSpinCLI
}
type MockPipelineControllerAPI struct{}
type MockApplicationControllerAPI struct{}
type MockApplicationControllerAPIWithEmptyID struct{}

func (a *MockApplicationControllerAPI) GetPipelineConfigUsingGET(ctx context.Context, application string, pipelineName string) (map[string]interface{}, *http.Response, error) {
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

func (p *MockPipelineControllerAPI) SavePipelineUsingPOST(ctx context.Context, pipeline interface{}, localVarOptionals *spinGateApi.PipelineControllerApiSavePipelineUsingPOSTOpts) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK}, nil
}

func (p *MockPipelineControllerAPI) InvokePipelineConfigUsingPOST1(ctx context.Context, application string, pipelineNameOrID string, localVarOptionals *spinGateApi.PipelineControllerApiInvokePipelineConfigUsingPOST1Opts) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK}, nil
}

func (s *MockCustomSpinCli) ExecutePipeline(application, pipelineName string, args io.Reader) (*ExecutePipelineResponse, *http.Response, error) {
	req, _ := http.NewRequest("POST", "url", args)
	return &ExecutePipelineResponse{Ref: "/pipeline/1234"}, &http.Response{StatusCode: http.StatusOK, Request: req}, nil
}

func (s *MockCustomSpinCli) PipelineExecutionDetails(refID string, args io.Reader) (*PipelineExecutionDetailsResponse, *http.Response, error) {
	return &PipelineExecutionDetailsResponse{
		Application: "application",
		Stages: []map[string]interface{}{
			{
				"name":   "testedname",
				"status": "SUCCEEDED",
				"outputs": map[string]interface{}{
					"test": "123",
				},
			},
		},
		PipelineName: "pipeline",
	}, &http.Response{StatusCode: http.StatusOK}, nil
}

func (a *MockApplicationControllerAPIWithEmptyID) GetPipelineConfigUsingGET(ctx context.Context, application string, pipelineName string) (map[string]interface{}, *http.Response, error) {
	res := map[string]interface{}{}

	return res, &http.Response{StatusCode: http.StatusOK}, nil
}
