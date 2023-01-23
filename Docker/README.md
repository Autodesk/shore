# Shore Dev Images

These images are meant for development purposes.

Contents of each of the images are:

* `full` - `shore`, `jsonnet-bundler` (`jb`), `jsonnet`, `jsonnetfmt`, `spin`. Published.
* `local` - full, but also `shore` is built from local source. Not published.

The name of the images is `shore-dev`. Each image is tagged with the shore's version and suffix based on the image, for instance `v0.0.10-full`.

Images can be pulled from:
```shell
docker pull ***REMOVED***/shore/shore-dev:v0.0.10-full
```

## Local

`local.Dockerfile` is meant to be built locally. At the root of shore's repo, run:
```shell
docker build -t local-shore --file Docker/local.Dockerfile .
```
And then to start it:
```shell
docker run --name local-shore -it local-shore /bin/sh
````