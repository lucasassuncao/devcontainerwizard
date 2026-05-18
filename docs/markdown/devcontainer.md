# DevContainer

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| $schema | URL of the JSON schema that describes the format of this file. | string | No | - |
| name | Name of the dev container. | string | Yes | - |
| image | Docker image to use for the dev container. | string | No | - |
| [build](#build) | Configuration for building the image. | object | No | - |
| dockerFile | Deprecated: legacy path to the Dockerfile. Use build.dockerfile instead. | string | No | - |
| dockerComposeFile | List of Docker Compose files to use. | array[string] | No | - |
| service | Specific service to run from Docker Compose. | string | No | - |
| runServices | Docker Compose services to start automatically alongside the dev container service. | array[string] | No | - |
| workspaceFolder | Path to the workspace folder inside the container. | string | No | - |
| workspaceMount | Mount type for the workspace folder. | string | No | - |
| remoteUser | User that tools run as inside the container (VS Code, extensions, terminals). | string | No | - |
| containerUser | User for all processes inside the container, including the entrypoint. | string | No | - |
| updateRemoteUserUID | Sync the container user UID/GID with the local user on Linux to avoid permission issues. | boolean | No | - |
| userEnvProbe | Shell type used to probe user environment variables. | string | No | - |
| containerEnv | Environment variables to set in the container. | map[string]string | No | - |
| remoteEnv | Environment variables for remote connections (like SSH). | map[string]string | No | - |
| forwardPorts | Ports that are forwarded from the container to the local machine. Can be an integer port number, or a string of the format "host:port_number" | array[object] | No | - |
| appPort | Legacy: ports to publish from the container. Prefer forwardPorts instead. | array[object] | No | - |
| [portsAttributes](#portsattributes-value) | Additional attributes for forwarded ports. | map[string]object | No | - |
| [otherPortsAttributes](#otherportsattributes) | Default attributes applied to all forwarded ports not defined in portsAttributes. | object | No | - |
| [mounts](#mounts-item) | Mount points inside the container. Each entry can be a Mount object or a Docker --mount string. | array[object] | No | - |
| runArgs | Additional arguments to pass to 'docker run'. | array[string] | No | - |
| overrideCommand | Whether to override the container's default startup command with the devcontainer lifecycle commands. | boolean | No | - |
| init | Whether to run an init process inside the container. | boolean | No | - |
| privileged | Run the container in privileged mode. | boolean | No | - |
| capAdd | Linux capabilities to add to the container. | array[string] | No | - |
| securityOpt | Security options for the container. | array[string] | No | - |
| devices | Devices to expose to the container. | array[string] | No | - |
| [hostRequirements](#hostrequirements) | Minimum host hardware requirements for the dev container. | object | No | - |
| overrideFeatureInstallOrder | Order to install features inside the container, overriding defaults. | array[string] | No | - |
| features | Features to install in the container and their options. | map[string]object | No | - |
| initializeCommand | Command to run on the host before the container is created or started. Can be a string, an array of strings, or a named command object. | map[string]object | No | - |
| onCreateCommand | Command to run after the container is created. Can be a string, an array of strings, or a named command object. | map[string]object | No | - |
| updateContentCommand | Command to run when the container content is updated. Can be a string, an array of strings, or a named command object. | map[string]object | No | - |
| postCreateCommand | Command to run after the container is created and initialized. Can be a string, an array of strings, or a named command object. | map[string]object | No | - |
| postStartCommand | Command to run after the container starts. Can be a string, an array of strings, or a named command object. | map[string]object | No | - |
| postAttachCommand | Command to run after attaching to the container. Can be a string, an array of strings, or a named command object. | map[string]object | No | - |
| waitFor | Lifecycle command to wait for before the tool considers the container ready. | string | No | - |
| [watch](#watch) | Configuration for files/processes to watch for restarts. | object | No | - |
| [customizations](#customizations) | Editor/IDE customizations inside the container. | object | No | - |
| [secrets](#secrets-value) | Secrets to pass to the container. | map[string]object | No | - |
| shutdownAction | Action to take when the container is stopped. Use none or stopContainer (stopCompose is only valid in compose variants). | string | No | - |

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
| options | Additional CLI options passed to docker build (e.g. --no-cache). | array[string] | No | - |

### portsAttributes Value

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| label | Human-readable label for the port. | string | No | - |
| onAutoForward | Behavior when the port is auto-forwarded (notify, openBrowser, ignore). | string | No | - |
| protocol | Network protocol (http/https) for the port. | string | No | - |
| elevateIfNeeded | Prompt for elevated privileges if the port requires it (e.g. ports below 1024). | boolean | No | - |
| requireLocalPort | Require the local port to match the remote port. Shows a modal if not available. | boolean | No | - |

### otherPortsAttributes

Default attributes applied to all forwarded ports not defined in portsAttributes.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| label | Human-readable label for the port. | string | No | - |
| onAutoForward | Behavior when the port is auto-forwarded (notify, openBrowser, ignore). | string | No | - |
| protocol | Network protocol (http/https) for the port. | string | No | - |
| elevateIfNeeded | Prompt for elevated privileges if the port requires it (e.g. ports below 1024). | boolean | No | - |
| requireLocalPort | Require the local port to match the remote port. Shows a modal if not available. | boolean | No | - |

### mounts Item

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| [Mount](#mount) |  | object | No | - |
| Str |  | string | No | - |

#### Mount

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| type | Type of mount: bind, volume, or tmpfs. | string | Yes | - |
| source | Source path or volume name. Not required for tmpfs mounts. | string | No | - |
| target | Target path inside the container. | string | Yes | - |
| readonly | Whether the mount is read-only. | boolean | No | - |

### hostRequirements

Minimum host hardware requirements for the dev container.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| cpus | Minimum number of CPUs required. | integer | No | - |
| memory | Minimum memory required (e.g. "4gb"). | string | No | - |
| storage | Minimum disk storage required (e.g. "32gb"). | string | No | - |
| [gpu](#gpu) | GPU requirement: true/false, "optional", or object with cores/memory. | object | No | - |

#### gpu

GPU requirement: true/false, "optional", or object with cores/memory.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| Bool |  | boolean | No | - |
| StringVal |  | string | No | - |
| [Requirement](#requirement) |  | object | No | - |

##### Requirement

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| cores | Minimum number of GPU cores required. | integer | No | - |
| memory | Minimum GPU memory required (e.g. "4gb"). | string | No | - |

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

#### vscode

VS Code specific customizations.

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| settings | Key-value settings for VS Code. | object | No | - |
| extensions | List of VS Code extensions to install. | array[string] | No | - |
| devPort | Port on which the VS Code server listens inside the container. | integer | No | - |

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

### secrets Value

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| description | Human-readable description of the secret. | string | No | - |
| documentationUrl | URL pointing to documentation for this secret. | string | No | - |

