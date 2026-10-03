# Configuration Reference

`config.yaml` is the source file that `devcontainerwizard convert` translates
into `devcontainer.json`. It follows the
[Dev Containers specification](https://containers.dev/implementors/spec/), one
YAML key per specification property.

For the full field-by-field reference — types, defaults, allowed values — see
the generated [attribute pages](/devcontainerwizard/attributes/devcontainer.md),
or press **h** on any field in `devcontainerwizard edit`.

---

## Editor completion

A JSON Schema for `config.yaml` is generated alongside the docs. Reference it
from the top of your file and your editor will complete keys and flag invalid
values as you type:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/lucasassuncao/devcontainerwizard/main/docs/devcontainerwizard/schema/devcontainer.schema.json
```

Regenerate it locally with `make docs`.

---

## Choosing a container source

Exactly one of `image`, `build`, or `dockerComposeFile` must be present. The
editor and `convert` both reject a config that sets more than one.

### Image-based container

```yaml
name: my-project
image: mcr.microsoft.com/devcontainers/base:ubuntu
remoteUser: vscode
```

### Dockerfile-based container

```yaml
name: my-project
build:
  dockerfile: Dockerfile
  context: .
  args:
    NODE_VERSION: "20"
```

### Docker Compose

`service` names the container the IDE attaches to, and is required whenever
`dockerComposeFile` is set. `runServices` starts additional services alongside it.

```yaml
name: my-project
dockerComposeFile: docker-compose.yml
service: app
runServices:
  - db
workspaceFolder: /workspaces/my-project
```

---

## Environment variables

### containerEnv

Set for every process in the container, baked in at creation time.

```yaml
containerEnv:
  NODE_ENV: development
  API_URL: http://localhost:3000
```

### remoteEnv

Set for remote connections (the IDE, SSH) and resolved at attach time, so it can
reference other variables. `${localEnv:VAR}` reads from your host environment
and `${containerEnv:VAR}` from the container's:

```yaml
remoteEnv:
  PATH: "${containerEnv:PATH}:/custom/bin"
  MY_TOKEN: "${localEnv:MY_TOKEN}"
```

`${localEnv:...}` is the supported way to pull a value off the host — there is
no top-level `localEnv` block.

---

## Port forwarding

```yaml
forwardPorts:
  - 3000
  - 5432

portsAttributes:
  "3000":
    label: Web Server
    onAutoForward: openBrowser
    protocol: http
  "5432":
    label: PostgreSQL
    onAutoForward: silent

otherPortsAttributes:
  onAutoForward: silent
```

Port keys must be quoted — YAML would otherwise read `3000` as an integer.
`otherPortsAttributes` applies to every forwarded port `portsAttributes` does
not match.

---

## Features

Pre-built additions from the
[Dev Containers Features catalog](https://containers.dev/features). An empty
map means "defaults".

```yaml
features:
  ghcr.io/devcontainers/features/git:1: {}
  ghcr.io/devcontainers/features/github-cli:1:
    version: latest
  ghcr.io/devcontainers/features/docker-in-docker:2:
    moby: true

overrideFeatureInstallOrder:
  - ghcr.io/devcontainers/features/git:1
  - ghcr.io/devcontainers/features/docker-in-docker:2
```

---

## Lifecycle commands

Each command accepts a string, a list of strings, or a map of named commands
that run in parallel.

```yaml
initializeCommand: ./scripts/prepare-host.sh
onCreateCommand: echo "Container created"
updateContentCommand: npm install
postCreateCommand:
  - npm install
  - npm run build
postStartCommand: npm run dev
postAttachCommand: echo "Welcome!"

waitFor: postCreateCommand
```

`waitFor` names the stage the IDE waits for before reporting the container ready.

---

## Mounts

Each entry is either a structured object or a raw Docker `--mount` string.

```yaml
mounts:
  - type: bind
    source: ./data
    target: /data
  - type: volume
    source: node_modules
    target: /workspace/node_modules
```

The object form carries exactly `type`, `source`, and `target` — that is the
whole of the specification's Mount definition, and it declares
`additionalProperties: false`, so anything else is invalid. `type` and `target`
are required.

For options the object form has no property for — `readonly`, `consistency`,
`bind-propagation` — use the string form, which is passed to Docker unchanged:

```yaml
mounts:
  - source=./configs,target=/etc/app/config,type=bind,readonly
```

---

## Host requirements

```yaml
hostRequirements:
  cpus: 4
  memory: 8gb
  storage: 32gb
  gpu: true
```

`cpus` must be at least 1. `memory` and `storage` are a whole number with an
optional **lowercase** unit — `kb`, `mb`, `gb`, or `tb` — so `8gb` is valid and
`8GB` is not. A bare number means bytes.

`gpu` accepts `true`, `false`, the string `optional`, or an object with `cores`
(at least 1) and `memory`.

---

## IDE customizations

```yaml
customizations:
  vscode:
    extensions:
      - dbaeumer.vscode-eslint
      - esbenp.prettier-vscode
    settings:
      editor.formatOnSave: true
      editor.defaultFormatter: esbenp.prettier-vscode
  jetbrains:
    plugins:
      - org.rust.lang
```

---

## Secrets

Keyed by environment variable name. The IDE prompts the user for each value —
nothing secret is stored in the file.

```yaml
secrets:
  GITHUB_TOKEN:
    description: Personal access token for the registry
    documentationUrl: https://github.com/settings/tokens
```

Keys follow environment-variable naming: letters, digits and underscore, not
starting with a digit. Both sub-fields are optional, but a secret with no
`description` leaves the user staring at a bare variable name, so the editor
pre-checks it.

---

## File watching

```yaml
watch:
  waitFor:
    - postCreateCommand
  restart:
    - "**/*.go"
```

---

## Users and permissions

```yaml
remoteUser: vscode
containerUser: vscode
updateRemoteUserUID: true
userEnvProbe: loginInteractiveShell
```

---

## Runtime and security

```yaml
runArgs:
  - --network=host
init: true
privileged: false
capAdd:
  - SYS_PTRACE
securityOpt:
  - seccomp=unconfined
devices:
  - /dev/video0:/dev/video0
overrideCommand: true
shutdownAction: stopContainer
```

`overrideCommand` defaults to `true` for image and Dockerfile configs and
`false` for Docker Compose. `shutdownAction` accepts `none`, `stopContainer`,
or `stopCompose` — use `stopCompose` with a Compose config, where it is also
the default.
