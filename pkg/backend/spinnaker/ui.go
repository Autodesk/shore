package spinnaker

// A fake UI implementation to mimick https://github.com/mitchellh/cli which SpinCli depends on

import "fmt"

// UI A fake UI
type UI struct{}

// Ask asks the user for input using the given query. The response is
// returned as the given string, or an error.
func (ui *UI) Ask(string) (string, error) {
	return "", fmt.Errorf("Not implemented")
}

// AskSecret asks the user for input using the given query, but does not echo
// the keystrokes to the terminal.
func (ui *UI) AskSecret(string) (string, error) {
	return "", fmt.Errorf("Not implemented")
}

// Output is called for normal standard output.
func (ui *UI) Output(string) {}

// Info is called for information related to the previous output.
// In general this may be the exact same as Output, but this gives
// Ui implementors some flexibility with output formats.
func (ui *UI) Info(string) {}

// Error is used for any error messages that might appear on standard
// error.
func (ui *UI) Error(string) {}

// Warn is used for any warning messages that might appear on standard
// error.
func (ui *UI) Warn(string) {}

// Success is used for success cases (probably?)
func (ui *UI) Success(message string) {}

// JsonOutput is used to output data in JSON format (probably?)
// Seems like the spin cli maintainers didn't follow standards, the method name should have been `JSONOutput` according to Golang best practices
func (ui *UI) JsonOutput(data interface{}) {}
