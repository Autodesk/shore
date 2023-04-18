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
func (d *Dependencies) Load(profileName, execConfigName string) error {
	d.ShoreConfigOpts = config.ShoreConfigOpts{
		ProfileName:        profileName,
		ExecutorConfigName: execConfigName,
	}

	var err error
	d.ShoreConfig, err = config.LoadShoreConfig(d.Project)
	if err != nil {
		return err
	}

	d.Renderer, err = d.initRenderer(d.ShoreConfig)
	if err != nil {
		return err
	}

	d.Backend, err = d.initBackend(d.ShoreConfig)
	return err
}

// initRenderer initializes the Renderer based on the shore config
func (d *Dependencies) initRenderer(shoreConfig config.ShoreConfig) (renderer.Renderer, error) {
	switch strings.ToLower(shoreConfig.Renderer[`type`].(string)) {
	case JSONNET:
		d.Logger.Debug("Using the Jsonnet Renderer")
		return jsonnet.NewRenderer(d.Project.FS, d.Project.Log), nil
	default:
		return nil, fmt.Errorf("the following Renderer is undefined: %s", shoreConfig.Renderer[`type`].(string))
	}
}

// initRenderer initializes the Backend based on the shore config
func (d *Dependencies) initBackend(shoreConfig config.ShoreConfig) (backend.Backend, error) {
	execConfig := shoreConfig.GetExecutorConfig(d.ShoreConfigOpts.ExecutorConfigName)

	switch strings.ToLower(shoreConfig.Executor[`type`].(string)) {
	case SPINNAKER:
		d.Logger.Debug("Using the Spinnaker Backend")
		return spinnaker.NewClient(execConfig, d.Project.Log), nil
	default:
		return nil, fmt.Errorf("the following Executor is undefined: %s", shoreConfig.Executor[`type`].(string))
	}
}
