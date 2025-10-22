# Mount

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| type | Type of mount (e.g., bind, volume). | string | Yes | - |
| source | Source path of the mount. | string | Yes | - |
| target | Target path inside the container. | string | Yes | - |
| consistency | Consistency mode for the mount (e.g., cached, delegated, consistent). | string | No | - |
| readonly | Whether the mount is read-only. | boolean | No | - |

