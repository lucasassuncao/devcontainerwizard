package model

// This file is the single source of field descriptions, defaults, allowed
// values, snippets, and examples for the devcontainer schema. It is consumed
// by the TUI editor's hint panel, by the FromMetadata validator family, and by
// the documentation generator - all through Metadata(), which each reader
// composes by calling it on every nested struct.
//
// Because the keys are validated against the struct's yaml tags, a typo here
// fails at startup rather than silently producing dead metadata.

// Metadata declares the top-level devcontainer fields.
//
// Two fields set Children explicitly instead of letting composition fill them:
//
//   - mounts, because MountOrString declares its own shape - composition skips
//     those (reflection cannot read their fields), so the mount attributes are
//     declared here for the hint panel.
//   - portsAttributes, because it is a map keyed by port number and needs
//     different snippets from otherPortsAttributes, which shares the same Go
//     type. Explicit children win over composition, so otherPortsAttributes
//     picks up the unkeyed PortAttributes.Metadata() and this one does not.
func (DevContainer) Metadata() map[string]any {
	return map[string]any{
		"name": meta{
			Description: "Name for the dev container, displayed to the user by the IDE.",
			Required:    true,
			Example:     "name: my-project",
		},
		"image": meta{
			Description: "Docker image used to create the container. Mutually exclusive with build and dockerComposeFile.",
			Example:     "image: mcr.microsoft.com/devcontainers/base:ubuntu",
		},
		"build": meta{
			Description: "Build the image from a Dockerfile instead of pulling one. Mutually exclusive with image and dockerComposeFile.",
			Example:     "build:\n  dockerfile: Dockerfile\n  context: .",
		},
		"dockerFile": meta{
			Description: "Path to a Dockerfile. Deprecated shorthand for build.dockerfile - prefer the build block for new configs.",
			Deprecated:  "use the build block instead: build.dockerfile",
			Example:     "dockerFile: Dockerfile",
		},
		"dockerComposeFile": meta{
			Description: "Docker Compose file(s) used to start the services. Later files override earlier ones. Mutually exclusive with image and build.",
			MinCount:    1,
			Unique:      true,
			Example:     "dockerComposeFile:\n  - docker-compose.yml\n  - docker-compose.override.yml",
		},
		"service": meta{
			Description: "Service in the Docker Compose file the IDE attaches to - the primary container of the dev environment. Required when dockerComposeFile is set.",
			Example:     "service: app",
		},
		"runServices": meta{
			Description: "Additional Compose services started and stopped alongside the primary service. Defaults to every service in the file.",
			Unique:      true,
			Example:     "runServices:\n  - db\n  - redis",
		},
		"workspaceFolder": meta{
			Description: "Path the IDE opens inside the container. Supports ${localWorkspaceFolderBasename} and the other devcontainer variables.",
			Default:     "the automatic source-code mount location; \"/\" for Docker Compose",
			Formats:     []string{"directory"},
			Example:     "workspaceFolder: /workspaces/${localWorkspaceFolderBasename}",
		},
		"workspaceMount": meta{
			Description: "Overrides the default workspace mount, as a docker run --mount string. Set it together with workspaceFolder.",
			Example:     "workspaceMount: source=${localWorkspaceFolder},target=/workspaces/my-project,type=bind,consistency=cached",
		},
		"remoteUser": meta{
			Description: "User the IDE and every process it spawns run as. Does not change the container's own entrypoint user - see containerUser.",
			Default:     "the user the container runs as, often root",
			Example:     "remoteUser: vscode",
		},
		"containerUser": meta{
			Description: "User the container itself is started with, entrypoint included.",
			Default:     "root, or the last USER instruction in the Dockerfile",
			Example:     "containerUser: vscode",
		},
		"updateRemoteUserUID": meta{
			Description: "On Linux, update the container user's UID/GID to match the local user so bind-mounted files keep the right ownership.",
			Default:     "true",
			Example:     "updateRemoteUserUID: true",
		},
		"userEnvProbe": meta{
			Description: "Shell used to probe the user's environment variables on container start. A login shell picks up /etc/profile and ~/.profile; an interactive one picks up ~/.bashrc.",
			OneOf:       []string{"none", "loginShell", "loginInteractiveShell", "interactiveShell"},
			Default:     "loginInteractiveShell",
			Example:     "userEnvProbe: loginInteractiveShell",
		},
		"containerEnv": meta{
			Description: "Environment variables baked into the container at creation time, visible to every process. Changing one requires a rebuild.",
			Example:     "containerEnv:\n  GOPATH: /go\n  GOMODCACHE: /root/.cache/go/mod",
		},
		"remoteEnv": meta{
			Description: "Environment variables for processes the IDE spawns, resolved at attach time rather than build time. The only place ${containerEnv:VAR} works; ${localEnv:VAR} reads the host environment.",
			Example:     "remoteEnv:\n  PATH: ${containerEnv:PATH}:/usr/local/mybin\n  MY_TOKEN: ${localEnv:MY_TOKEN}",
		},
		"forwardPorts": meta{
			Description: "Ports forwarded from the container to the local machine. Each entry is a port number (0-65535) or a \"host:port\" string naming another Compose service.",
			Type:        "[]int | []string",
			Unique:      true,
			Example:     "forwardPorts:\n  - 3000\n  - 5432\n  - \"db:5432\"",
		},
		"appPort": meta{
			Description: "Legacy: publishes container ports on the host, like docker run -p. Unlike forwardPorts it cannot be changed without a rebuild - prefer forwardPorts.",
			Type:        "int | string | []int | []string",
			Unique:      true,
			Example:     "appPort:\n  - 3000",
		},
		"portsAttributes": meta{
			Description: "Per-port forwarding behaviour. Keys are quoted port numbers (\"3000\"), ranges (\"40000-55000\"), or regular expressions matching a port.",
			Example:     "portsAttributes:\n  \"3000\":\n    label: Web App\n    onAutoForward: openBrowser",
			Children:    portAttributesChildren(true),
		},
		"otherPortsAttributes": meta{
			Description: "Fallback attributes for every forwarded port that portsAttributes does not match.",
			Example:     "otherPortsAttributes:\n  onAutoForward: silent",
		},
		"mounts": meta{
			Description: "Extra mounts set up when the container is created. Each entry is a Mount object or a raw docker --mount string.",
			Type:        "[]object | []string",
			Example:     "mounts:\n  - type: bind\n    source: ${localWorkspaceFolder}/.cache\n    target: /home/vscode/.cache",
			Children:    mountChildren(),
		},
		"runArgs": meta{
			Description: "Extra arguments passed to docker run. Ignored for Docker Compose configs, where the Compose file owns the runtime settings.",
			Example:     "runArgs:\n  - --network=host\n  - --cap-add=SYS_PTRACE",
		},
		"overrideCommand": meta{
			Description: "Replace the image's default command with an infinite sleep so the container stays up for the IDE to attach to.",
			Default:     "true for image/Dockerfile, false for Docker Compose",
			Example:     "overrideCommand: true",
		},
		"init": meta{
			Description: "Pass docker run --init, adding a tiny PID 1 that reaps zombie processes.",
			Default:     "false",
			Example:     "init: true",
		},
		"privileged": meta{
			Description: "Pass docker run --privileged, giving the container full access to the host. Prefer capAdd for a specific capability.",
			Default:     "false",
			Example:     "privileged: false",
		},
		"capAdd": meta{
			Description: "Linux capabilities to grant the container, such as SYS_PTRACE for attaching a debugger.",
			Unique:      true,
			Example:     "capAdd:\n  - SYS_PTRACE\n  - NET_ADMIN",
		},
		"securityOpt": meta{
			Description: "Values for docker run --security-opt, most often seccomp=unconfined so a debugger can attach.",
			Unique:      true,
			Example:     "securityOpt:\n  - seccomp=unconfined",
		},
		"devices": meta{
			Description: "Host devices to expose inside the container, as docker run --device values. A Docker passthrough, not part of the devcontainer specification.",
			Unique:      true,
			Example:     "devices:\n  - /dev/video0:/dev/video0",
		},
		"hostRequirements": meta{
			Description: "Minimum host hardware. Codespaces picks a machine type from this; locally it is a pre-flight check.",
			Example:     "hostRequirements:\n  cpus: 4\n  memory: 8gb\n  storage: 32gb",
		},
		"overrideFeatureInstallOrder": meta{
			Description: "Force the install order of features by listing their IDs. Without it the order comes from each feature's own installsAfter property.",
			Unique:      true,
			Example:     "overrideFeatureInstallOrder:\n  - ghcr.io/devcontainers/features/git:1\n  - ghcr.io/devcontainers/features/node:1",
		},
		"features": meta{
			Description: "Dev Container Features to install, keyed by feature ID. An empty object takes the feature's defaults.",
			Example:     "features:\n  ghcr.io/devcontainers/features/git:1: {}\n  ghcr.io/devcontainers/features/node:1:\n    version: \"18\"",
		},
		// Lifecycle command fields: CommandValue declares no shape and all its
		// fields are yaml:"-", so it contributes no children to compose.
		"initializeCommand": meta{
			Description: "Runs on the HOST, before the container is created or started - the only lifecycle command that does. Use it to prepare files the build needs.",
			Type:        lifecycleType,
			Example:     "initializeCommand: echo initializing\n\n# Named parallel commands:\ninitializeCommand:\n  setup: ./scripts/setup.sh\n  check: ./scripts/check.sh",
		},
		"onCreateCommand": meta{
			Description: "First command run inside the container, during creation. In prebuild-capable tools this is the last step baked into the prebuilt image.",
			Type:        lifecycleType,
			Example:     "onCreateCommand: npm install",
		},
		"updateContentCommand": meta{
			Description: "Runs after onCreateCommand, and again whenever the workspace content is refreshed. The default stage waitFor waits on.",
			Type:        lifecycleType,
			Example:     "updateContentCommand: npm install",
		},
		"postCreateCommand": meta{
			Description: "Runs once the container is created and every feature is installed - the usual place to install project dependencies.",
			Type:        lifecycleType,
			Example:     "postCreateCommand: bash .devcontainer/setup.sh",
		},
		"postStartCommand": meta{
			Description: "Runs every time the container starts, including after a stop/start cycle. Keep it quick - the IDE waits on it.",
			Type:        lifecycleType,
			Example:     "postStartCommand: nohup bash -c 'sleep infinity' &",
		},
		"postAttachCommand": meta{
			Description: "Runs each time a tool attaches to the container, so it can fire several times for one container.",
			Type:        lifecycleType,
			Example:     "postAttachCommand: source ~/.bashrc",
		},
		"waitFor": meta{
			Description: "Lifecycle stage the IDE waits for before it reports the container ready and opens the workspace.",
			OneOf:       []string{"initializeCommand", "onCreateCommand", "updateContentCommand", "postCreateCommand", "postStartCommand"},
			Default:     "updateContentCommand",
			Example:     "waitFor: postCreateCommand",
		},
		"watch": meta{
			Description: "File-watch configuration for automatic process restarts inside the container.",
			Example:     "watch:\n  waitFor:\n    - postCreateCommand\n  restart:\n    - '**/*.go'",
		},
		"customizations": meta{
			Description: "Tool-specific configuration, one sub-object per tool. Unknown tools are ignored rather than rejected.",
			Example:     "customizations:\n  vscode:\n    extensions:\n      - golang.go\n    settings:\n      editor.formatOnSave: true",
		},
		"secrets": meta{
			Description: "Secrets the container needs, keyed by environment variable name (letters, digits and underscore, not starting with a digit). The IDE prompts for the values - none of them are stored in this file.",
			Example:     "secrets:\n  MY_TOKEN:\n    description: Personal access token for the registry\n    documentationUrl: https://example.com/tokens",
		},
		"shutdownAction": meta{
			Description: "What happens to the container when the IDE disconnects. Use none to leave it running. stopCompose is the Docker Compose counterpart of stopContainer.",
			OneOf:       []string{"none", "stopContainer", "stopCompose"},
			Default:     "stopContainer for image/Dockerfile, stopCompose for Docker Compose",
			Example:     "shutdownAction: stopContainer",
		},
	}
}

