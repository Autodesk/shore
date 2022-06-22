package shore_testing

// TestsConfig - describes the tests to run agains a pipeline
type TestsConfig struct {
	// RenderArgs
	Application string                `json:"application"`
	Pipeline    string                `json:"pipeline"`
	Timeout     int                   `json:"timeout"`
	Parallel    bool                  `json:"parallel"`
	Tests       map[string]TestConfig `json:"tests"`
	Ordering    []string              `json:"ordering"`
}

// TestConfig - describes a high level test config for a pipeline
type TestConfig struct {
	ExecArgs   map[string]interface{} `json:"execution_args"`
	Assertions map[string]Assertion   `json:"assertions"`
}

// Assertion - describes supported stage assertions
type Assertion struct {
	ExpectedStatus string                 `json:"expected_status"`
	ExpectedOutput map[string]interface{} `json:"expected_output"`
}
