# Mount

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| type | Type of mount: bind, volume, or tmpfs. | string | Yes | - |
| source | Source path or volume name. Not required for tmpfs mounts. | string | No | - |
| target | Target path inside the container. | string | Yes | - |
| readonly | Whether the mount is read-only. | boolean | No | - |

