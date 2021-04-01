package project

import (
	"bytes"
	"errors"
	"html/template"
	"regexp"
	"strconv"
	"strings"

	"github.com/Autodeskshore/internal/gocmd"
	v1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	v1Dependencies "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"
	"github.com/sirupsen/logrus"
)

var testTpl = `local deployment = import '../main.pipeline.jsonnet';

local tests = [
  deployment({application: 'application', pipeline: 'pipeline', example_value: 'example_value'}),
];

local assertions = [
  {application: 'application', pipeline: 'pipeline', message: "Hello example_value!"}
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
`

var jsonnetTpl = `/**
    Creates a pipeline.
**/

function(params={}) (
	{
		"application": params["application"],
		"pipeline": params["pipeline"],
		"message": "Hello %s!" % [params["example_value"]],
	}
)
`

var readMeTpl = `# {{ .ProjectName }}
A {{ .Renderer }} project for {{ .Backend }}, initialized by Shore.
`

var gitIgnoreTpl = `# Go specific
vendor/*
!vendor/modules.txt

# IDE Specific
.vscode/
.idea
`

var e2eTpl = `application: "{{ .ProjectName }}"
pipeline: "{{ .ShortName }}-pipeline"
tests:
  "Test Success":
    execution_args:
      parameters:
        my_pipeline_param: "Test Parameter"
    assertions: { }
`

var execTpl = `application: "{{ .ProjectName }}"
pipeline: "{{ .ShortName }}-pipeline"
parameters:
  my_pipeline_param: "Example Value"
`

var renderTpl = `application: "{{ .ProjectName }}"
pipeline: "{{ .ShortName }}-pipeline"
example_value: "World"
`

// ShoreProjectInit - Common data structure to initialize a shore project.
type ShoreProjectInit struct {
	ProjectName string
	Renderer    string
	Backend     string
	Libraries   []string
}

// ShortName Creates a golang matching package name.
// github.com/Autodesk/test-init becomes testinit
func (s ShoreProjectInit) ShortName() string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	split := strings.Split(s.ProjectName, "/")
	projectName := reg.ReplaceAllString(split[len(split)-1], "")

	return projectName
}

type ProjectInitialize struct {
	Log     logrus.FieldLogger
	GoCmd   gocmd.GoWrapper
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
	if err := pInit.createFileFromTemplate("README.md", readMeTpl, shoreInit); err != nil {
		return err
	}
	if err := pInit.createFileFromTemplate("E2E.yml", e2eTpl, shoreInit); err != nil {
		return err
	}
	if err := pInit.createFileFromTemplate("render.yml", renderTpl, shoreInit); err != nil {
		return err
	}
	if err := pInit.createFileFromTemplate("exec.yml", execTpl, shoreInit); err != nil {
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
	pInit.Project.WriteFile("main.pipeline.jsonnet", jsonnetTpl)
	pInit.Project.WriteFile(".gitignore", gitIgnoreTpl)
	pInit.Project.WriteFile("tests/example_test.libsonnet", testTpl)

	return nil
}

func (pInit *ProjectInitialize) createFileFromTemplate(fileName string, fileTpl string, shoreInit ShoreProjectInit) error {
	t, err := template.New(fileName).Parse(fileTpl)

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
