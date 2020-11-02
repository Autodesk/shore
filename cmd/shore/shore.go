package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Autodesk/shore/pkg/controller"
	"github.com/Autodesk/shore/pkg/fs"
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
		// TODO: `getShoreProjectDir()` should be extracted to a shared config and passed to each command
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
		// TODO: `getShoreProjectDir()` should be extracted to a shared config and passed to each command
		projectPath, err := getShoreProjectDir()

		if err != nil {
			log.Fatal(err)
		}

		pipeline, err := controller.Render(projectPath)

		if err != nil {
			log.Fatal(err)
		}

		res, err := controller.SavePipeline(pipeline)

		if err != nil {
			log.Println(err)
		}

		log.Println(res)
	},
}

func init() {
	// TODO: Add global validations to init.
	// cobra.OnInitialize()
	viper.AutomaticEnv()
	fs.InitFs(fs.OS)
	rootCmd.AddCommand(render)
	rootCmd.AddCommand(savePipeline)

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
