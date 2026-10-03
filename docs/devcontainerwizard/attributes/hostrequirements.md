# hostRequirements

## Examples

For usage examples, see [hostRequirements presets](../examples/hostrequirements.md).

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| cpus | int | Minimum number of CPU cores. | No | - |
| memory | string | Minimum memory. A whole number with an optional lowercase unit: kb, mb, gb, or tb. A bare number means bytes. | No | - |
| storage | string | Minimum free disk space, in the same format as memory. | No | - |
| gpu | bool \| "optional" \| object | GPU requirement: true, false, "optional" to use one when present, or an object with cores (at least 1) and memory. | No | - |

