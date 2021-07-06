package project

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	v1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	v1Dependencies "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"
	"github.com/sirupsen/logrus"
)

// Holds the content of the whole `templates` directory
//go:embed templates/*
var templates embed.FS

// ShoreProjectInit - Common data structure to initialize a shore project.
type ShoreProjectInit struct {
	projectName string
	Renderer    string
	Backend     string
	Libraries   []string
}

// NewShoreProjectInit - Creates a ShoreProjectInit
func NewShoreProjectInit(projectName, render, backend string, libs []string) ShoreProjectInit {
	return ShoreProjectInit{
		projectName: projectName,
		Renderer:    render,
		Backend:     backend,
		Libraries:   libs,
	}
}

// ProjectName Creates a golang matching package name. It is santized.
// Aligning with what Spinnaker expects:
// https://github.com/spinnaker/deck/blob/5a9768bc6db2f527a73d6b1f5fb3120c101e094b/app/scripts/modules/core/src/pipeline/create/CreatePipelineModal.tsx#L290
// Example - github.com/Autodesk/test-init becomes test-init
func (s ShoreProjectInit) ProjectName() string {
	reg, _ := regexp.Compile(`[\^/\\?%#]*`)
	split := strings.Split(s.projectName, "/")
	projectName := reg.ReplaceAllString(split[len(split)-1], "")

	return projectName
}

// AppName - Creates an Application name that is santized.
// Aligning with what Spinnaker expects, see comment here:
// https://github.com/Autodeskshore/issues/108#issuecomment-1607689
func (s ShoreProjectInit) AppName() string {
	reg, _ := regexp.Compile(`[\W_]*`)
	split := strings.Split(s.projectName, "/")
	projectName := reg.ReplaceAllString(split[len(split)-1], "")

	return projectName
}

type ProjectInitialize struct {
	Log     logrus.FieldLogger
	Project Project
}

/*
Init - Initializes a shore project
	This all or nothing method wraps all the necessary required steps to prep a shore project for a user.

	Creates the following files:
	- README.md
	- E2E.yml
	- render.yml
	- exec.yml
	- jsonnetfile.json
	- main.pipeline.jsonnet
	- .gitignore
	- tests/example_test.libsonnet

	Does not run jsonnent-bundler install (`jb install`).
*/
func (pInit *ProjectInitialize) Init(shoreInit ShoreProjectInit) error {
	if err := pInit.createFileFromTemplate("README.md", shoreInit); err != nil {
		return err
	}
	if err := pInit.createFileFromTemplate("E2E.yml", shoreInit); err != nil {
		return err
	}
	if err := pInit.createFileFromTemplate("render.yml", shoreInit); err != nil {
		return err
	}
	if err := pInit.createFileFromTemplate("exec.yml", shoreInit); err != nil {
		return err
	}
	if err := pInit.createFileFromTemplate("main.pipeline.jsonnet", shoreInit); err != nil {
		return err
	}
	if err := pInit.createFileFromTemplate(".gitignore", shoreInit); err != nil {
		return err
	}
	if err := pInit.createFileFromTemplate("tests/example_test.libsonnet", shoreInit); err != nil {
		return err
	}

	jsonnetDependencies := make(map[string]v1Dependencies.Dependency)
	for _, libURL := range shoreInit.Libraries {
		// Strip the protocal
		depName := strings.Replace(libURL, "http://", "", 1)
		depName = strings.Replace(libURL, "https://", "", 1)
		jsonnetDependencies[depName] = *v1Dependencies.Parse("", libURL)

	}
	jsonnetFileStruct := v1.JsonnetFile{
		LegacyImports: false,
		Dependencies:  jsonnetDependencies,
	}
	jsonnetFileBytes, err := jsonnetFileStruct.MarshalJSON()
	if err != nil {
		return err
	}

	pInit.Project.WriteFile("jsonnetfile.json", string(jsonnetFileBytes))

	return nil
}

func (pInit *ProjectInitialize) createFileFromTemplate(fileName string, shoreInit ShoreProjectInit) error {
	templateContent, err := templates.ReadFile(fmt.Sprintf("templates/%s", fileName))
	if err != nil {
		return err
	}

	t, err := template.New(fileName).Parse(string(templateContent))

	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, shoreInit); err != nil {
		return err
	}

	pInit.Project.WriteFile(fileName, tpl.String())

	return nil
}

// IsValidGoVersion - Checks if the currently installed response from Golang binary version command has GOMOD support
func IsValidGoVersion(golangVersion string) (bool, error) {
	reg := regexp.MustCompile(`\d+\.(?P<minor>\d+)`)
	matches := reg.FindAllStringSubmatch(string(golangVersion), -1)

	if len(matches) < 1 || len(matches) > 1 {
		return false, errors.New("the text provided doesn't seem to be from a Golang binary, please make sure the Golang binary is installed correctly and is on your path")
	}

	if len(matches) == 1 {
		minorVersion, _ := strconv.Atoi(matches[0][1])

		if minorVersion < 11 {
			return false, errors.New("Golang version must be 11 or higher (must support GOMOD)")
		}
	}

	return true, nil
}
