package controller

import (
	"github.com/Autodesk/shore/pkg/backend"
	"github.com/Autodesk/shore/pkg/project"
	"github.com/Autodesk/shore/pkg/renderer"
	log "github.com/sirupsen/logrus"
)

// Dependencies - Shared dependencies all controller MAY require
type Dependencies struct {
	Renderer renderer.Renderer
	Project  project.ProjectOrganizer
	Backend  backend.Backend
	Logger   *log.Logger
}
