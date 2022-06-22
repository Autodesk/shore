package integration_tests

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/Autodeskshore/pkg/command"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulRemoteTestWithConfigFileNoArgs(t *testing.T) {

	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		e2eConfig := `
		{
			"application": "cosv3-state-buckets",
			"pipeline": "cosv3-dynamodb-table",
			"tests": {
				"Test Success": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				}
			}
		}
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "E2E.json"), []byte(e2eConfig), os.ModePerm)

		// Test
		testRemoteCmd := command.NewTestRemoteCommand(deps)
		testRemoteCmd.SilenceErrors = true
		testRemoteCmd.SilenceUsage = true
		err := testRemoteCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessfulRemoteTestWithConfigFileWithArgs(t *testing.T) {

	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		e2eConfig := `
		{
			"application": "cosv3-state-buckets",
			"pipeline": "cosv3-dynamodb-table",
			"tests": {
				"Test Success 1": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				},
				"Test Success 2": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				}
			}
		}
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "E2E.json"), []byte(e2eConfig), os.ModePerm)

		// Test
		testRemoteCmd := command.NewTestRemoteCommand(deps)
		testRemoteCmd.Flags().Set("test-names", "\"Test Success 1\",\"Test Success 2\"")
		testRemoteCmd.Flags().Set("stringify", "true")
		err := testRemoteCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessfulRemoteTestsConcurrent(t *testing.T) {

	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		e2eConfig := `
		{
			"application": "cosv3-state-buckets",
			"pipeline": "cosv3-dynamodb-table",
			"tests": {
				"Test Success 1": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				},
				"Test Success 2": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				}
			}
		}
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "E2E.json"), []byte(e2eConfig), os.ModePerm)

		// Test
		testRemoteCmd := command.NewTestRemoteCommand(deps)
		testRemoteCmd.Flags().Set("test-names", "\"Test Success 1\",\"Test Success 2\"")
		testRemoteCmd.Flags().Set("concurrent", "true")
		testRemoteCmd.Flags().Set("stringify", "true")
		err := testRemoteCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessfulRemoteTestWithConfigFileWithStringifyFalse(t *testing.T) {

	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		e2eConfig := `
		{
			"application": "cosv3-state-buckets",
			"pipeline": "cosv3-dynamodb-table",
			"tests": {
				"Test Success 1": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				},
				"Test Success 2": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				}
			}
		}
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "E2E.json"), []byte(e2eConfig), os.ModePerm)

		// Test
		testRemoteCmd := command.NewTestRemoteCommand(deps)
		testRemoteCmd.Flags().Set("test-names", "\"Test Success 1\",\"Test Success 2\"")
		testRemoteCmd.Flags().Set("stringify", "false")
		err := testRemoteCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessfulRemoteTestWithConfigFileWithArgsNotAll(t *testing.T) {

	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		e2eConfig := `
		{
			"application": "cosv3-state-buckets",
			"pipeline": "cosv3-dynamodb-table",
			"tests": {
				"Test Success 1": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				},
				"Test Success 2": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				},
				"Test Success 3": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				}
			}
		}
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "E2E.json"), []byte(e2eConfig), os.ModePerm)

		// Test
		testRemoteCmd := command.NewTestRemoteCommand(deps)
		testRemoteCmd.Flags().Set("test-names", "\"Test Success 1\",\"Test Success 3\"")
		// No way to check if "Test Success 2" ran or not.
		err := testRemoteCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestFailRemoteTestIncorrectArgs(t *testing.T) {

	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		e2eConfig := `
		{
			"application": "cosv3-state-buckets",
			"pipeline": "cosv3-dynamodb-table",
			"tests": {
				"Test Success 1": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				}
			}
		}
		`
		failingTestName := "test that doesn't exist"
		execError := fmt.Sprintf("the provided E2E configuration does not contain [%s] as a test", failingTestName)

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "E2E.json"), []byte(e2eConfig), os.ModePerm)

		// Test
		testRemoteCmd := command.NewTestRemoteCommand(deps)
		testRemoteCmd.Flags().Set("test-names", fmt.Sprintf("\"%s\"", failingTestName))
		// No way to check if "Test Success 2" ran or not.
		err := testRemoteCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, execError, err.Error())
	})
}

func TestSuccessfulRemoteTestBlankArg(t *testing.T) {

	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		e2eConfig := `
		{
			"application": "cosv3-state-buckets",
			"pipeline": "cosv3-dynamodb-table",
			"tests": {
				"Test Success 1": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				}
			}
		}
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "E2E.json"), []byte(e2eConfig), os.ModePerm)

		// Test
		testRemoteCmd := command.NewTestRemoteCommand(deps)
		testRemoteCmd.Flags().Set("test-names", "")
		testRemoteCmd.Flags().Set("stringify", "")
		// No way to check if "Test Success 2" ran or not.
		err := testRemoteCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessfulRemoteTestWithConfigFileWithOrdering(t *testing.T) {

	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		e2eConfig := `
		{
			"application": "cosv3-state-buckets",
			"pipeline": "cosv3-dynamodb-table",
			"ordering": ["Test Success 1", "Test Success 2"],
			"tests": {
				"Test Success 1": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				},
				"Test Success 2": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				}
			}
		}
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "E2E.json"), []byte(e2eConfig), os.ModePerm)

		// Test
		testRemoteCmd := command.NewTestRemoteCommand(deps)
		err := testRemoteCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessfulRemoteTestWithConfigFileWithEmptyOrdering(t *testing.T) {

	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		e2eConfig := `
		{
			"application": "cosv3-state-buckets",
			"pipeline": "cosv3-dynamodb-table",
			"ordering": [],
			"tests": {
				"Test Success 1": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				},
				"Test Success 2": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				}
			}
		}
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "E2E.json"), []byte(e2eConfig), os.ModePerm)

		// Test
		testRemoteCmd := command.NewTestRemoteCommand(deps)
		err := testRemoteCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestFailRemoteTestWithConfigFileWithIncorrectOrdering(t *testing.T) {

	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		e2eConfig := `
		{
			"application": "cosv3-state-buckets",
			"pipeline": "cosv3-dynamodb-table",
			"ordering": ["yolo"],
			"tests": {
				"Test Success 1": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				},
				"Test Success 2": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				}
			}
		}
		`
		execError := "the provided E2E configuration does not contain [yolo] as a test"

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "E2E.json"), []byte(e2eConfig), os.ModePerm)

		// Test
		testRemoteCmd := command.NewTestRemoteCommand(deps)
		err := testRemoteCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, execError, err.Error())
	})
}
