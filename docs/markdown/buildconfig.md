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
| output | Output location of the build. | string | No | - |
| ssh | SSH mount sources to use during build. | array[string] | No | - |
| [secrets](#secrets-item) | Secrets to pass to the build process. | array[object] | No | - |

### secrets Item

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| id | Identifier for the secret. | string | No | - |
| src | Path or source of the secret. | string | No | - |