// lifecycleType is the hint-panel type label for every lifecycle command.
// CommandValue is a union reflection cannot describe: a string runs in a shell,
// an array runs without shell expansion, and a map runs its entries in parallel.
const lifecycleType = "string | []string | map[string]string"

// Metadata declares the fields of the build block.
func (BuildConfig) Metadata() map[string]any {
	return map[string]any{
		"dockerfile": meta{
			Description: "Path to the Dockerfile, relative to devcontainer.json. The only required field of the build block.",
			Required:    true,
			Default:     "Dockerfile",
			Example:     "dockerfile: Dockerfile",
			PreChecked:  true,
			Snippet:     "  dockerfile: Dockerfile\n",
		},
		"context": meta{
			Description: "Directory docker build runs from, relative to devcontainer.json. Paths in COPY and ADD resolve against it.",
			Default:     "the folder containing devcontainer.json",
			Formats:     []string{"directory"},
			Example:     "context: .",
			PreChecked:  true,
			Snippet:     "  context: .\n",
		},
		"args": meta{
			Description: "Build arguments, passed as --build-arg key=value. Each one must have a matching ARG in the Dockerfile.",
			Example:     "args:\n  NODE_VERSION: \"18\"\n  APP_ENV: dev",
			Snippet:     "  args:\n    MY_ARG: value\n",
		},
		"target": meta{
			Description: "Stage to stop at in a multi-stage Dockerfile, as docker build --target.",
			Example:     "target: dev",
			Snippet:     "  target: dev\n",
		},
		"cacheFrom": meta{
			Description: "Images to pull layers from as a build cache, as docker build --cache-from. Usually a registry image a CI job pushed.",
			Unique:      true,
			Example:     "cacheFrom:\n  - myregistry/image:cache",
			Snippet:     "  cacheFrom:\n    - myregistry/image:cache\n",
		},
		"options": meta{
			Description: "Extra flags passed straight to docker build for anything the fields above do not cover.",
			Example:     "options:\n  - --no-cache\n  - --pull",
			Snippet:     "  options:\n    - --no-cache\n",
		},
	}
}

