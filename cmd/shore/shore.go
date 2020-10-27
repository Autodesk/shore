package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/Autodesk/shore/pkg/backend/spinnaker"
	"github.com/Autodesk/shore/pkg/controller"
	"github.com/spf13/cobra"
)

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
		res, err := controller.Render()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(strings.Join(res, "\n"))
	},
}

var savePipeline = &cobra.Command{
	Use:   "save",
	Short: "save the pipelines",
	Long:  "Walk through the `pipelines` directory, render & save the pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		// All business logic should be abstracted to business requirement specific functions (AKA controllers or similar)
		res, err := controller.Render()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(strings.Join(res, "\n"))

		cli, err := spinnaker.NewClient()
		for _, pipeline := range res {
			res, err := cli.SavePipeline(pipeline)
			runtime.Breakpoint()
			if err != nil {
				log.Println(err)
			}

			log.Println(res)
		}
	},
}

func init() {
	// TODO: Add global validations to init.
	// cobra.OnInitialize()
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
