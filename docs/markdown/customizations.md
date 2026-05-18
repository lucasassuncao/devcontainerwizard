# Customizations

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| [vscode](#vscode) | VS Code specific customizations. | object | No | - |
| [codespaces](#codespaces) | Codespaces specific customizations. | object | No | - |
| [jetbrains](#jetbrains) | JetBrains IDE specific customizations. | object | No | - |

### vscode

VS Code specific customizations.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| settings | Key-value settings for VS Code. | object | No | - |
| extensions | List of VS Code extensions to install. | array[string] | No | - |
| devPort | Port on which the VS Code server listens inside the container. | integer | No | - |

### codespaces

Codespaces specific customizations.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| settings | Key-value settings for Codespaces. | object | No | - |
| extensions | List of Codespaces extensions to install. | array[string] | No | - |

### jetbrains

JetBrains IDE specific customizations.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| plugins | List of JetBrains plugins to install. | array[string] | No | - |

