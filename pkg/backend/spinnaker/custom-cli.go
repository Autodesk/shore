package spinnaker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

/*
A custom API layer implementation for specific usecases when gateclient doesn't provide a good interface
*/

// CustomSpinCLI is a wrapper the implementes specific requests that are either broken or unsupported by SpinCLI.
type CustomSpinCLI interface {
	Do(req *http.Request) ([]byte, *http.Response, error)
	Post(url string, args *bytes.Buffer) ([]byte, *http.Response, error)
	Get(url string, args *bytes.Buffer) ([]byte, *http.Response, error)
	ExecutePipeline(application, pipelineName string, args *bytes.Buffer) (*ExecutePipelineResponse, *http.Response, error)
	PipelineExecutionDetails(refID string, args *bytes.Buffer) (*PipelineExecutionDetailsResponse, *http.Response, error)
}

// CustomSpinClient is a wrapper the implementes specific requests that are either broken or unsupported by SpinCLI.
type CustomSpinClient struct {
	CustomSpinCLI
	Endpoint   string
	HTTPClient http.Client
}

func (cli *CustomSpinClient) Do(req *http.Request) ([]byte, *http.Response, error) {
	req.Header.Set("Content-Type", "application/json")

	res, err := cli.HTTPClient.Do(req)

	if err != nil {
		return nil, res, err
	}

	defer res.Body.Close()

	resBuf, err := ioutil.ReadAll(res.Body)

	return resBuf, res, err
}

// Generic post request using the Spinnaker HTTP Client
func (cli *CustomSpinClient) Post(url string, args *bytes.Buffer) ([]byte, *http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, args)

	if err != nil {
		return nil, &http.Response{}, err
	}

	return cli.Do(req)
}

func (cli *CustomSpinClient) Get(url string, args *bytes.Buffer) ([]byte, *http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, args)

	if err != nil {
		return nil, &http.Response{}, err
	}

	return cli.Do(req)
}

type ExecutePipelineResponse struct {
	Ref string
}

// ExecutePipeline calls the POST endpoint of a pipeline to execute it
func (cli *CustomSpinClient) ExecutePipeline(application, pipelineName string, args *bytes.Buffer) (*ExecutePipelineResponse, *http.Response, error) {
	var execPipelineResponse ExecutePipelineResponse

	url := fmt.Sprintf("%s/pipelines/%s/%s", cli.Endpoint, application, pipelineName)
	body, res, nil := cli.Post(url, args)
	err := jsoniter.Unmarshal(body, &execPipelineResponse)

	if err != nil {
		return &ExecutePipelineResponse{}, res, err
	}

	return &execPipelineResponse, res, nil
}

type PipelineExecutionDetailsResponse struct {
	Type         string                   `json:"type"`
	Status       string                   `json:"status"`
	Canceled     bool                     `json:"canceled"`
	BuildTime    int                      `json:"buildTime"`
	StartTime    int                      `json:"startTime"`
	Application  string                   `json:"application"`
	Stages       []map[string]interface{} `json:"stages"`
	PipelineName string                   `json:"pipelineName"`
}

func (cli *CustomSpinClient) PipelineExecutionDetails(refID string, args *bytes.Buffer) (*PipelineExecutionDetailsResponse, *http.Response, error) {
	var pipelineExecutionDetails PipelineExecutionDetailsResponse

	url := fmt.Sprintf("%s/pipelines/%s/", cli.Endpoint, refID)
	body, res, nil := cli.Get(url, args)
	err := jsoniter.Unmarshal(body, &pipelineExecutionDetails)

	if err != nil {
		return &PipelineExecutionDetailsResponse{}, res, err
	}

	return &pipelineExecutionDetails, res, nil
}
