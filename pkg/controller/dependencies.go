package controller

import (
	"github.com/Autodesk/shore/pkg/backend"
	"github.com/Autodesk/shore/pkg/project"
	"github.com/Autodesk/shore/pkg/renderer"
	"github.com/sirupsen/logrus"
)

// Dependencies - Shared dependencies all controller MAY require
type Dependencies struct {
	Renderer renderer.Renderer
	Backend  backend.Backend
	Logger   logrus.FieldLogger
	Project  *project.Project
}
