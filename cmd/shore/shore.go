package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Autodesk/shore/pkg/backend/spinnaker"
	"github.com/Autodesk/shore/pkg/controller"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getShoreProjectDir() (string, error) {
	// Magic variable to allow working shore without actually beign in the path
	// For dev purposes only!!
	if isLocal := viper.GetBool("LOCAL"); isLocal == true {
		projectPath := viper.GetString("SHORE_PROJECT_PATH")

		if projectPath == "" {
			return "", fmt.Errorf("env variable `SHORE_PROJECT_PATH` is not set")
		}

		return projectPath, nil
	}

	return os.Getwd()
}

var rootCmd = &cobra.Command{
	Use:   "shore",
	Short: "Shore - Pipeline Framework",
	Long:  "A Pipeline development framework for integrated pipelines.",
}

var render = &cobra.Command{
	Use:   "render",
	Short: "render a pipeline",
	Long:  "Walk through the `pipelines` directory, renderer the pipelines and output to STDOUT",
	Run: func(cmd *cobra.Command, args []string) {
		// All business logic should be abstracted to business requirement specific functions (AKA controllers or similar)
		projectPath, err := getShoreProjectDir()

		if err != nil {
			log.Fatal(err)
		}

		res, err := controller.Render(projectPath)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res)
	},
}

var savePipeline = &cobra.Command{
	Use:   "save",
	Short: "save the pipelines",
	Long:  "Walk through the `pipelines` directory, render & save the pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		// All business logic should be abstracted to business requirement specific functions (AKA controllers or similar)
		projectPath, err := getShoreProjectDir()

		if err != nil {
			log.Fatal(err)
		}

		pipeline, err := controller.Render(projectPath)

		if err != nil {
			log.Fatal(err)
		}

		cli, err := spinnaker.NewClient()

		if err != nil {
			log.Fatal(err)
		}

		res, err := cli.SavePipeline(pipeline)

		if err != nil {
			log.Println(err)
		}

		log.Println(res)
	},
}

func init() {
	// TODO: Add global validations to init.
	// cobra.OnInitialize()
	rootCmd.AddCommand(render)
	rootCmd.AddCommand(savePipeline)
	viper.AutomaticEnv()

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
