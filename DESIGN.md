# High level technical design doc

## Reading/Rendering files

JSONNET/{INSERT LANGUAGE} files will read from `./{project_path}/pipelines/`.

Only top level files that generate a `Pipeline` object will be rendered.

```bash
{project_path}/pipelines/main.jsonnet # RENDERS A PIPELINE
{project_path}/pipelines/another-main.jsonnet  # RENDERS ANOTHER PIPELINE

# ---- DOES NOT RENDER A PIPELINE ----
{project_path}/pipelines/MYSUBDIRECTORY/not-really-main.jsonnet
```

`Pipeline` objects can be identified using a validatio method that conforms to one of the supported backends.

## Saving to a backend

The rendered output is stored in Memory and is passed on to the correct backend service provider.

As of today (20 Oct 2020) only Spinnaker is supported as a backend service.

The framework will not validate the input before pushing to a backend as combinations may be very tricky to validate.

Instead the framework will try to provide known good values for a specific backend configuration (I.E. Spinnaker)

## Tools

The framework will provide a few packages and functions OOTB for customer's to consume.

These packages will be made available through the common resources and identified at runtime.