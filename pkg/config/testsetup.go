package config

import (
	"os"
	"testing"

	"github.com/Autodeskshore/pkg/project"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
)

var testPath = "/test"

// SetupTest is a general purpose set up for the config tests.
func SetupTest(t *testing.T, f func(*testing.T, *project.Project)) {
	os.Setenv("LOCAL", "true")
	os.Setenv("SHORE_PROJECT_PATH", testPath)

	memFs := afero.NewMemMapFs()
	memFs.Mkdir(testPath, os.ModePerm)

	logger, _ := test.NewNullLogger()

	proj := project.NewShoreProject(memFs, logger)

	f(t, proj)
}
