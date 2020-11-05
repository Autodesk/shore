package backend

import "net/http"

// Backend - an interface that describes a generic backend pipeline
type Backend interface {
	// TODO: Return type needs to be a custom wrapper in the future.
	// We cannot assume that every backend-cli implementation will return an HTTP object.
	SavePipeline(pipelineJSON string) (*http.Response, error)
	// ExecutePipeline(application, pipelineName string, parameters ...interface{}) (*http.Response, error)
}