// Metadata declares the fields of a port-attributes block. It is the unkeyed
// variant used by otherPortsAttributes; see portAttributesChildren.
func (PortAttributes) Metadata() map[string]any {
	return portAttributesChildren(false)
}

// portAttributesChildren builds the child map for a port-attributes block.
//
// keyed selects the portsAttributes variant, which is a map keyed by port
// number: its first snippet has to open the port key before the attribute, and
// the label is pre-checked because an unlabelled port entry is not useful.
// The otherPortsAttributes variant applies to every unmatched port at once, so
// it needs no key and only onAutoForward is worth pre-checking.
func portAttributesChildren(keyed bool) map[string]any {
	labelSnippet, attrIndent := "  label: Other Port\n", "  "
	if keyed {
		labelSnippet, attrIndent = "  \"3000\":\n    label: Web App\n", "    "
	}

	// The two blocks do not share an onAutoForward enum: the specification
	// omits openBrowserOnce from otherPortsAttributes, because "once" has no
	// meaning for a rule that applies to every unmatched port at once.
	forwardActions := []string{"notify", "openBrowser", "openPreview", "silent", "ignore"}
	forwardDesc := "What happens when a port is auto-forwarded: notify shows a toast, openBrowser opens the system browser, openPreview uses the IDE's built-in browser, silent forwards quietly, and ignore skips the port entirely."
	if keyed {
		forwardActions = []string{"notify", "openBrowser", "openBrowserOnce", "openPreview", "silent", "ignore"}
		forwardDesc = "What happens when this port is auto-forwarded: notify shows a toast, openBrowser opens the system browser every time, openBrowserOnce only the first time in a session, openPreview uses the IDE's built-in browser, silent forwards quietly, and ignore skips the port entirely."
	}

	return map[string]any{
		"label": meta{
			Description: "Name shown next to the port in the IDE's forwarded-ports list.",
			Default:     "Application",
			Example:     "label: Web App",
			PreChecked:  keyed,
			Snippet:     labelSnippet,
		},
		"onAutoForward": meta{
			Description: forwardDesc,
			OneOf:       forwardActions,
			Default:     "notify",
			Example:     "onAutoForward: notify",
			PreChecked:  true,
			Snippet:     attrIndent + "onAutoForward: notify\n",
		},
		"protocol": meta{
			Description: "Scheme the IDE uses when it builds the URL for this port.",
			OneOf:       []string{"http", "https"},
			Example:     "protocol: http",
			Snippet:     attrIndent + "protocol: http\n",
		},
		"elevateIfNeeded": meta{
			Description: "Ask for elevated privileges when the local port needs them - ports below 1024 such as 22, 80 and 443.",
			Default:     "false",
			Example:     "elevateIfNeeded: false",
			Snippet:     attrIndent + "elevateIfNeeded: false\n",
		},
		"requireLocalPort": meta{
			Description: "Insist on the same port number locally. When false the IDE silently maps to a free port instead; when true it warns that the port is taken.",
			Default:     "false",
			Example:     "requireLocalPort: false",
			Snippet:     attrIndent + "requireLocalPort: false\n",
		},
	}
}

