package edit

const (
	tmplName = `name: my-devcontainer
`
	tmplImage = `image: ubuntu:22.04
`
	tmplBuild = `build:
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
`
	tmplDockerComposeFile = `dockerComposeFile:
  - docker-compose.yml
  # - docker-compose.override.yml
`
	tmplService = `service: app
`
	tmplRunServices = `runServices:
  - db
  - redis
  # - queue
`
	tmplWorkspaceFolder = `workspaceFolder: /workspace
`
	tmplWorkspaceMount = `workspaceMount: source=${localWorkspaceFolder},target=/workspace,type=bind,consistency=cached
`
	tmplRemoteUser = `remoteUser: vscode
`
	tmplContainerUser = `containerUser: vscode
`
	tmplUpdateRemoteUserUID = `updateRemoteUserUID: true
`
	tmplUserEnvProbe = `# Options: none | loginShell | loginInteractiveShell | interactiveShell
userEnvProbe: loginInteractiveShell
`
	tmplContainerEnv = `containerEnv:
  MY_VAR: value
  # ANOTHER_VAR: other-value
`
	tmplRemoteEnv = `remoteEnv:
  MY_REMOTE_VAR: value
  # PATH: ${containerEnv:PATH}:/my/custom/bin
`
	tmplLocalEnv = `localEnv:
  MY_LOCAL_VAR: ${env:MY_LOCAL_VAR}
`
	tmplAppPort = `# Legacy field — prefer forwardPorts instead.
appPort:
  - 3000
  # - 8080
`
	tmplForwardPorts = `forwardPorts:
  - 3000
  # - 5432
  # - "host:8080"
`
	tmplPortsAttributes = `portsAttributes:
  "3000":
    label: Web App
    onAutoForward: notify
    # protocol: http
  # "5432":
  #   label: PostgreSQL
  #   onAutoForward: silent
`
	tmplOtherPortsAttributes = `otherPortsAttributes:
  onAutoForward: silent
  # label: Other Port
`
	tmplMounts = `mounts:
  - type: bind
    source: ${localWorkspaceFolder}/.cache
    target: /home/vscode/.cache
    # consistency: cached
    # readonly: false
  # - type: volume
  #   source: myvolume
  #   target: /data
`
	tmplRunArgs = `runArgs:
  - "--network=host"
  # - "--cap-add=SYS_PTRACE"
  # - "--security-opt=seccomp=unconfined"
`
	tmplStartupCommand = `startupCommand: "echo 'Container started'"
`
	tmplOverrideCommand = `overrideCommand: true
`
	tmplCommand = `command: sleep infinity
`
	tmplEntrypoint = `entrypoint: /usr/local/bin/docker-entrypoint.sh
`
	tmplInit = `init: true
`
	tmplPrivileged = `privileged: false
`
	tmplCapAdd = `capAdd:
  - SYS_PTRACE
  # - NET_ADMIN
`
	tmplCapDrop = `capDrop:
  - ALL
`
	tmplSecurityOpt = `securityOpt:
  - seccomp=unconfined
  # - apparmor=unconfined
`
	tmplDevices = `devices:
  - /dev/net/tun
`
	tmplHostRequirements = `hostRequirements:
  cpus: 4
  memory: 8gb
  storage: 32gb
  # gpu: true
  # gpu:
  #   cores: 4
  #   memory: 4gb
`
	tmplFeatures = `features:
  ghcr.io/devcontainers/features/git:1: {}
  # ghcr.io/devcontainers/features/node:1:
  #   version: lts
  # ghcr.io/devcontainers/features/docker-in-docker:2:
  #   version: latest
`
	tmplOverrideFeatureInstallOrder = `overrideFeatureInstallOrder:
  - ghcr.io/devcontainers/features/git:1
`
	tmplInitializeCommand = `initializeCommand: echo 'Initializing on host'
# initializeCommand:
#   - /bin/sh
#   - -c
#   - echo 'Initializing on host'
`
	tmplUpdateContentCommand = `updateContentCommand: echo 'Content updated'
# updateContentCommand:
#   - /bin/sh
#   - -c
#   - pip install -r requirements.txt
`
	tmplWaitFor = `# Options: initializeCommand | onCreateCommand | updateContentCommand | postCreateCommand | postStartCommand
waitFor: updateContentCommand
`
	tmplOnCreateCommand = `onCreateCommand: echo 'Container created'
# onCreateCommand:
#   - /bin/sh
#   - -c
#   - echo 'Container created'
`
	tmplPostCreateCommand = `postCreateCommand: pip install -r requirements.txt
# postCreateCommand:
#   - /bin/sh
#   - -c
#   - pip install -r requirements.txt
`
	tmplPostStartCommand = `postStartCommand: echo 'Container started'
`
	tmplPostAttachCommand = `postAttachCommand: echo 'Attached to container'
`
	tmplWatch = `watch:
  waitFor: postCreateCommand
  restart: true
`
	tmplCustomizations = `customizations:
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
`
	tmplSecrets = `secrets:
  MY_SECRET:
    description: "Description of the secret"
    # default: ""
` // #nosec G101 -- YAML example template, not a real credential
	tmplShutdownAction = `# Options: none | stopContainer | stopCompose
shutdownAction: stopContainer
`
)

