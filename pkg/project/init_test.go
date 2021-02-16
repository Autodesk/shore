package project_test

import (
	"errors"
	"testing"

	"github.com/Autodeskshore/pkg/project"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestIsValidGoVersion(t *testing.T) {
	// Given
	goVersionStr := "go version go1.16.0 darwin/amd64"

	// Test
	res, err := project.IsValidGoVersion(goVersionStr)

	// Assert
	assert.Equal(t, true, res)
	assert.Nil(t, err)
}

func TestIsValidGoVersionNotValid(t *testing.T) {
	// Given
	goVersionStr := "go version go1.10.10 darwin/amd64"

	// Test
	res, err := project.IsValidGoVersion(goVersionStr)

	// Assert
	assert.Equal(t, false, res)
	assert.Error(t, err)
}

func TestIsValidGoVersionBrokenString(t *testing.T) {
	// Given
	goVersionStr := "goversiongo1.0.10amd64"

	// Test
	res, err := project.IsValidGoVersion(goVersionStr)

	// Assert
	assert.Equal(t, false, res)
	assert.Error(t, err)
}

func TestIsValidGoVersionNoVersion(t *testing.T) {
	// Given
	goVersionStr := "goversiongoamd64"

	// Test
	res, err := project.IsValidGoVersion(goVersionStr)

	// Assert
	assert.Equal(t, false, res)
	assert.Error(t, err)
}

func TestIsValidGoVersionReallyBrokenVersion(t *testing.T) {
	// Given
	goVersionStr := "1.11 1.12 1.13 1.1.1.1.1 10.10.10.10 13.13.13.13"

	// Test
	res, err := project.IsValidGoVersion(goVersionStr)

	// Assert
	assert.Equal(t, false, res)
	assert.Error(t, err)
}

type MockGoCmd struct {
	Dir string
	Env []string
}

func (m *MockGoCmd) Init(name string) (string, error) {
	return "", nil
}

func (m *MockGoCmd) Get(packages []string) (string, error) {
	return "", nil
}

func (m *MockGoCmd) Vendor() (string, error) {
	return "", nil
}

func (m *MockGoCmd) Version() (string, error) {
	return "", nil
}

type MockGoCmdSuccess struct {
	MockGoCmd
}

func TestInitSuccess(t *testing.T) {
	// Given
	localFs := afero.NewMemMapFs()

	init := project.ShoreProjectInit{
		ProjectName: "my-project",
		Renderer:    "jsonnet",
		Backend:     "spinnaker",
		Libraries:   []string{"github.com/Autodesk/sponnet", "github.com/***REMOVED***managed_infra"},
	}

	pInit := project.ProjectInitialize{
		Log:   Logger,
		GoCmd: &MockGoCmd{},
		Project: project.Project{
			FS:   localFs,
			Log:  Logger,
			Path: "/tmp/test/",
		},
	}

	// Test
	err := pInit.Init(init)

	// Assert
	requireExists, _ := afero.Exists(localFs, "/tmp/test/require.go")
	mainExists, _ := afero.Exists(localFs, "/tmp/test/main.pipeline.jsonnet")

	assert.True(t, requireExists)
	assert.True(t, mainExists)
	assert.Nil(t, err)
}

type MockGoCmdInitFailure struct {
	MockGoCmd
}

func (m *MockGoCmdInitFailure) Init(name string) (string, error) {
	return "", errors.New("Failed due to an internal logic error")
}

func TestInitVersionInitFailure(t *testing.T) {
	// Given
	localFs := afero.NewMemMapFs()

	init := project.ShoreProjectInit{
		ProjectName: "my-project",
		Renderer:    "jsonnet",
		Backend:     "spinnaker",
		Libraries:   []string{"github.com/Autodesk/sponnet", "github.com/***REMOVED***managed_infra"},
	}

	pInit := project.ProjectInitialize{
		Log:   Logger,
		GoCmd: &MockGoCmdInitFailure{},
		Project: project.Project{
			FS:   localFs,
			Log:  Logger,
			Path: "/tmp/test/",
		},
	}

	// Test
	err := pInit.Init(init)

	// Assert
	requireExists, _ := afero.Exists(localFs, "/tmp/test/require.go")
	mainExists, _ := afero.Exists(localFs, "/tmp/test/main.pipeline.jsonnet")

	assert.False(t, requireExists)
	assert.False(t, mainExists)
	assert.Error(t, err)
}

type MockGoCmdGetFailure struct {
	MockGoCmd
}

func (m *MockGoCmdGetFailure) Get(packages []string) (string, error) {
	return "", errors.New("Library name is wrong")
}

func TestInitVersionGetFailure(t *testing.T) {
	// Given
	localFs := afero.NewMemMapFs()

	init := project.ShoreProjectInit{
		ProjectName: "my-project",
		Renderer:    "jsonnet",
		Backend:     "spinnaker",
		Libraries:   []string{"github.com/Autodesk/sponnet", "github.com/***REMOVED***managed_infra"},
	}

	pInit := project.ProjectInitialize{
		Log:   Logger,
		GoCmd: &MockGoCmdGetFailure{},
		Project: project.Project{
			FS:   localFs,
			Log:  Logger,
			Path: "/tmp/test/",
		},
	}

	// Test
	err := pInit.Init(init)

	// Assert
	requireExists, _ := afero.Exists(localFs, "/tmp/test/require.go")
	mainExists, _ := afero.Exists(localFs, "/tmp/test/main.pipeline.jsonnet")

	assert.True(t, requireExists)
	assert.False(t, mainExists)
	assert.Error(t, err)
}

type MockGoCmdVendorFailure struct {
	MockGoCmd
}

func (m *MockGoCmdVendorFailure) Vendor() (string, error) {
	return "", errors.New("vendor failed due to an internal logic problem")
}

func TestInitVersionVendorFailure(t *testing.T) {
	// Given
	localFs := afero.NewMemMapFs()

	init := project.ShoreProjectInit{
		ProjectName: "my-project",
		Renderer:    "jsonnet",
		Backend:     "spinnaker",
		Libraries:   []string{"github.com/Autodesk/sponnet", "github.com/***REMOVED***managed_infra"},
	}

	pInit := project.ProjectInitialize{
		Log:   Logger,
		GoCmd: &MockGoCmdVendorFailure{},
		Project: project.Project{
			FS:   localFs,
			Log:  Logger,
			Path: "/tmp/test/",
		},
	}

	// Test
	err := pInit.Init(init)

	// Assert
	requireExists, _ := afero.Exists(localFs, "/tmp/test/require.go")
	mainExists, _ := afero.Exists(localFs, "/tmp/test/main.pipeline.jsonnet")

	assert.True(t, requireExists)
	assert.False(t, mainExists)
	assert.Error(t, err)
}

func TestShortName(t *testing.T) {
	// Given
	init := project.ShoreProjectInit{
		ProjectName: "my-project",
		Renderer:    "jsonnet",
		Backend:     "spinnaker",
		Libraries:   []string{"github.com/Autodesk/sponnet", "github.com/***REMOVED***managed_infra"},
	}
	// Tests
	shortName := init.ShortName()
	// Assert
	assert.Equal(t, "myproject", shortName)
}
