package retry_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Autodesk/shore/internal/retry"
	"github.com/stretchr/testify/assert"
)

var quickConfig retry.Config = retry.Config{
	Tries: 1,
	Delay: time.Nanosecond * 1,
}

func TestRetryNoErrorConfig(t *testing.T) {
	// Given
	tries := 0
	// Test
	err := retry.Retry(func() error {
		tries += 1
		return nil
	}, quickConfig)

	// Assert
	assert.Equal(t, quickConfig.Tries, tries)
	assert.Nil(t, err)
}

func TestRetryRetriesExceeded(t *testing.T) {
	// Given
	// Test
	err := retry.Retry(func() error {
		return retry.ErrRetry
	}, quickConfig)

	assert.IsType(t, &retry.ErrMaxRetries{}, err)
}

func TestDelayFunc(t *testing.T) {
	// Given
	config := retry.Config{
		Tries:     1,
		DelayFunc: func(try int) time.Duration { return time.Nanosecond * time.Duration(try) },
	}

	// Test
	err := retry.Retry(func() error {
		return retry.ErrRetry
	}, config)

	assert.Error(t, err, "max retries exceeded: %w")
}

type CustomErr struct{}

func (e *CustomErr) Error() string {
	return "custom error"
}

func TestRetryCustomErrBreak(t *testing.T) {
	// Given
	config := retry.Config{
		Tries: 3,
		Delay: time.Nanosecond * time.Duration(1),
	}

	tries := 0
	// Test
	err := retry.Retry(func() error {
		tries += 1
		if tries > 2 {
			return &CustomErr{}
		}

		return retry.ErrRetry
	}, config)

	assert.IsType(t, &CustomErr{}, err)
}

func TestRetryErrMissingDelays(t *testing.T) {
	// Given
	config := retry.Config{
		Tries: 3,
	}

	// Test
	err := retry.Retry(func() error { return nil }, config)

	assert.EqualError(t, err, "Config must include EITHER a non zero Delay, or a non nil DelayFunc")
}

func TestRetryErrMissingNilDelays(t *testing.T) {
	// Given
	config := retry.Config{
		Tries:     3,
		Delay:     0,
		DelayFunc: nil,
	}

	// Test
	err := retry.Retry(func() error { return nil }, config)

	assert.EqualError(t, err, "Config must include EITHER a non zero Delay, or a non nil DelayFunc")
}

func ExampleRetry() {
	// Store results from the closure into outer-scope variables.
	var res http.Response
	var err error

	// Retry Config
	config := retry.Config{
		// Retry 10 times
		Tries: 10,
		// Wait for 1 Nanosecond between each retry
		Delay: time.Nanosecond * time.Duration(1),
	}

	// Keeping track of retries for testing purposes only...
	// THIS IS NOT NEEDED (OR RECOMMENDED) FOR PRODUCTION CODE
	tries := 0

	httpCall := func() (http.Response, error) {
		tries += 1

		if tries > 9 {
			return http.Response{Body: ioutil.NopCloser(strings.NewReader("I lived!")), Status: "200"}, nil
		}

		return http.Response{}, fmt.Errorf("I died...")
	}

	retryErr := retry.Retry(func() error {
		// Retry some request
		// Fake response
		res, err = httpCall()

		if err != nil {
			return retry.ErrRetry
		}

		return nil
	}, config)

	if retryErr != nil {
		fmt.Fprintf(os.Stderr, "%s", retryErr)
	}

	fmt.Println(res, err)
}

func ExampleConfig_DelayFunc() {
	// Store results from the closure into outer-scope variables.
	var res http.Response
	var err error

	// Retry Config
	config := retry.Config{
		// Retry 10 times
		Tries: 10,
		// Linear regression instead of a simple delay.
		DelayFunc: func(try int) time.Duration { return time.Nanosecond * time.Duration(try) },
	}

	// Keeping track of retries for testing purposes only...
	// THIS IS NOT NEEDED (OR RECOMMENDED) FOR PRODUCTION CODE
	tries := 0

	fn := func() error {
		// Retry some request
		// Fake response
		tries += 1

		if tries > 9 {
			res = http.Response{}
			err = fmt.Errorf("I died...")
		} else {
			res = http.Response{Body: ioutil.NopCloser(strings.NewReader("I lived!")), Status: "200"}
			err = fmt.Errorf("I died...")
		}

		if err != nil {
			return retry.ErrRetry
		}

		return nil
	}

	if retryErr := retry.Retry(fn, config); retryErr != nil {
		fmt.Fprintf(os.Stderr, "%s", retryErr)
	}

	fmt.Println(res, err)
}
