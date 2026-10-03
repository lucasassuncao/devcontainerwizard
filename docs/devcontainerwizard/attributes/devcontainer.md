# DevContainer

## Examples

For usage examples, see [DevContainer presets](../examples/devcontainer.md).

## Arguments

The following arguments are supported:

| Name | Type | Format | Description | Required | Default |
|------|------|--------|-------------|----------|---------|
| name | string | - | Name for the dev container, displayed to the user by the IDE. | Yes | - |
| image | string | - | Docker image used to create the container. Mutually exclusive with build and dockerComposeFile. | No | - |
| [build](./build.md) | object | - | Build the image from a Dockerfile instead of pulling one. Mutually exclusive with image and dockerComposeFile. | No | - |
| dockerFile | string | - | Path to a Dockerfile. Deprecated shorthand for build.dockerfile - prefer the build block for new configs. | No | - |
| dockerComposeFile | []string | - | Docker Compose file(s) used to start the services. Later files override earlier ones. Mutually exclusive with image and build. | No | - |
| service | string | - | Service in the Docker Compose file the IDE attaches to - the primary container of the dev environment. Required when dockerComposeFile is set. | No | - |
| runServices | []string | - | Additional Compose services started and stopped alongside the primary service. Defaults to every service in the file. | No | - |
| workspaceFolder | string | directory | Path the IDE opens inside the container. Supports ${localWorkspaceFolderBasename} and the other devcontainer variables. | No | the automatic source-code mount location; "/" for Docker Compose |
| workspaceMount | string | - | Overrides the default workspace mount, as a docker run --mount string. Set it together with workspaceFolder. | No | - |
| remoteUser | string | - | User the IDE and every process it spawns run as. Does not change the container's own entrypoint user - see containerUser. | No | the user the container runs as, often root |
| containerUser | string | - | User the container itself is started with, entrypoint included. | No | root, or the last USER instruction in the Dockerfile |
| updateRemoteUserUID | bool | - | On Linux, update the container user's UID/GID to match the local user so bind-mounted files keep the right ownership. | No | true |
| userEnvProbe | string | - | Shell used to probe the user's environment variables on container start. A login shell picks up /etc/profile and ~/.profile; an interactive one picks up ~/.bashrc. | No | loginInteractiveShell |
| containerEnv | map[string]string | - | Environment variables baked into the container at creation time, visible to every process. Changing one requires a rebuild. | No | - |
| remoteEnv | map[string]string | - | Environment variables for processes the IDE spawns, resolved at attach time rather than build time. The only place ${containerEnv:VAR} works; ${localEnv:VAR} reads the host environment. | No | - |
| forwardPorts | []int \| []string | - | Ports forwarded from the container to the local machine. Each entry is a port number (0-65535) or a "host:port" string naming another Compose service. | No | - |
| appPort | int \| string \| []int \| []string | - | Legacy: publishes container ports on the host, like docker run -p. Unlike forwardPorts it cannot be changed without a rebuild - prefer forwardPorts. | No | - |
| [portsAttributes](./portsattributes.md) | map[string]object | - | Per-port forwarding behaviour. Keys are quoted port numbers ("3000"), ranges ("40000-55000"), or regular expressions matching a port. | No | - |
| [otherPortsAttributes](./otherportsattributes.md) | object | - | Fallback attributes for every forwarded port that portsAttributes does not match. | No | - |
| [mounts](./mounts.md) | []object \| []string | - | Extra mounts set up when the container is created. Each entry is a Mount object or a raw docker --mount string. | No | - |
| runArgs | []string | - | Extra arguments passed to docker run. Ignored for Docker Compose configs, where the Compose file owns the runtime settings. | No | - |
| overrideCommand | bool | - | Replace the image's default command with an infinite sleep so the container stays up for the IDE to attach to. | No | true for image/Dockerfile, false for Docker Compose |
| init | bool | - | Pass docker run --init, adding a tiny PID 1 that reaps zombie processes. | No | false |
| privileged | bool | - | Pass docker run --privileged, giving the container full access to the host. Prefer capAdd for a specific capability. | No | false |
| capAdd | []string | - | Linux capabilities to grant the container, such as SYS_PTRACE for attaching a debugger. | No | - |
| securityOpt | []string | - | Values for docker run --security-opt, most often seccomp=unconfined so a debugger can attach. | No | - |
| devices | []string | - | Host devices to expose inside the container, as docker run --device values. A Docker passthrough, not part of the devcontainer specification. | No | - |
| [hostRequirements](./hostrequirements.md) | object | - | Minimum host hardware. Codespaces picks a machine type from this; locally it is a pre-flight check. | No | - |
| overrideFeatureInstallOrder | []string | - | Force the install order of features by listing their IDs. Without it the order comes from each feature's own installsAfter property. | No | - |
| features | map[string]object | - | Dev Container Features to install, keyed by feature ID. An empty object takes the feature's defaults. | No | - |
| initializeCommand | string \| []string \| map[string]string | - | Runs on the HOST, before the container is created or started - the only lifecycle command that does. Use it to prepare files the build needs. | No | - |
| onCreateCommand | string \| []string \| map[string]string | - | First command run inside the container, during creation. In prebuild-capable tools this is the last step baked into the prebuilt image. | No | - |
| updateContentCommand | string \| []string \| map[string]string | - | Runs after onCreateCommand, and again whenever the workspace content is refreshed. The default stage waitFor waits on. | No | - |
| postCreateCommand | string \| []string \| map[string]string | - | Runs once the container is created and every feature is installed - the usual place to install project dependencies. | No | - |
| postStartCommand | string \| []string \| map[string]string | - | Runs every time the container starts, including after a stop/start cycle. Keep it quick - the IDE waits on it. | No | - |
| postAttachCommand | string \| []string \| map[string]string | - | Runs each time a tool attaches to the container, so it can fire several times for one container. | No | - |
| waitFor | string | - | Lifecycle stage the IDE waits for before it reports the container ready and opens the workspace. | No | updateContentCommand |
| [watch](./watch.md) | object | - | File-watch configuration for automatic process restarts inside the container. | No | - |
| [customizations](./customizations.md) | object | - | Tool-specific configuration, one sub-object per tool. Unknown tools are ignored rather than rejected. | No | - |
| [secrets](./secrets.md) | map[string]object | - | Secrets the container needs, keyed by environment variable name (letters, digits and underscore, not starting with a digit). The IDE prompts for the values - none of them are stored in this file. | No | - |
| shutdownAction | string | - | What happens to the container when the IDE disconnects. Use none to leave it running. stopCompose is the Docker Compose counterpart of stopContainer. | No | stopContainer for image/Dockerfile, stopCompose for Docker Compose |

