# customizations

## Examples

For usage examples, see [customizations presets](../examples/customizations.md).

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| vscode | object | VS Code-specific settings and extensions. | No | - |
| codespaces | object | GitHub Codespaces-specific settings and extensions. | No | - |
| jetbrains | object | JetBrains IDE plugins to install. | No | - |

### vscode

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| settings | map[string]object | Default settings written into the container's workspace settings.json. Users can still override them locally. | No | - |
| extensions | []string | Extension IDs installed when the container is created, in publisher.name form. | No | - |

### codespaces

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| settings | map[string]object | Settings applied only when the container runs as a GitHub Codespace. | No | - |
| extensions | []string | Extension IDs installed only in Codespaces, on top of the vscode list. | No | - |

### jetbrains

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| plugins | []string | JetBrains plugin IDs installed when the container is created, in reverse-domain form. | No | - |

