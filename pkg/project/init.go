package project

import (
	"bytes"
	"errors"
	"html/template"
	"regexp"
	"strconv"
	"strings"

	"github.com/Autodesk/shore/internal/gocmd"
	"github.com/sirupsen/logrus"
)

var requireGoTpl = `// +build require
// Package {{ .ShortName }} contains required external dependencies for JSONNET code.
package {{ .ShortName }}

{{if .Libraries}}
import ({{range .Libraries}}
    _ "{{.}}"{{end}}
)
{{end}}
`

var jsonnetTpl = `
function(params={}) (
	{
		"Hello": "World!"
	}
)
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
	Steps:
		1. Creates a go project (`go mod init`)
		2. Downloads required packages (`go get .. .. ..`)
		3. Creates a vendor directory (`go mod vendor`)
		4. Creates a main.pipeline.jsonnet file.
*/
func (pInit *ProjectInitialize) Init(shoreInit ShoreProjectInit) error {
	stdout, err := pInit.GoCmd.Init(shoreInit.ProjectName)

	if err != nil {
		return err
	}

	pInit.Log.Info(stdout)

	if err := pInit.createRequireGoFile(shoreInit); err != nil {
		return err
	}

	if len(shoreInit.Libraries) > 0 {
		pInit.Log.Debug("Calling go with ", shoreInit.Libraries)
		stdout, err := pInit.GoCmd.Get(shoreInit.Libraries)

		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		pInit.Log.Info(string(stdout))

		stdout, err = pInit.GoCmd.Vendor()

		if err != nil {
			return err
		}

		pInit.Log.Info(string(stdout))
	}

	pInit.Project.WriteFile("main.pipeline.jsonnet", jsonnetTpl)

	return nil
}

func (pInit *ProjectInitialize) createRequireGoFile(shoreInit ShoreProjectInit) error {
	t, err := template.New("require.go").Parse(requireGoTpl)

	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, shoreInit); err != nil {
		return err
	}

	pInit.Project.WriteFile("require.go", tpl.String())

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
