# Shore

Shore (verb) is a tool used to develop deployment pipelines for different pipeline based products like Spinnaker.

## Installing shore

<!-- TODO Fill these with the correct information -->

### As a binary

#### Released

1. Shore released binaries can be found in archived format at shore's github releases page:  <https://github.com/Autodesk/shore/releases>
2. To download the binary relevant for you platform click on the file name and save the file to the disk.
3. Un-archive, e.g. for tar.gz for Mac OS:

```shell
cd ~/Downloads
tar -xzf shore_0.0.2_Darwin_x86_64.tar.gz -C /usr/local/bin/
chmod +x /usr/local/bin/shore
shore -h
```

#### Current/Nightly build

<!-- TODO Fill these with the correct information -->

### As a docker image

1. Shore released docker images links can be found at shore's github releases page:  <https://github.com/Autodesk/shore/releases>

```shell
docker pull ***REMOVED***/shore/shore:v0.0.2
docker run ***REMOVED***/shore/shore:v0.0.2 -h
```

## Building Shore

```bash
git clone git@github.com:Autodesk/shore.git

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

As of today (30 Mar 2021) only Spinnaker is supported as a backend service.

The framework will not validate the input before pushing to a backend as combinations may be very tricky to validate.

Instead the framework will try to provide known good values for a specific backend configuration (I.E. Spinnaker)

### Tools

The framework will provide a few packages and functions for customer's to consume.

These packages will be made available through the common resources and identified at runtime.

To get these common resources, we recommend using [Jsonnet-Bundler](https://github.com/jsonnet-bundler/jsonnet-bundler/)

## Release

<!-- TODO Fill these with the correct information -->
