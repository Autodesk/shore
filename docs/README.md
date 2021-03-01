# Shore DOCS

Documentation on the `Shore`'s tooling and libraries.

# Concepts

## Package

A package in shore is essentially shared code (Typically Jsonnet). This code can be used as is, or extended in order to best build out your pipeline.

Since shore piggybacks on GoLang for package management, you can think of shore packages as go libraries which download shore files. [Read more Here](https://golang.org/ref/mod)

Each package has its own git repository and can be developed independently of anything else. Packages should use semantic versioning and tagging in order to protect users of the package and for best practices sake.

**Note**: Packages have a tight coupling to their supported backends and renderer. Meaning each package developed will have to be recoded/certified for additional renderers and backends.

### Using Packages

Using `golang`'s built in `go mod` package manager, we can package our JSONNET code to share as either libraries or common pipelines.

To require a pipeline:

```bash
go get github.com/Autodesksponnet
```

and add a new line to the `require.go` file:

```go
// +build require
// Package mynewpipeline contains required external dependencies for JSONNET code.
package mynewpipeline

import (
 _ "github.com/Autodeskspin-fe-deployer"
)
```

then run:

```bash
go mod download && go mod vendor
```

### Structure of a package

A `Shore` project will always follow this directory structure:

```bash
./
- go.sum
- go.mod
- main.pipeline.jsonnet
- require.go
- vendor/
```

> Ignore the `go` related files for now. These are used for package mangement described next.

When running `shore render` or `shore save`, the **main.pipeline.jsonnet** file is invoked.

> TODO: Automate these steps with a single command in the future.

## Renderer

Abstraction representing the templating language that will be used. We will be using `Jsonnet` by default. Each additional language will need to implement the interface:

[Renderer](https://github.com/Autodeskshore/tree/master/pkg/renderer)

## Backend

Abstraction representing the pipeline engine. By default, we will be supporting Spinnaker.

[Backend](https://github.com/Autodeskshore/tree/master/pkg/backend)
