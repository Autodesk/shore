# Jsonnet Renderer

## Preface

This renderer is both a Reference Implementation & Minimal Viable Product.

The `Jsonnet Renderer` is designed as a "plugin" which doesn't need to be shipped with `shore core`, but side-loaded.

## Context

The `Jsonnet Renderer` ("plugin") integrates the [Jsonnet](jsonnet.org) programming language with the `Shore` Framework.

`Jsonnet` is a templating language that process's expressions and outputs `JSON` formatted text. (see more at [https://jsonnet.org](https://jsonnet.org))

Following the `shore-core` interfaces, this implements the [`Renderer` Interface](../shore-core.md#renderer---interface)

## Implementation specific details

### Embedded `Go-Jsonnet`

The `go-jsonnet` engine is embedded into the plugin and shipped with `shore-core`.

#### Embedded `Go-Jsonnet` - Pros

1. It's already implemented in `Golang` - easy to embed.
2. Easy to debug, test, validate - the runtime engine is embedded.
3. Simple to setup a development environment - `go mod` handles the setup
4. Customers don't need to install multiple tools.
5. Many of the `shared libraries` in common use are already implemented in `Jsonnet` - this can help drive initial adoption.

#### Embedded `Go-Jsonnet` - Cons

1. Only selected features are exposed via `Shore` - Rendering & Libraries.
2. Version is pinned to the `Renderer`.

### Other Solutions Considered

1. Python Renderer.

`Python Renderer` Pros:

1. *Believed to be* common enough on most developer's machines (I.E. no setup cost)
2. Python has built-in testing framework (which can be extended with `pytest`)
3. Python has many 3rd party libraries that deal with many edge cases.

`Python Renderer` Cons:

1. Not a "configuration language" - may lead to complex solutions based on Python's ability to do "everything".
2. Will require re-implementing many libraries.
3. Large Foot Print - The Python runtime is large, while for pipeline development only a fraction of it's capabilities will be used (I.E. `json.dumps()`)

### Side Notes


Although renderers are pluggable by design, the initial support will focus on `Jsonnet` due to it's implementation simplicity. In the future we can look to extend to further renderers as needed.

NOTE: If you'd like to "bring your own" pipeline and leverage Shore's testing, saving, and executing capabilities, please follow the [bring your own pipeline tutorial](https://github.com/Autodesk/shore-tutorials/tree/master/tutorials/bring-your-own-pipeline).
