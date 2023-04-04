package command

import (
	"fmt"

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
	Renderer    renderer.Renderer
	Backend     backend.Backend
	Logger      logrus.FieldLogger
	Project     *project.Project
	ShoreConfig config.ShoreConfig
}

// NewDependencies - Creates a Dependencies struct.
func NewDependencies(p *project.Project) (*Dependencies, error) {
	var chosenRenderer renderer.Renderer
	var chosenBackend backend.Backend

	shoreConfig, err := config.LoadShoreConfig(p)
	if err != nil {
		return &Dependencies{}, err
	}

	// Select the Renderer
	switch shoreConfig.Renderer[`type`] {
	case JSONNET:
		p.Log.Debug("Using the Jsonnet Renderer")
		chosenRenderer = jsonnet.NewRenderer(p.FS, p.Log)
	default:
		return &Dependencies{}, fmt.Errorf("the following Renderer is undefined: %s", shoreConfig.Renderer[`type`].(string))
	}

	// Select the Backend
	switch shoreConfig.Executor[`type`] {
	case SPINNAKER:
		p.Log.Debug("Using the Spinnaker Backend")
		chosenBackend = spinnaker.NewClient(p.Log)
	default:
		return &Dependencies{}, fmt.Errorf("the following Backend is undefined: %s", shoreConfig.Executor[`type`].(string))
	}

	return &Dependencies{
		Project:     p,
		Renderer:    chosenRenderer,
		Backend:     chosenBackend,
		Logger:      p.Log,
		ShoreConfig: shoreConfig,
	}, nil
}
