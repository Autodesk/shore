package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"sync"

	jsonnettest "github.com/Autodesk/shore/jsonnet/tools/jsonnet-test/pkg/jsonnet-test"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

var (
	// Version - Shore CLI version
	// This variable is populated during compile time with a variable substitution.
	// The variable should be a `const`, but `ldflags` can only operate on `var+string`.
	version = "local"

	rootCmd = &cobra.Command{
		Use:          "jt",
		Short:        "A Unit Testing Tool for Jsonnet",
		Long:         `A Unit Testing Tool for Jsonnet`,
		Version:      version,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var path string
			searchPath, _ := cmd.Flags().GetString("path")
			libs, err := cmd.Flags().GetStringSlice("libs")

			if err != nil {
				return err
			}

			if filepath.IsAbs(searchPath) {
				path = searchPath
			} else {
				workingDir, _ := os.Getwd()
				path = filepath.Join(workingDir, searchPath)
			}

			libPaths := make([]string, 0)
			for _, lib := range libs {
				if filepath.IsAbs(lib) {
					libPaths = append(libPaths, lib)
				} else {
					workingDir, _ := os.Getwd()
					libPaths = append(libPaths, filepath.Join(workingDir, lib))
				}
			}

			testFiles, err := getTestFiles(path, libPaths)

			if err != nil {
				return err
			}

			if len(testFiles) == 0 {
				return errors.New("no test file found")
			}

			if err := unitTest(testFiles, libPaths); err != nil {
				return err
			}

			plural := "s"
			if len(testFiles) == 1 {
				plural = ""
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Passed: %d test%s passed\n", len(testFiles), plural)
			return nil
		},
	}
)

func getTestFiles(path string, libPaths []string) ([]string, error) {
	var testFiles []string

	fileInfo, err := os.Lstat(path)

	if err != nil {
		return testFiles, err
	}

	if fileInfo.IsDir() == false {
		testFiles = append(testFiles, path)
		return testFiles, nil
	}

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		// Don't look grep files in tests libs.
		for _, libPath := range libPaths {
			if strings.HasPrefix(path, libPath) {
				return nil
			}
		}

		matched, err := regexp.Match(`^.*_test\.(jsonnet|libsonnet)$`, []byte(info.Name()))
		if err != nil {
			return err
		}

		if !matched {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		testFiles = append(testFiles, path)
		return nil
	})

	return testFiles, err
}

func unitTest(testFiles, libPaths []string) error {
	var testErrors error
	var wait sync.WaitGroup

	for _, testFile := range testFiles {
		wait.Add(1)

		go func(testFile string) {
			defer wait.Done()

			err := jsonnettest.JsonnetTest(testFile, libPaths...)
			if err != nil {
				testErrors = multierror.Append(testErrors, err)
			}
		}(testFile)
	}

	wait.Wait()

	return testErrors
}

func init() {
	rootCmd.Flags().StringP("path", "p", "tests", `Path to search for tests.
If a directory is provided, all files that match '*_test.(jsonnet|libsonnet)' will be included.
If a file is passed, will only run the file.`)

	rootCmd.Flags().StringSliceP("libs", "l", []string{"vendor"}, "Comma separated list of library paths")
	// Make the version easily parsable when invoking `jt --version`
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
