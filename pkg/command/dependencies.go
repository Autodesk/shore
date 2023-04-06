package command

import (
	"fmt"
	"strings"

	"github.com/Autodesk/shore/pkg/backend"
	"github.com/Autodesk/shore/pkg/backend/spinnaker"
	"github.com/Autodesk/shore/pkg/config"
	"github.com/Autodesk/shore/pkg/project"
	"github.com/Autodesk/shore/pkg/renderer"
	"github.com/Autodesk/shore/pkg/renderer/jsonnet"
	"github.com/sirupsen/logrus"
)

// Renderer Enum
const (
	// JSONNET - Jsonnet Renderer
	JSONNET string = "jsonnet"
)

// Backend Enum
const (
	// SPINNAKER - Spinnaker Backend
	SPINNAKER string = "spinnaker"
)

// Dependencies - Shared dependencies all controller MAY require
type Dependencies struct {
	Renderer        renderer.Renderer
	Backend         backend.Backend
	Logger          logrus.FieldLogger
	Project         *project.Project
	ShoreConfig     config.ShoreConfig
	ShoreConfigOpts config.ShoreConfigOpts
}

// Load - loads the shore config and sets the renderer and backend
func (d *Dependencies) Load() error {
	var chosenRenderer renderer.Renderer
	var chosenBackend backend.Backend

	shoreConfig, err := config.LoadShoreConfig(d.Project)
	if err != nil {
		return err
	}

	// Select the Renderer
	switch strings.ToLower(shoreConfig.Renderer[`type`].(string)) {
	case JSONNET:
		d.Logger.Debug("Using the Jsonnet Renderer")
		chosenRenderer = jsonnet.NewRenderer(d.Project.FS, d.Project.Log)
	default:
		return fmt.Errorf("the following Renderer is undefined: %s", shoreConfig.Renderer[`type`].(string))
	}
	d.Renderer = chosenRenderer

	// Select the Backend
	switch strings.ToLower(shoreConfig.Executor[`type`].(string)) {
	case SPINNAKER:
		d.Logger.Debug("Using the Spinnaker Backend")
		chosenBackend = spinnaker.NewClient(d.Project.Log)
	default:
		return fmt.Errorf("the following Executor is undefined: %s", shoreConfig.Executor[`type`].(string))
	}
	d.Backend = chosenBackend

	return nil
}
