# Get Pipeline ID

This example shows how to invoke and get remote Pipeline IDs defined as an `application/pipeline` properties.

The example is complex by testing multiple scenarios:

* `main-pipeline`:
  * Remote Pipeline Execution stage
  * Create a nested pipeline stage with:
    * Remote Pipeline Execution stage
    * Resolve Artifact From Pipeline
* `target-pipeline`:
  * Evaluate Artifact
* `remote-pipeline`:
  * Trigger Pipeline on success of `main-pipeline`.
  * Resolve Artifact From Pipeline

## Explaining The Example

### Main Invokes The Nested Pipeline

![main-pipeline](./assets/01-main-pipeline.png)

### Nested Pipeline Invokes `get-pipeline-id-target-pipeline`

![main-pipeline](./assets/02-nested-pipeline.png)

### `get-pipeline-id-target-pipeline` evaluates a new artifact

![main-pipeline](./assets/03-triger-evaluate.png)

### Nested Pipeline Retrieves the evaluated artifact from `get-pipeline-id-target-pipeline`

![main-pipeline](./assets/04-resolve-artifact.png)

### Main Pipeline Invokes `get-pipeline-id-target-pipeline` (show top level invocation works as well)

![main-pipeline](./assets/05-trigger-targer-top-level.png)

### `get-pipeline-id-remote-pipeline` is invoked via Spinnaker when `get-pipeline-id-main-pipeline` ends successfully

![main-pipeline](./assets/06-trigger-on-main-pipeline-success-different-app.png)

### `get-pipeline-id-remote-pipeline` configuration

![main-pipeline](./assets/07-trigger-on-pipeline-configuration.png)
