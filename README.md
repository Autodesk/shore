# Shore

[![Build Status](https://***REMOVED***.***REMOVED***/buildStatus/icon?job=forge-cd-services%2Fshore%2Fmaster)](https://***REMOVED***.***REMOVED***/job/forge-cd-services/job/shore/job/master/)
[![Codacy Badge](https://code-quality.autodesk.com:16006/project/badge/Grade/6089dc45142b46b29cc22b9b8e357a75)](https://code-quality.autodesk.com:443/manual/***REMOVED***shore?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=forge-cd-services/shore&amp;utm_campaign=Badge_Grade)

Shore (verb) is a tool used to develop deployment pipelines for different pipeline based products like Spinnaker.

## Building Shore

```bash
git clone git@github.com:forge-cd-service/shore.git

export GOPRIVATE="github.com"
export GOPROXY="https://:@***REMOVED***/***REMOVED***/gocenter/"

go mod download
go mod vendor
go build -o shore cmd/shore/shore.go
./shore
```

### Reading/Rendering files

JSONNET/{INSERT LANGUAGE} files will read from `./{project_path}/main.pipeline.jsonnet`.

Only top level files that generate a `Pipeline` object will be rendered.

```bash
{project_path}/main.pipeline.jsonnet # RENDERS A PIPELINE

# ---- DOES NOT RENDER A PIPELINE ----
{project_path}/pipelines/MYSUBDIRECTORY/not-really-main.jsonnet
```

`Pipeline` objects can be identified using a validation method that conforms to one of the supported backends.

### Saving to a backend

The rendered output is stored in Memory and is passed on to the correct backend service provider.

As of today (20 Oct 2020) only Spinnaker is supported as a backend service.

The framework will not validate the input before pushing to a backend as combinations may be very tricky to validate.

Instead the framework will try to provide known good values for a specific backend configuration (I.E. Spinnaker)

### Tools

The framework will provide a few packages and functions for customer's to consume.

These packages will be made available through the common resources and identified at runtime.

To get these common resources, we recommend using [Jsonnet-Bundler](https://github.com/jsonnet-bundler/jsonnet-bundler/)

# Release

[Jenkins Job](https://master-11.***REMOVED***/job/***REMOVED***job/shore/)

For `master` branch merges the [`Jenkins`]('./Jenkinsfile') will create a new file in [Artifactory](https://***REMOVED***.dev.***REMOVED***/***REMOVED***/webapp/#/artifacts/browse/tree/General/SHORE-dist)

The format is `shore-${version}-${branch_name}-${build_number}-${architecture}`.

These builds are recommended for testing, debugging and sharing with other contributors for validations.
