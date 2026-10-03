# mounts

## Examples

For usage examples, see [mounts presets](../examples/mounts.md).

## Arguments

The following arguments are supported:

| Name | Type | Format | Description | Required | Default |
|------|------|--------|-------------|----------|---------|
| source | string | - | Host path for a bind mount, or the volume name for a volume mount. Omit it for tmpfs. | No | - |
| target | string | directory | Absolute path the mount appears at inside the container. | Yes | - |
| type | string | - | Mount type. bind maps a host path, volume uses a named Docker volume that survives rebuilds. tmpfs is a Docker extension, not part of the specification. | Yes | - |