// mountChildren describes one entry of the mounts list. MountOrString declares
// its own shape, so metadata composition skips it and these entries are
// attached to the mounts field by hand; they are not validated by reflection.
func mountChildren() map[string]any {
	return map[string]any{
		"type": meta{
			// The specification's enum is bind|volume. tmpfs is accepted here as
			// well because Docker supports it and rejecting it would break configs
			// that already work; being more permissive than the spec is safe,
			// being stricter is not.
			Description: "Mount type. bind maps a host path, volume uses a named Docker volume that survives rebuilds. tmpfs is a Docker extension, not part of the specification.",
			Required:    true,
			OneOf:       []string{"bind", "volume", "tmpfs"},
			Example:     "type: bind",
			PreChecked:  true,
			Snippet:     "  - type: bind\n",
		},
		"source": meta{
			Description: "Host path for a bind mount, or the volume name for a volume mount. Omit it for tmpfs.",
			Example:     "source: ${localWorkspaceFolder}/.cache",
			Snippet:     "    source: ${localWorkspaceFolder}/.cache\n",
		},
		"target": meta{
			Description: "Absolute path the mount appears at inside the container.",
			Required:    true,
			Formats:     []string{"directory"},
			Example:     "target: /home/vscode/.cache",
			PreChecked:  true,
			Snippet:     "    target: /home/vscode/.cache\n",
		},
	}
}

