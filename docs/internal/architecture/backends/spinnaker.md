# Spinnaker Backend

## Preface

This backend is both a Reference Implementation & Minimum Viable Product.
The `Spinnaker Backend` is designed as a "plugin" which doesn't need to be shipped with `shore core`, but side-loaded.

## Context

The `Spinnaker Backend` ("plugin") integrates the [Spinnaker](https://spinnaker.io) Platform with the `Shore` Framework.

> Spinnaker is an open source, multi-cloud continuous delivery platform for releasing software changes with high velocity and confidence.

(source: [https://spinnaker.io/](https://spinnaker.io/))

Following the `shore-core` interfaces, this "plugin" implements the [`Backend` Interface](../shore-core.md#backend---interface)

### SavePipeline

Given a `JSON` string (AKA "Spinnaker Pipeline") and application name, saves the pipeline in Spinnaker.

The only custom feature this backend supports is handling `NestedPipeline` during save ([see more](#nestedpipeline-support))

### ExecutePipeline

Given a `PipelineName` & `ApplicationName` executes a pipeline.
You may optionally pass parameters to the execution.

### WaitForPipelineToFinish

Given an `ExecutionId`, waits for a pipeline to finish executing (I.E. `status != RUNNING`).

### TestPipeline

A Reference implementation that implements `E2E`, `remote-test` feature set for `shore` remote testing.

When provided with a data structure ([DTO reference `TestsConfig`](../../../../pkg/backend/spinnaker/testing.go)), will go through each test case, execute a pipeline and wait for it to finish (ttl 20 minutes).

When a pipeline finishes, each stage assertion is validate against the stages `output` & `Spinnaker Status`.

For more information please see [spinnaker/backend.go `TestPipeline()`](spinnaker/../../../../../pkg/backend/spinnaker/backend.go)

## Implementation specific details

### Embedded `spin-cli` package

The `spin-cli` engine is embedded into the plugin core

#### Embedded `spin-cli` package - Pros

1. Community bug fixes and improvements.
2. Built in Authorization (based-on spin-cli `auth`)

#### Embedded `spin-cli` package - Cons

1. Relatively large dependency.
2. Currently only using the `auth, GetPipeline, SavePipeline` module.
3. Pinned to a single `spin-cli` version.

### Other solutions considered

> Call out to `spin-cli` and have it as a Shore requirement.

It doesn't make sense to require a CLI tool to communicate with an API.

Examples:

* Terraform Providers
* K8S Operators

#### Custom CLI

Implementing pieces of the `Spinnaker` API through a custom CLI.

This is required for missing implementation or bad code-generation on `spin-cli` (which uses `swagger-code-gen`).

Specific examples [`InvokePipelineConfigUsingPOST1` Never Returns Successfully](https://github.com/spinnaker/spin/blob/master/gateapi/api_pipeline_controller.go#L818-L838).

Due to the nature of `swagger-code-gen`, these bugs aren't easy to solve by contributing, and may lead to a lot of overhead.

For these cases, we have implemented a very small subset of API endpoints

##### Custom CLI - Pros

1. Small footprint - only implementing what is required and isn't well defined in `spin-cli`.
2. Easily extensible - Using `spin-cli HTTPClient` as a base after authenticating.

##### Custom CLI - Cons

1. Isn't maintained by a community - semi-con.
2. Add development complexity - which CLI to choose from? Spin-CLI or Spinnaker-CustomCli?

The `Spinnaker Backend` ended up with both solutions for 2 main reasons:

1. Time constraints - Fixing `spin-cli` wasn't feasible as it meant fixing the `swagger.yaml` file which means fixing the main `Spinnaker` APIs.
2. Footprint - Authentication to Spinnaker is the main crux. Calling APIs is less of a burden. For both implementation the same base authentication mechanism is used.

### NestedPipeline Support

A `NestedPipeline` is a pipeline which invokes a pipeline as a Spinnaker `stage`.

`NestedPipeline` is helpful when:

1. Consolidating many **features** in 1 Application - 1 Main pipeline can orchestrate the execution of other pipelines.
2. Separating environments/regions of "features".

> Example Features: Code Deployment, Secrets Management, Infra Management, etc..
