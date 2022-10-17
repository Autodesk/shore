package jsonnettest

import (
	"fmt"

	"encoding/json"

	"github.com/google/go-jsonnet"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var differ *diffmatchpatch.DiffMatchPatch = diffmatchpatch.New()

// TestCase - An internal representation of the Jsonnet Test Format
type TestCase struct {
	Pass       bool        `json:"bool"`
	Tests      interface{} `json:"tests"`
	Assertions interface{} `json:"assertions"`
}

// ErrTestFailure - An internal representation of the Jsonnet Test Failure
type ErrTestFailure struct {
	FileName string
	Diffs    string
}

func (e ErrTestFailure) Error() string {
	return fmt.Sprintf("File: %s:\nDiff:\n%s", e.FileName, ColorizeDiffOutput(e.Diffs))
}

// JsonnetTest - Given a file path, will run the Jsonnet Test file and diff the output.
func JsonnetTest(filepath string, libPaths ...string) error {
	// Create a Jsonnet VM locally to avoid possible race conditions when the function is called in a GoRoutine.
	var jsonnetVM *jsonnet.VM = jsonnet.MakeVM()

	if len(libPaths) > 0 {
		jsonnetVM.Importer(&jsonnet.FileImporter{JPaths: libPaths})
	}

	res, err := jsonnetVM.EvaluateFile(filepath)

	if err != nil {
		return err
	}

	var testJSON TestCase

	if err := json.Unmarshal([]byte(res), &testJSON); err != nil {
		return err
	}

	if testJSON.Pass {
		return nil
	}

	tests, _ := json.MarshalIndent(testJSON.Tests, "", "    ")
	assertions, _ := json.MarshalIndent(testJSON.Assertions, "", "    ")

	// Add new line to make the diff algorithm output happy
	testsStr := string(tests) + "\n"
	// Add new line to make the diff algorithm output happy
	assertionsStr := string(assertions) + "\n"

	// Assertions are the base we are testing against, that is why we have the tests diffed against them
	edits := myers.ComputeEdits(span.URIFromPath("assetions"), assertionsStr, testsStr)

	if len(edits) > 0 {
		diff := fmt.Sprint(gotextdiff.ToUnified("assertions", "tests", assertionsStr, edits))

		return ErrTestFailure{
			FileName: filepath,
			Diffs:    diff,
		}
	}

	return nil
}
