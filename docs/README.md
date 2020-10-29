# Shore DOCS

Documentation on the `Shore`'s tooling and libraries.

## Project

A `Shore` project will always follow this directory structure:

```bash
./
- go.sum
- go.mod
- main.pipeline.jsonnet
- require.go
- vendor/
```

> Ignore the `go` related files for now.

When running `shore render` or `shore save`, the `main.pipeline.jsonnet` file is invoked.

## Package Management

Using `golang`'s built in `go mod` package manager, we can package our JSONNET code to share as either libraries or common pipelines.

To require a pipeline:

```bash
go get github.com/Autodesk/sponnet
```

and add a new line to the `require.go` file:

```go
// +build require
// Package mynewpipeline contains required external dependencies for JSONNET code.
package mynewpipeline

import (
 _ "github.com/Autodesk/sponnet"
)
```

then run:

```bash
go mod download && go mod vendor
```

> TODO: Automate these steps with a single command in the future.
