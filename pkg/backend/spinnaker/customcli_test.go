package spinnaker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockedHTTPClientOK struct{}

func (c *mockedHTTPClientOK) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(`{"Ref": "/pipeline/1234"}`)),
	}, nil
}

type mockedHTTPClientUnmarshalError struct{}

func (c *mockedHTTPClientUnmarshalError) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(`{"Ref": "/pipelin`)),
	}, nil
}

type mockedHTTPClient404 struct{}

func (c *mockedHTTPClient404) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

type mockedHTTPClientAPIError struct{}

func (c *mockedHTTPClientAPIError) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{}, errors.New(fmt.Sprint("API error"))
}

func TestExecutePipelineOKStatus(t *testing.T) {
	// Test
	application := "TestApplication"
	pipelineName := "Some pipeline"
	expectedResponseBody := &ExecutePipelineResponse{Ref: "/pipeline/1234"}

	httpClient := &mockedHTTPClientOK{}
	mockedCustomSpinCLI := &CustomSpinClient{Endpoint: "none", HTTPClient: httpClient}

	body, res, err := mockedCustomSpinCLI.ExecutePipeline(application, pipelineName, strings.NewReader(""))

	assert.Equal(t, body, expectedResponseBody)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestExecutePipeline404Status(t *testing.T) {
	// Test
	application := "TestApplication"
	pipelineName := "Some pipeline"
	expectedResponseBody := &ExecutePipelineResponse{}
	expectedError := &CustomCliError{
		PipelineName:    pipelineName,
		ApplicationName: application,
		StatusCode:      404,
		Err:             error(nil),
	}
	httpClient := &mockedHTTPClient404{}
	mockedCustomSpinCLI := &CustomSpinClient{Endpoint: "none", HTTPClient: httpClient}

	body, res, err := mockedCustomSpinCLI.ExecutePipeline(application, pipelineName, strings.NewReader(""))

	assert.Equal(t, body, expectedResponseBody)
	assert.Equal(t, err, expectedError)
	assert.Equal(t, res.StatusCode, http.StatusNotFound)
}

func TestExecutePipelineUnmarshalError(t *testing.T) {
	// Test
	application := "TestApplication"
	pipelineName := "Some pipeline"

	httpClient := &mockedHTTPClientUnmarshalError{}
	mockedCustomSpinCLI := &CustomSpinClient{Endpoint: "none", HTTPClient: httpClient}

	_, res, err := mockedCustomSpinCLI.ExecutePipeline(application, pipelineName, strings.NewReader(""))

	assert.NotNil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestExecutePipelineAPIError(t *testing.T) {
	// Test
	application := "TestApplication"
	pipelineName := "Some pipeline"
	expectedError := &CustomCliError{
		PipelineName:    pipelineName,
		ApplicationName: application,
		StatusCode:      0,
		Err:             errors.New(fmt.Sprint("API error")),
	}
	httpClient := &mockedHTTPClientAPIError{}
	mockedCustomSpinCLI := &CustomSpinClient{Endpoint: "none", HTTPClient: httpClient}

	_, res, err := mockedCustomSpinCLI.ExecutePipeline(application, pipelineName, strings.NewReader(""))

	assert.Equal(t, err, expectedError)
	assert.Equal(t, res.StatusCode, 0)
}
