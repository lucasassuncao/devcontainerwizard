# Customizations

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| [vscode](#vscode) | VS Code specific customizations. | object | No | - |
| [codespaces](#codespaces) | Codespaces specific customizations. | object | No | - |
| [jetbrains](#jetbrains) | JetBrains IDE specific customizations. | object | No | - |
| [neovim](#neovim) | Neovim specific customizations. | object | No | - |

### vscode

VS Code specific customizations.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| settings | Key-value settings for VS Code. | object | No | - |
| extensions | List of VS Code extensions to install. | array[string] | No | - |
| remoteUser | Remote user for VS Code container setup. | string | No | - |

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

### neovim

Neovim specific customizations.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| plugins | List of Neovim plugins to install. | array[string] | No | - |

