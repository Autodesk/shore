package controller

import (
	"fmt"
	"time"

	"github.com/Autodesk/shore/internal/gocmd"
	"github.com/Autodesk/shore/pkg/project"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// NewProjectCommand - Creates the `project` subcommand
func NewProjectCommand(d *Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Collection of project related commands",
	}

	cmd.AddCommand(NewProjectInitCommand(d))

	return cmd
}

var projectNamePrompt = promptui.Prompt{
	Label: "Project Name",
	Validate: func(i string) error {
		if len(i) <= 0 {
			return fmt.Errorf("`Project Name` cannot be empty")
		}

		return nil
	},
}

// Should be setup through some automation of selected plugins
var rendererPrompt = promptui.Select{
	Label: "Frontend Renderer",
	Items: []string{"Jsonnnet"},
}

// Should be setup through some automation of selected plugins
var backendPrompt = promptui.Select{
	Label: "Pipeline Backend",
	Items: []string{"Spinnaker"},
}

var addLibsPrompt = promptui.Prompt{
	Label:     "Add libraries?",
	IsConfirm: true,
	Default:   "y",
}

var libraryPrompt = promptui.Prompt{
	Label: "Library path (leave empty to continue)",
}

// NewProjectInitCommand - initialize a `shore` project.
func NewProjectInitCommand(d *Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new shore project",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := d.Project.GetProjectPath()

			if err != nil {
				return err
			}

			goCmd := gocmd.NewGoCmd(path)

			if version, err := goCmd.Version(); err != nil {
				if isValid, err := project.IsValidGoVersion(string(version)); isValid != true && err != nil {
					return err
				}
			}

			shoreProjectInit, err := getShoreInitValues()

			if err != nil {
				return err
			}

			pInit := &project.ProjectInitialize{
				Log:     d.Logger,
				GoCmd:   goCmd,
				Project: *d.Project,
			}

			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Suffix = " Setting up the environment, this may take a few moments (depending on Internet traffic!)"
			s.Start() // Start the spinner

			if err = pInit.Init(shoreProjectInit); err != nil {
				return err
			}

			s.Stop()

			color.Green("Project %s has been created successfully!", shoreProjectInit.ProjectName)
			color.Cyan("Try running `shore render`")

			return nil
		},
	}

	return cmd
}

func getShoreInitValues() (project.ShoreProjectInit, error) {
	projectName, err := projectNamePrompt.Run()

	if err != nil {
		return project.ShoreProjectInit{}, err
	}

	_, renderer, err := rendererPrompt.Run()

	if err != nil {
		return project.ShoreProjectInit{}, err
	}

	_, backend, err := backendPrompt.Run()

	if err != nil {
		return project.ShoreProjectInit{}, err
	}

	// Confirmation prompt returns a string, but we don't care.
	// Err is used as a bool here for some reason.
	// TODO: check if this prompt could be customized to return `bool, err`.
	_, err = addLibsPrompt.Run()

	var libs []string
	if err == nil {
		for {
			lib, err := libraryPrompt.Run()

			if err != nil {
				// d.Logger.Error(err)
			}

			if err == promptui.ErrInterrupt {
				return project.ShoreProjectInit{}, err
			}

			if lib == "" {
				break
			}

			libs = append(libs, lib)
		}
	}

	return project.ShoreProjectInit{
		ProjectName: projectName,
		Renderer:    renderer,
		Backend:     backend,
		Libraries:   libs,
	}, nil
}
