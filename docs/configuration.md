# Configuration Reference

`config.yaml` is the source file that `devcontainerwizard convert` translates into `devcontainer.json`. It follows the [Dev Containers specification](https://containers.dev/implementors/spec/) with a few additions described below.

For field-level documentation on each block, run `devcontainerwizard show-docs` or browse the [model docs](index.md).

---

## Image-based container

```yaml
name: my-project
image: mcr.microsoft.com/devcontainers/base:ubuntu
remoteUser: vscode
```

---

## Dockerfile-based container

```yaml
name: my-project
build:
  dockerfile: Dockerfile
  context: .
  args:
    NODE_VERSION: "20"
```

---

## Docker Compose

```yaml
name: my-project
dockerComposeFile: docker-compose.yml
service: app
workspaceFolder: /workspaces/my-project
```

---

## Environment variables

### containerEnv

Variables available inside the container.

```yaml
containerEnv:
  NODE_ENV: development
  API_URL: http://localhost:3000
```

### remoteEnv

Variables for remote connections (SSH, VS Code Remote).

```yaml
remoteEnv:
  PATH: "${containerEnv:PATH}:/custom/bin"
```

### localEnv

Variables resolved from the host environment before conversion. Not written to `devcontainer.json`.

```yaml
localEnv:
  MY_SECRET: ${env:MY_SECRET}
```

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
```

---

## Features

Add pre-built features from the [Dev Containers Features catalog](https://containers.dev/features).

```yaml
features:
  ghcr.io/devcontainers/features/git:1: {}
  ghcr.io/devcontainers/features/github-cli:1:
    version: latest
  ghcr.io/devcontainers/features/docker-in-docker:2:
    moby: true
```

---

## Lifecycle hooks

Values can be a string or a list of strings.

```yaml
onCreateCommand: echo "Container created"
postCreateCommand:
  - npm install
  - npm run build
postStartCommand: npm run dev
postAttachCommand: echo "Welcome!"
```

---

## Mounts

```yaml
mounts:
  - type: bind
    source: ./data
    target: /data
  - type: volume
    source: node_modules
    target: /workspace/node_modules
```

---

## VS Code customizations

```yaml
customizations:
  vscode:
    extensions:
      - dbaeumer.vscode-eslint
      - esbenp.prettier-vscode
    settings:
      editor.formatOnSave: true
      editor.defaultFormatter: esbenp.prettier-vscode
```
