package controller

import (
	"log"
	"net/http"

	"github.com/Autodesk/shore/pkg/backend/spinnaker"
)

// Creates a new Client and saves a pipeline.
func SavePipeline(pipeline string) (*http.Response, error) {
	// TODO: The backend client (AKA service provider) should either be DI'ed or imported from a "global" context.
	// We don't know which cli the customer may choose to use in the future.
	cli, err := spinnaker.NewClient()

	if err != nil {
		log.Fatal(err)
	}

	return cli.SavePipeline(pipeline)
}
