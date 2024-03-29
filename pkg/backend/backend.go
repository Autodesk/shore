package backend

import (
	"net/http"

	"github.com/Autodesk/shore/pkg/shore_testing"
)

// Backend - an interface that describes a generic backend pipeline
type Backend interface {
	// TODO: Return type needs to be a custom wrapper in the future.
	// We cannot assume that every backend-cli implementation will return an HTTP object.
	SavePipeline(pipelineJSON string) (*http.Response, error)
	ExecutePipeline(parameters string, stringify bool) (string, *http.Response, error)
	WaitForPipelineToFinish(id string, timeout int) (string, *http.Response, error)
	// TODO: Reconsider `onChange`, it may be a channel to communicate data between `shore-cli` & the Testing process in an async fashion.
	TestPipeline(testConfig shore_testing.TestsConfig, onChange func(), stringify bool) error
	GetPipeline(application string, pipelineName string) (map[string]interface{}, *http.Response, error)
	DeletePipeline(pipelineJSON string) (*http.Response, error)
	GetPipelinesNamesAndApplication(pipelineJSON string) ([]string, string, error)
}