// sizePattern is the devcontainer specification's regex for hostRequirements
// memory and storage: a whole number with an optional lowercase unit suffix.
// Uppercase ("8GB") is rejected by the specification, so it is rejected here.
const sizePattern = `^\d+([tgmk]b)?$`

// Metadata declares the fields of the hostRequirements block.
func (HostRequirements) Metadata() map[string]any {
	return map[string]any{
		"cpus": meta{
			Description: "Minimum number of CPU cores.",
			Min:         "1",
			Example:     "cpus: 4",
			PreChecked:  true,
			Snippet:     "  cpus: 4\n",
		},
		"memory": meta{
			Description: "Minimum memory. A whole number with an optional lowercase unit: kb, mb, gb, or tb. A bare number means bytes.",
			Pattern:     sizePattern,
			Example:     "memory: 8gb",
			PreChecked:  true,
			Snippet:     "  memory: 8gb\n",
		},
		"storage": meta{
			Description: "Minimum free disk space, in the same format as memory.",
			Pattern:     sizePattern,
			Example:     "storage: 32gb",
			Snippet:     "  storage: 32gb\n",
		},
		// gpu is *GPUValue: a union with no yaml-tagged fields, so it declares no
		// metadata of its own and stays a leaf in the hint panel.
		"gpu": meta{
			Description: "GPU requirement: true, false, \"optional\" to use one when present, or an object with cores (at least 1) and memory.",
			Type:        "bool | \"optional\" | object",
			Example:     "gpu: true\n\n# Or with specific requirements:\ngpu:\n  cores: 2\n  memory: 4gb",
			Snippet:     "  gpu: true\n",
		},
	}
}

