/**
Package retry

This package implements a common retry mechanism for easily implementing retries.
*/
package retry

import (
	"errors"
	"fmt"
	"time"
)

// Config - config parameters for the retry function
type Config struct {
	Tries     int
	Delay     time.Duration
	DelayFunc func(try int) time.Duration
}

// Func - retry function type def
type Func func() error

// ErrRetry - custom retry error to return when the process should be retried
var ErrRetry = errors.New("Retry Error")

// ErrMaxRetries - custom error thrown when retries were exceeded.
type ErrMaxRetries struct {
	// The last run error encountered when max retries was reached
	RunErr error
}

func (e *ErrMaxRetries) Error() string {
	return "max retries exceeded"
}

// RetryDefaultConfig - Default config to use when no config is supplied
var retryDefaultConfig = Config{
	Tries: 5,
	Delay: time.Second * 6,
}

/*
Retry - closure runner for retries

Prefers using Config.Delay over Config.DelayFunc when both are set in the Config.

When implementing a retryable closure, the expected user is expected to know in advance:
	1. To retry, return the `retry.ErrRetry` error.
	2. To break the retry loop successfully, return `nil`.
	3. To break the retry loop with an Error, return a non-nil error that will be propagated to your code.
*/
func Retry(run Func, configs ...Config) error {
	var config Config
	var runErr error

	if len(configs) > 0 {
		config = configs[0]
	} else {
		config = retryDefaultConfig
	}

	if config.Delay <= 0 && config.DelayFunc == nil {
		return fmt.Errorf("Config must include EITHER a non zero Delay, or a non nil DelayFunc")
	}

	tries := 0

	for tries < config.Tries {
		runErr = run()

		if runErr != nil && errors.Is(runErr, ErrRetry) {
			tries++

			if config.Delay != 0 {
				time.Sleep(time.Second * time.Duration(config.Delay))
				continue
			}

			if config.DelayFunc != nil {
				time.Sleep(config.DelayFunc(tries))
				continue
			}
		}

		if runErr != nil {
			return runErr
		}

		return nil
	}

	return &ErrMaxRetries{
		RunErr: runErr,
	}
}
