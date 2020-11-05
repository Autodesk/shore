package spinnaker

import (
	"context"
	"net/http"
	"testing"

	spinGateApi "github.com/spinnaker/spin/gateapi"
	"github.com/stretchr/testify/assert"
)

type mockApplicationControllerAPI struct{}

func (a *mockApplicationControllerAPI) GetPipelineConfigUsingGET(ctx context.Context, application string, pipelineName string) (map[string]interface{}, *http.Response, error) {
	res := map[string]interface{}{
		"id": "1234",
	}

	return res, &http.Response{StatusCode: http.StatusOK}, nil
}

type mockPipelineControllerAPI struct{}

func (p *mockPipelineControllerAPI) SavePipelineUsingPOST(ctx context.Context, pipeline interface{}, localVarOptionals *spinGateApi.PipelineControllerApiSavePipelineUsingPOSTOpts) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK}, nil
}

type mockSpinRequester struct{}

func (s *mockSpinRequester) ChangeBasePath(path string) {}

func TestSaveSuccess(t *testing.T) {
	// Given
	cli := &SpinClient{
		ApplicationControllerAPI: &mockApplicationControllerAPI{},
		PipelineControllerAPI:    &mockPipelineControllerAPI{},
		Context:                  context.Background(),
	}

	// Test
	res, err := cli.SavePipeline(`{"application": "test", "name": "test"}`)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestSaveFailedApplication(t *testing.T) {
	// Given
	cli := &SpinClient{
		ApplicationControllerAPI: &mockApplicationControllerAPI{},
		PipelineControllerAPI:    &mockPipelineControllerAPI{},
		Context:                  context.Background(),
	}
	// Test
	_, err := cli.SavePipeline(`{"name": "test"}`)

	// Assert
	assert.EqualError(t, err, "required pipeline key 'application' missing")
}

func TestSaveFailedName(t *testing.T) {
	// Given
	cli := &SpinClient{
		ApplicationControllerAPI: &mockApplicationControllerAPI{},
		PipelineControllerAPI:    &mockPipelineControllerAPI{},
		Context:                  context.Background(),
	}
	// Test
	_, err := cli.SavePipeline(`{"application": "test"}`)

	// Assert
	assert.EqualError(t, err, "required pipeline key 'name' missing")
}
