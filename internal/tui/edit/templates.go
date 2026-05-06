package edit

// guidedTemplates maps each known top-level key to a YAML snippet that shows
// every available field with example values. Optional fields are commented out
// so the user can uncomment what they need and delete what they don't.
var guidedTemplates = map[string]string{ // #nosec G101 -- YAML example templates, not real credentials
	"name": `name: my-devcontainer
`,

	"image": `image: ubuntu:22.04
`,

	"build": `build:
  dockerfile: Dockerfile
  context: .
  # args:
  #   MY_ARG: value
  # target: dev
  # cacheFrom:
  #   - myregistry/myimage:cache
  # output: type=local,dest=./out
  # ssh:
  #   - default
  # secrets:
  #   - id=mysecret,src=/path/to/secret
`,

	"dockerComposeFile": `dockerComposeFile:
  - docker-compose.yml
  # - docker-compose.override.yml
`,

	"service": `service: app
`,

	"runServices": `runServices:
  - db
  - redis
  # - queue
`,

	"workspaceFolder": `workspaceFolder: /workspace
`,

	"workspaceMount": `workspaceMount: source=${localWorkspaceFolder},target=/workspace,type=bind,consistency=cached
`,

	"remoteUser": `remoteUser: vscode
`,

	"containerUser": `containerUser: vscode
`,

	"updateRemoteUserUID": `updateRemoteUserUID: true
`,

	"userEnvProbe": `# Options: none | loginShell | loginInteractiveShell | interactiveShell
userEnvProbe: loginInteractiveShell
`,

	"containerEnv": `containerEnv:
  MY_VAR: value
  # ANOTHER_VAR: other-value
`,

	"remoteEnv": `remoteEnv:
  MY_REMOTE_VAR: value
  # PATH: ${containerEnv:PATH}:/my/custom/bin
`,

	"localEnv": `localEnv:
  MY_LOCAL_VAR: ${env:MY_LOCAL_VAR}
`,

	"appPort": `# Legacy field — prefer forwardPorts instead.
appPort:
  - 3000
  # - 8080
`,

	"forwardPorts": `forwardPorts:
  - 3000
  # - 5432
  # - "host:8080"
`,

	"portsAttributes": `portsAttributes:
  "3000":
    label: Web App
    onAutoForward: notify
    # protocol: http
  # "5432":
  #   label: PostgreSQL
  #   onAutoForward: silent
`,

	"otherPortsAttributes": `otherPortsAttributes:
  onAutoForward: silent
  # label: Other Port
`,

	"mounts": `mounts:
  - type: bind
    source: ${localWorkspaceFolder}/.cache
    target: /home/vscode/.cache
    # consistency: cached
    # readonly: false
  # - type: volume
  #   source: myvolume
  #   target: /data
`,

	"runArgs": `runArgs:
  - "--network=host"
  # - "--cap-add=SYS_PTRACE"
  # - "--security-opt=seccomp=unconfined"
`,

	"startupCommand": `startupCommand: "echo 'Container started'"
`,

	"overrideCommand": `overrideCommand: true
`,

	"command": `command: sleep infinity
`,

	"entrypoint": `entrypoint: /usr/local/bin/docker-entrypoint.sh
`,

	"init": `init: true
`,

	"privileged": `privileged: false
`,

	"capAdd": `capAdd:
  - SYS_PTRACE
  # - NET_ADMIN
`,

	"capDrop": `capDrop:
  - ALL
`,

	"securityOpt": `securityOpt:
  - seccomp=unconfined
  # - apparmor=unconfined
`,

	"devices": `devices:
  - /dev/net/tun
`,

	"hostRequirements": `hostRequirements:
  cpus: 4
  memory: 8gb
  storage: 32gb
  # gpu: true
  # gpu:
  #   cores: 4
  #   memory: 4gb
`,

	"features": `features:
  ghcr.io/devcontainers/features/git:1: {}
  # ghcr.io/devcontainers/features/node:1:
  #   version: lts
  # ghcr.io/devcontainers/features/docker-in-docker:2:
  #   version: latest
`,

	"overrideFeatureInstallOrder": `overrideFeatureInstallOrder:
  - ghcr.io/devcontainers/features/git:1
`,

	"initializeCommand": `initializeCommand: echo 'Initializing on host'
# initializeCommand:
#   - /bin/sh
#   - -c
#   - echo 'Initializing on host'
`,

	"updateContentCommand": `updateContentCommand: echo 'Content updated'
# updateContentCommand:
#   - /bin/sh
#   - -c
#   - pip install -r requirements.txt
`,

	"waitFor": `# Options: initializeCommand | onCreateCommand | updateContentCommand | postCreateCommand | postStartCommand
waitFor: updateContentCommand
`,

	"onCreateCommand": `onCreateCommand: echo 'Container created'
# onCreateCommand:
#   - /bin/sh
#   - -c
#   - echo 'Container created'
`,

	"postCreateCommand": `postCreateCommand: pip install -r requirements.txt
# postCreateCommand:
#   - /bin/sh
#   - -c
#   - pip install -r requirements.txt
`,

	"postStartCommand": `postStartCommand: echo 'Container started'
`,

	"postAttachCommand": `postAttachCommand: echo 'Attached to container'
`,

	"watch": `watch:
  waitFor: postCreateCommand
  restart: true
`,

	"customizations": `customizations:
  vscode:
    extensions:
      - ms-python.python
      # - esbenp.prettier-vscode
    settings:
      editor.formatOnSave: true
      # terminal.integrated.shell.linux: /bin/zsh
  # jetbrains:
  #   plugins:
  #     - org.rust.lang
  # codespaces:
  #   openFiles:
  #     - README.md
`,

	"secrets": `secrets:
  MY_SECRET:
    description: "Description of the secret"
    # default: ""
`,

	"shutdownAction": `# Options: none | stopContainer | stopCompose
shutdownAction: stopContainer
`,
}

// GuidedTemplate returns the guided YAML snippet for a key, or a minimal
// fallback if no template is defined for that key.
func GuidedTemplate(key string) string {
	if t, ok := guidedTemplates[key]; ok {
		return t
	}
	return key + ": \n"
}
