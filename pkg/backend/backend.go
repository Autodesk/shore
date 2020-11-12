package backend

import "net/http"

// Backend - an interface that describes a generic backend pipeline
type Backend interface {
	// TODO: Return type needs to be a custom wrapper in the future.
	// We cannot assume that every backend-cli implementation will return an HTTP object.
	SavePipeline(pipelineJSON string) (*http.Response, error)
	ExecutePipeline(parameters string) (interface{}, *http.Response, error)
	// TODO: Reconsider `onChange`, it may be a channel to communicate data between `shore-cli` & the Testing process in an async fashion.
	TestPipeline(testConfig string, onChange func()) error
}
