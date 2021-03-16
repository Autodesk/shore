package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewExecCommand - Using a Project, Renderer & Backend, executes a pipeline pipeline.
func NewExecCommand(d *Dependencies) *cobra.Command {
	var withWait bool
	var withSilent bool
	var waitTimeout int

	cmd := &cobra.Command{
		Use:   "exec",
		Short: "executes the pipeline",
		Long:  "Executes the selected pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigName("exec")

			d.Logger.Info("Exec command invoked")

			var payloadErr error

			if viper.IsSet("payload") {
				payloadErr = viper.ReadConfig(strings.NewReader(viper.GetString("payload")))
			} else {
				payloadErr = viper.ReadInConfig()
			}

			if payloadErr != nil {
				if _, ok := payloadErr.(viper.ConfigFileNotFoundError); ok {
					d.Logger.Warn(payloadErr)
				} else {
					d.Logger.Error("Failed to load the payload.")
					return payloadErr
				}
			}

			payload := viper.AllSettings()
			payloadBytes, errSerialize := jsoniter.Marshal(payload)

			if errSerialize != nil {
				d.Logger.Error("Failed serialize the payload, returned an error ", errSerialize)
				return errSerialize
			}

			// A bit of a hack, rather change this to an object later on.
			execArgs := string(payloadBytes)

			d.Logger.Debug("Calling `Backend.ExecutePipeline`")
			refID, res, err := d.Backend.ExecutePipeline(execArgs)

			if err != nil {
				return err
			}

			if !withWait {
				// Skip if the output should be silent.
				if !withSilent {
					fmt.Println(res)
				}

				return nil
			}

			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Writer = color.Error
			s.Suffix = fmt.Sprintf(" Waiting for pipeline to finish executing (%d Seconds)", waitTimeout)
			s.Start() // Start the spinner
			execDetails, res, err := d.Backend.WaitForPipelineToFinish(refID, waitTimeout)
			s.Stop() // Stop the spinner

			if err != nil {
				return err
			}
			// Return early if the output should be silent.
			if withSilent {
				return nil
			}

			fmt.Println(execDetails)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&withWait, "wait", "w", false, "Wait for the pipeline to finish execution")
	cmd.Flags().BoolVarP(&withSilent, "silent", "s", false, "Do not print JSON response to STDOUT")
	cmd.Flags().IntVarP(&waitTimeout, "timeout", "t", 60, "how long to wait (Seconds) for the pipeline to finish in Seconds. Yes Seconds.")
	cmd.Flags().StringP("payload", "p", "", "A JSON payload string. If not provided the exec.[json/yml/yaml] file is used.")
	viper.BindPFlag("payload", cmd.Flags().Lookup("payload"))

	return cmd
}
