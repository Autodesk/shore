package command

import (
	"github.com/Autodeskshore/pkg/backend"
	"github.com/Autodeskshore/pkg/project"
	"github.com/Autodeskshore/pkg/renderer"
	"github.com/sirupsen/logrus"
)

// Dependencies - Shared dependencies all controller MAY require
type Dependencies struct {
	Renderer renderer.Renderer
	Backend  backend.Backend
	Logger   logrus.FieldLogger
	Project  *project.Project
}
