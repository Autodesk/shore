package controller

import (
	"github.com/Autodesk/shore/pkg/backend"
	"github.com/Autodesk/shore/pkg/project"
	"github.com/Autodesk/shore/pkg/renderer"
)

type Dependencies struct {
	Renderer renderer.Renderer
	Project  project.ProjectOrganizer
	Backend  backend.Backend
}
