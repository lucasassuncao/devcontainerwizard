# watch

## Examples

For usage examples, see [watch presets](../examples/watch.md).

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| waitFor | []string | Lifecycle commands that must finish before file watching starts, so a rebuild does not race the watcher. | No | - |
| restart | []string | Glob patterns that trigger a restart of the watched process when a matching file changes. Quote them - a leading * is not valid unquoted YAML. | No | - |

