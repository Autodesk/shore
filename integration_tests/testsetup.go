package integration_tests

import (
	"context"
	"os"
	"testing"

	"github.com/Autodesk/shore/pkg/backend/spinnaker"
	"github.com/Autodesk/shore/pkg/command"
	"github.com/Autodesk/shore/pkg/project"
	"github.com/Autodesk/shore/pkg/renderer/jsonnet"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
)

var testPath string = "/test"

func SetupTest(t *testing.T, f func(*testing.T, *command.Dependencies)) {
	os.Setenv("LOCAL", "true")
	os.Setenv("SHORE_PROJECT_PATH", testPath)

	memFs := afero.NewMemMapFs()
	memFs.Mkdir(testPath, os.ModePerm)

	logger, _ := test.NewNullLogger()
	s := spinnaker.NewClient(logger)
	s.CustomSpinCLI = &spinnaker.MockCustomSpinCli{}
	s.SpinCLI = &spinnaker.SpinCLI{
		ApplicationControllerAPI: &spinnaker.MockApplicationControllerAPI{},
		PipelineControllerAPI:    &spinnaker.MockPipelineControllerAPI{},
		Context:                  context.Background(),
	}

	deps := &command.Dependencies{
		Project:  project.NewShoreProject(memFs, logger),
		Renderer: jsonnet.NewRenderer(memFs, logger),
		Backend:  s,
		Logger:   logger,
	}

	f(t, deps)
}