// #nosec G101 -- YAML example templates, not real credentials
var templates = map[string]string{
	"name":                        tmplName,
	"image":                       tmplImage,
	"build":                       tmplBuild,
	"dockerComposeFile":           tmplDockerComposeFile,
	"service":                     tmplService,
	"runServices":                 tmplRunServices,
	"workspaceFolder":             tmplWorkspaceFolder,
	"workspaceMount":              tmplWorkspaceMount,
	"remoteUser":                  tmplRemoteUser,
	"containerUser":               tmplContainerUser,
	"updateRemoteUserUID":         tmplUpdateRemoteUserUID,
	"userEnvProbe":                tmplUserEnvProbe,
	"containerEnv":                tmplContainerEnv,
	"remoteEnv":                   tmplRemoteEnv,
	"localEnv":                    tmplLocalEnv,
	"appPort":                     tmplAppPort,
	"forwardPorts":                tmplForwardPorts,
	"portsAttributes":             tmplPortsAttributes,
	"otherPortsAttributes":        tmplOtherPortsAttributes,
	"mounts":                      tmplMounts,
	"runArgs":                     tmplRunArgs,
	"startupCommand":              tmplStartupCommand,
	"overrideCommand":             tmplOverrideCommand,
	"command":                     tmplCommand,
	"entrypoint":                  tmplEntrypoint,
	"init":                        tmplInit,
	"privileged":                  tmplPrivileged,
	"capAdd":                      tmplCapAdd,
	"capDrop":                     tmplCapDrop,
	"securityOpt":                 tmplSecurityOpt,
	"devices":                     tmplDevices,
	"hostRequirements":            tmplHostRequirements,
	"features":                    tmplFeatures,
	"overrideFeatureInstallOrder": tmplOverrideFeatureInstallOrder,
	"initializeCommand":           tmplInitializeCommand,
	"updateContentCommand":        tmplUpdateContentCommand,
	"waitFor":                     tmplWaitFor,
	"onCreateCommand":             tmplOnCreateCommand,
	"postCreateCommand":           tmplPostCreateCommand,
	"postStartCommand":            tmplPostStartCommand,
	"postAttachCommand":           tmplPostAttachCommand,
	"watch":                       tmplWatch,
	"customizations":              tmplCustomizations,
	"secrets":                     tmplSecrets,
	"shutdownAction":              tmplShutdownAction,
}

// Template returns the YAML snippet for a key, or a minimal fallback if no template is defined.
func Template(key string) string {
	if t, ok := templates[key]; ok {
		return t
	}
	return key + ": \n"
}
