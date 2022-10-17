# Jsonnet-Test

A simple CLI tool to make `jsonnet` unit tests a bit easier.

## How to run

### Docker

Running with docker is the easiest way to get up and running quickly.

```bash
# TODO Add instructions for Docker
```

Installing:

## Help

```docs
A Unit Testing Tool for Jsonnet

Usage:
  jt [flags]

Flags:
  -h, --help           help for jt
  -l, --libs strings   Comma separated list of library paths (default [vendor])
  -p, --path string    Path to search for tests.
                       If a directory is provided, all files that match '*_test.(jsonnet|libsonnet)' will be included.
                       If a file is passed, will only run the file. (default "tests")
      --version        version for jt
```

## Creating a test

The `jt` CLI tool expects tests to be formatted as:

```jsonnet
{
    pass: true // Boolean
    tests: [], // Any
    assertions: [], // Any
}
```

Each property is used to easily identify problems in generated JSON and assert that the output is correct.

A less contrived example:

```jsonnet
local myObj1 = {
    "abc": "123",
};

local myObj2 = {
    "efg": "456",
};

{
    pass: myObj1 == myObj2 // Jsonnet allows for Deep-Equal testing OOTB.
    tests: myObj1,
    assertions: myObj2,
}
```

Creating small tests like these and separating multiple files is obviously a horrible experience.

To Solve that, we can use `Arrays` to group related tests together:

```jsonnet
local fakeAWSResource1 = {
    "abc": "123",
};

local fakeAWSResource2 = {
    "efg": "456",
};

{
    pass: fakeAWSResource1 == fakeAWSResource2
    tests: [
        fakeAWSResource1,
        fakeAWSResource2,
    ],
    assertions: [
        {"abc": "123"},
        {"efg": "456"},
    ],
}
```

That looks a whole lot cleaner!
