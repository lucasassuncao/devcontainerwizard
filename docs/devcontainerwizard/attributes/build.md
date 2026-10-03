# build

## Examples

For usage examples, see [build presets](../examples/build.md).

## Arguments

The following arguments are supported:

| Name | Type | Format | Description | Required | Default |
|------|------|--------|-------------|----------|---------|
| dockerfile | string | - | Path to the Dockerfile, relative to devcontainer.json. The only required field of the build block. | Yes | Dockerfile |
| context | string | directory | Directory docker build runs from, relative to devcontainer.json. Paths in COPY and ADD resolve against it. | No | the folder containing devcontainer.json |
| args | map[string]string | - | Build arguments, passed as --build-arg key=value. Each one must have a matching ARG in the Dockerfile. | No | - |
| target | string | - | Stage to stop at in a multi-stage Dockerfile, as docker build --target. | No | - |
| cacheFrom | []string | - | Images to pull layers from as a build cache, as docker build --cache-from. Usually a registry image a CI job pushed. | No | - |
| options | []string | - | Extra flags passed straight to docker build for anything the fields above do not cover. | No | - |