// Metadata declares the fields of the watch block.
func (WatchConfig) Metadata() map[string]any {
	return map[string]any{
		"waitFor": meta{
			Description: "Lifecycle commands that must finish before file watching starts, so a rebuild does not race the watcher.",
			Unique:      true,
			Example:     "waitFor:\n  - postCreateCommand",
			PreChecked:  true,
			Snippet:     "  waitFor:\n    - postCreateCommand\n",
		},
		"restart": meta{
			Description: "Glob patterns that trigger a restart of the watched process when a matching file changes. Quote them - a leading * is not valid unquoted YAML.",
			Unique:      true,
			Example:     "restart:\n  - '**/*.go'\n  - '**/*.ts'",
			Snippet:     "  restart:\n    - '**/*.go'\n",
		},
	}
}

// Metadata declares the fields of the customizations block.
func (Customizations) Metadata() map[string]any {
	return map[string]any{
		"vscode": meta{
			Description: "VS Code-specific settings and extensions.",
			Example:     "vscode:\n  extensions:\n    - golang.go\n  settings:\n    editor.formatOnSave: true",
			Snippet:     "  vscode:\n    extensions:\n      - ms-python.python\n    settings:\n      editor.formatOnSave: true\n",
		},
		"codespaces": meta{
			Description: "GitHub Codespaces-specific settings and extensions.",
			Example:     "codespaces:\n  extensions:\n    - github.copilot",
			Snippet:     "  codespaces:\n    extensions:\n      - github.copilot\n    settings:\n      editor.formatOnSave: true\n",
		},
		"jetbrains": meta{
			Description: "JetBrains IDE plugins to install.",
			Example:     "jetbrains:\n  plugins:\n    - org.rust.lang",
			Snippet:     "  jetbrains:\n    plugins:\n      - org.rust.lang\n",
		},
	}
}

// Metadata declares the fields of customizations.vscode.
func (VSCodeCustomization) Metadata() map[string]any {
	return map[string]any{
		"settings": meta{
			Description: "Default settings written into the container's workspace settings.json. Users can still override them locally.",
			Example:     "settings:\n  editor.formatOnSave: true\n  go.toolsManagement.autoUpdate: true",
		},
		"extensions": meta{
			Description: "Extension IDs installed when the container is created, in publisher.name form.",
			Unique:      true,
			Example:     "extensions:\n  - golang.go\n  - ms-python.python",
		},
	}
}

// Metadata declares the fields of customizations.codespaces.
func (CodespacesCustomization) Metadata() map[string]any {
	return map[string]any{
		"settings": meta{
			Description: "Settings applied only when the container runs as a GitHub Codespace.",
			Example:     "settings:\n  editor.formatOnSave: true",
		},
		"extensions": meta{
			Description: "Extension IDs installed only in Codespaces, on top of the vscode list.",
			Unique:      true,
			Example:     "extensions:\n  - github.copilot",
		},
	}
}

// Metadata declares the fields of customizations.jetbrains.
func (JetBrainsCustomization) Metadata() map[string]any {
	return map[string]any{
		"plugins": meta{
			Description: "JetBrains plugin IDs installed when the container is created, in reverse-domain form.",
			Unique:      true,
			Example:     "plugins:\n  - org.rust.lang\n  - com.intellij.kubernetes",
		},
	}
}

// Metadata declares the fields of one entry of the secrets map.
func (Secret) Metadata() map[string]any {
	return map[string]any{
		// Neither field is required by the specification, but a secret with no
		// description gives the user a bare variable name to fill in, so it is
		// pre-checked rather than enforced.
		"description": meta{
			Description: "What the secret is for. Shown next to the input when the IDE prompts for its value.",
			Example:     "description: Personal access token for the registry",
			PreChecked:  true,
			Snippet:     "  MY_SECRET:\n    description: Description of the secret\n",
		},
		"documentationUrl": meta{
			Description: "Link shown alongside the prompt, pointing at instructions for obtaining the value.",
			Formats:     []string{"url"},
			Example:     "documentationUrl: https://example.com/docs/tokens",
			Snippet:     "    documentationUrl: https://example.com/docs/secrets\n",
		},
	}
}
