package spinnaker

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	spinGateApi "github.com/spinnaker/spin/gateapi"
)

type ApplicationControllerError struct {
	body     string
	err      error
	response *http.Response
}

// NewApplicationControllerError returns an error which wraps an error and HTTP
// response from gateapi's ApplicationController. This provides a consistent
// interface for reading the response body, without having to deal with the
// quirky ways Spinnaker wraps errors and manages response body lifetimes.
func NewApplicationControllerError(spinnakerErr error, response *http.Response) *ApplicationControllerError {
	// Sometimes gateapi will wrap the error in a `GenericSwaggerError`, in
	// which case it will have read and closed the response body...
	var swaggerErr spinGateApi.GenericSwaggerError
	if errors.As(spinnakerErr, &swaggerErr) {
		return &ApplicationControllerError{
			body:     string(swaggerErr.Body()),
			err:      spinnakerErr,
			response: response,
		}
	}

	// ... and sometimes gateapi will just return the error from
	// `http.Client.Do`, without modifying the response or error.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		// This typically means the response body was closed prematurely,
		// to be defensive we assume that the caller read the response body,
		// rather than the gateapi client, and use an empty string for the body.
		return &ApplicationControllerError{
			body:     "",
			err:      spinnakerErr,
			response: response,
		}
	}
	defer response.Body.Close()

	return &ApplicationControllerError{
		body:     string(body),
		err:      spinnakerErr,
		response: response,
	}
}

func (e ApplicationControllerError) Error() string {
	return fmt.Sprintf("response code %d: %q", e.response.StatusCode, e.body)
}

func (e ApplicationControllerError) StatusCode() int {
	return e.response.StatusCode
}
