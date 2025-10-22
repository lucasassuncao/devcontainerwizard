# DevContainer

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| $schema | URL of the JSON schema that describes the format of this file. | string | No | - |
| name | Name of the dev container. | string | Yes | - |
| image | Docker image to use for the dev container. | string | No | - |
| [build](#build) | Configuration for building the image. | object | No | - |
| dockerComposeFile | List of Docker Compose files to use. | array[string] | No | - |
| service | Specific service to run from Docker Compose. | string | No | - |
| workspaceFolder | Path to the workspace folder inside the container. | string | No | - |
| workspaceMount | Mount type for the workspace folder. | string | No | - |
| remoteUser | User to use inside the container. | string | No | - |
| userEnvProbe | Command to detect the default user inside the container. | string | No | - |
| containerEnv | Environment variables to set in the container. | map[string]string | No | - |
| remoteEnv | Environment variables for remote connections (like SSH). | map[string]string | No | - |
| forwardPorts | Ports that are forwarded from the container to the local machine. Can be an integer port number, or a string of the format "host:port_number" | array[object] | No | - |
| [portsAttributes](#portsattributes-value) | Additional attributes for forwarded ports. | map[string]object | No | - |
| [otherPortsAttributes](#otherportsattributes) | Default attributes applied to all forwarded ports not defined in portsAttributes. | object | No | - |
| [mounts](#mounts-item) | Mount points inside the container. | array[object] | No | - |
| runArgs | Additional arguments to pass to 'docker run'. | array[string] | No | - |
| startupCommand | Command to run on container startup. | string | No | - |
| command | Command to run inside the container instead of the default CMD. | string | No | - |
| entrypoint | Entrypoint to override in the container. | string | No | - |
| init | Whether to run an init process inside the container. | boolean | No | - |
| privileged | Run the container in privileged mode. | boolean | No | - |
| capAdd | Linux capabilities to add to the container. | array[string] | No | - |
| capDrop | Linux capabilities to drop from the container. | array[string] | No | - |
| securityOpt | Security options for the container. | array[string] | No | - |
| devices | Devices to expose to the container. | array[string] | No | - |
| overrideFeatureInstallOrder | Order to install features inside the container, overriding defaults. | array[string] | No | - |
| features | Features to install in the container and their options. | map[string]object | No | - |
| onCreateCommand | Command to run after the container is created. | string | No | - |
| postCreateCommand | Command to run after the container is created and initialized. | string | No | - |
| postStartCommand | Command to run after the container starts. | string | No | - |
| postAttachCommand | Command to run after attaching to the container. | string | No | - |
| [watch](#watch) | Configuration for files/processes to watch for restarts. | object | No | - |
| [customizations](#customizations) | Editor/IDE customizations inside the container. | object | No | - |
| [secrets](#secrets-value) | Secrets to pass to the container. | map[string]object | No | - |
| shutdownAction | Action to take when the container is stopped. | string | No | - |

### build

Configuration for building the image.

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

#### secrets Item

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| id | Identifier for the secret. | string | No | - |
| src | Path or source of the secret. | string | No | - |

### portsAttributes Value

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| label | Human-readable label for the port. | string | No | - |
| onAutoForward | Behavior when the port is auto-forwarded (notify, openBrowser, ignore). | string | No | - |
| protocol | Network protocol (tcp/udp) for the port. | string | No | - |

### otherPortsAttributes

Default attributes applied to all forwarded ports not defined in portsAttributes.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| label | Human-readable label for the port. | string | No | - |
| onAutoForward | Behavior when the port is auto-forwarded (notify, openBrowser, ignore). | string | No | - |
| protocol | Network protocol (tcp/udp) for the port. | string | No | - |

### mounts Item

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| type | Type of mount (e.g., bind, volume). | string | Yes | - |
| source | Source path of the mount. | string | Yes | - |
| target | Target path inside the container. | string | Yes | - |
| consistency | Consistency mode for the mount (e.g., cached, delegated, consistent). | string | No | - |
| readonly | Whether the mount is read-only. | boolean | No | - |

### watch

Configuration for files/processes to watch for restarts.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| waitFor | List of processes or files to wait for before starting. | array[string] | No | - |
| restart | List of files or events that trigger a restart. | array[string] | No | - |

### customizations

Editor/IDE customizations inside the container.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| [vscode](#vscode) | VS Code specific customizations. | object | No | - |
| [codespaces](#codespaces) | Codespaces specific customizations. | object | No | - |
| [jetbrains](#jetbrains) | JetBrains IDE specific customizations. | object | No | - |
| [neovim](#neovim) | Neovim specific customizations. | object | No | - |

#### vscode

VS Code specific customizations.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| settings | Key-value settings for VS Code. | object | No | - |
| extensions | List of VS Code extensions to install. | array[string] | No | - |
| remoteUser | Remote user for VS Code container setup. | string | No | - |

#### codespaces

Codespaces specific customizations.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| settings | Key-value settings for Codespaces. | object | No | - |
| extensions | List of Codespaces extensions to install. | array[string] | No | - |

#### jetbrains

JetBrains IDE specific customizations.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| plugins | List of JetBrains plugins to install. | array[string] | No | - |

#### neovim

Neovim specific customizations.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| plugins | List of Neovim plugins to install. | array[string] | No | - |

### secrets Value

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| description | Human-readable description of the secret. | string | No | - |
| default | Default value for the secret if none is provided. | string | No | - |

