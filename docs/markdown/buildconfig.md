# BuildConfig

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| dockerfile | Path to the Dockerfile to use for building the image. | string | Yes | Dockerfile |
| context | Build context directory. | string | Yes | . |
| args | Build arguments as key-value pairs. | map[string]string | No | - |
| target | Target stage for multi-stage Docker builds. | string | No | - |
| cacheFrom | List of images to cache from. | array[string] | No | - |
| options | Additional CLI options passed to docker build (e.g. --no-cache). | array[string] | No | - |

