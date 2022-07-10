package command

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/fatih/color"

	"github.com/Autodeskshore/pkg/renderer"
	"github.com/nsf/jsondiff"
	"github.com/spf13/cobra"
)

// NewDiffCommand - A Cobra wrapper for the common Diff function.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func NewDiffCommand(d *Dependencies) *cobra.Command {
	var renderValues string
	var skipMatches string

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Difference between current and desired state.",
		Long:  `Shows difference between current and desired state of the pipeline.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			settingsBytes, err := GetConfigFileOrFlag(d.Project, "render", renderValues)

			var confErr *DefaultConfErr

			if err != nil && !errors.As(errors.Unwrap(err), &confErr) {
				return err
			}

			err = Diff(d, settingsBytes, skipMatches, renderer.MainFileName)

			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&renderValues, "values", "r", "", "A JSON string for the render. If not provided the render.[json/yml/yaml] file is used.")
	cmd.Flags().StringVarP(&skipMatches, "skip", "s", "false", "If true, skip the matching parts in the command output, default is false.")

	return cmd
}

// Diff Using a Project & Renderer & Get, renders the pipeline and shows the difference between current and desired state.
func Diff(d *Dependencies, settings []byte, skipMatches string, renderType renderer.RenderType) error {
	// TODO: For future DevX, aggregate errors and return them together.
	d.Logger.Info("Diff function started")

	d.Logger.Debug("GetProjectPath")
	projectPath, err := d.Project.GetProjectPath()
	d.Logger.Debug("GetProjectPath returned ", projectPath)

	if err != nil {
		d.Logger.Error("GetProjectPath returned an error ", err)
		return err
	}

	// A bit of a hack, rather change this to an object later on.
	renderArgs := string(settings)
	d.Logger.Debug("Args returned:\n", renderArgs)
	d.Logger.Info("calling Renderer.Render with projectPath ", projectPath, " and renderArgs ", renderArgs)

	desiredPipelineString, desiredPipelineInterface := getDesiredPipeline(d, projectPath, renderArgs, renderType)

	IDToPipelineMap := make(map[string]interface{})
	fillPipelineMap(d, IDToPipelineMap, desiredPipelineInterface)

	args := make(map[string]interface{})
	err = json.UnmarshalFromString(renderArgs, &args)

	if err != nil {
		d.Logger.Error("json.Unmarshal Could not unmarshell the rendered file ", err)
		return err
	}

	currentPipelineString, _ := getCurrentPipeline(d, args, IDToPipelineMap)

	diffAndPrint(args, currentPipelineString, desiredPipelineString, skipMatches)
	return nil
}

/*fillPipelineMap Populate the IDToPipelineMap (map[id] -> current pipeline configuration)
This function is needed to swap the content of nested pipelines gathered from spinnaker API
from Ids of to their full objects or names respectively to the kind NestedPipelineStage vs PipelineStage.
*/
func fillPipelineMap(d *Dependencies, IDToPipelineMap map[string]interface{},
	parentPipeline map[string]interface{}) {

	var nestedApplicationName string
	var nestedPipelineName string
	var nestedPipelineObject map[string]interface{}

	stages, exists := parentPipeline["stages"]
	if !exists {
		return
	}

	for _, stage := range stages.([]interface{}) {
		stage := stage.(map[string]interface{})

		// if there is a nested pipeline in the local configuration
		if nestedPipeline, exists := stage["pipeline"]; exists {
			nestedApplicationName = stage["application"].(string)
			// if this pipeline is not a string it is a kind of NestedPipelineStage which is an object.
			// vs PipelineStage which is a string (the name of the pipeline)
			itIsFullObjectPipeline := reflect.TypeOf(nestedPipeline).String() != "string"

			if itIsFullObjectPipeline {
				nestedPipelineObject = nestedPipeline.(map[string]interface{})
				nestedPipelineName = nestedPipelineObject["name"].(string)
			} else {
				nestedPipelineName = nestedPipeline.(string)
			}

			currentNestedPipeline, _, _ := d.Backend.GetPipeline(nestedApplicationName, nestedPipelineName)
			if currentNestedPipeline != nil {
				id := currentNestedPipeline["id"].(string)
				IDToPipelineMap[id] = currentNestedPipeline
				// need further recursion only if is type of NestedPipelineStage
				if itIsFullObjectPipeline {
					IDToPipelineMap[id] = currentNestedPipeline
					fillPipelineMap(d, IDToPipelineMap, nestedPipelineObject)
				} else {
					IDToPipelineMap[id] = currentNestedPipeline["name"].(string)
				}
			}
		}
	}
}

// formatCurrentPipeline Format the current pipeline object recursively
func formatCurrentPipeline(d *Dependencies, IDToPipelineMap map[string]interface{},
	parentPipeline map[string]interface{}) {
	cleanKeys(parentPipeline)

	stages, exists := (parentPipeline)["stages"]
	if !exists {
		return
	}

	for _, stage := range stages.([]interface{}) {
		stage := stage.(map[string]interface{})
		if nestedPipelineID, exists := stage["pipeline"]; exists {
			stage["pipeline"] = (IDToPipelineMap)[nestedPipelineID.(string)]
			nestedPipeline := stage["pipeline"]
			if nestedPipeline != nil && reflect.TypeOf(nestedPipeline).String() != "string" {
				formatCurrentPipeline(d, IDToPipelineMap, nestedPipeline.(map[string]interface{}))
			}
		}
	}
}

// cleanKeys Clean some spinnaker generated added fields
func cleanKeys(pipeline map[string]interface{}) {
	keysToDelete := [5]string{"id", "index", "lastModifiedBy", "updateTs", "schema"}
	for _, key := range keysToDelete {
		delete(pipeline, key)
	}
}

// getDesiredPipeline Returns the desired pipeline configuration as string and map[string]interface{}
func getDesiredPipeline(d *Dependencies, projectPath string, renderArgs string, renderType renderer.RenderType) ([]byte, map[string]interface{}) {
	desiredPipelineString, err := d.Renderer.Render(projectPath, renderArgs, renderType)

	if err != nil {
		d.Logger.Error("Renderer.Render returned an error ", err)
		return nil, nil
	}

	var desiredPipelineInterface map[string]interface{}
	err = json.UnmarshalFromString(desiredPipelineString, &desiredPipelineInterface)

	if err != nil {
		d.Logger.Error("json.UnmarshalFromString Could not unmarshell ", err)
		return []byte(desiredPipelineString), nil
	}

	return []byte(desiredPipelineString), desiredPipelineInterface
}

// getCurrentPipeline Returns the current pipeline configuration as string and map[string]interface{}
func getCurrentPipeline(d *Dependencies, args map[string]interface{},
	IDToPipelineMap map[string]interface{}) ([]byte, map[string]interface{}) {

	application := args["application"].(string)
	pipeline := args["pipeline"].(string)
	currentPipelineInterface, _, err := d.Backend.GetPipeline(application, pipeline)

	if err != nil {
		d.Logger.Error("Backend.GetPipeline returned an error ", err)
		return nil, nil
	}

	formatCurrentPipeline(d, IDToPipelineMap, currentPipelineInterface)

	currentPipelineString, err := json.Marshal(currentPipelineInterface)

	if err != nil {
		d.Logger.Error("Diff returned an error ", err)
		return nil, currentPipelineInterface
	}

	return currentPipelineString, currentPipelineInterface
}

// diffAndPrint Creating a diff string from 2 pipeline json strings and print them nicely to stdout
func diffAndPrint(args map[string]interface{}, currentPipelineString []byte, desiredPipelineString []byte, skipMatches string) {

	application := args["application"].(string)
	pipeline := args["pipeline"].(string)

	diffOptions := jsondiff.DefaultConsoleOptions()
	diffOptions.SkippedObjectProperty = jsondiff.SkippedObjectProperty

	if skipMatches == "true" {
		diffOptions.SkipMatches = true
	}

	diffType, diffStr := jsondiff.Compare(currentPipelineString, desiredPipelineString, &diffOptions)

	boldUnderline := color.New(color.Bold, color.Underline)
	bold := color.New(color.Bold)

	boldUnderline.Println("\nShore Difference Output:")
	bold.Printf("Application: %v\nPipeline: %v\n\n", application, pipeline)
	if diffType == jsondiff.FullMatch {
		bold.Printf("There Are No Changes in Configuration!\n\n")
	}

	fmt.Println(diffStr)
}
