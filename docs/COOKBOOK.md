# Cookbook

Worked configs for setups that come up often. Each one is a complete
`config.yaml` — drop it in, run `devcontainerwizard convert`, and reopen the
folder in your IDE.

---

## A Go service with Postgres

Two containers, so Compose. The IDE attaches to `app`; `db` comes up alongside it.

```yaml
name: go-api
dockerComposeFile:
  - docker-compose.yml
service: app
runServices:
  - db
workspaceFolder: /workspaces/go-api
remoteUser: vscode

forwardPorts:
  - 8080
  - 5432
portsAttributes:
  "8080":
    label: API
    onAutoForward: openBrowser
  "5432":
    label: PostgreSQL
    onAutoForward: silent

postCreateCommand: go mod download
customizations:
  vscode:
    extensions:
      - golang.go
    settings:
      go.lintTool: golangci-lint
```

---

## Reuse the host's Docker daemon

The `docker-outside-of-docker` feature mounts the host socket, so builds inside
the container hit the daemon you already have — no nested daemon, no layer cache
to rebuild.

```yaml
name: build-tools
image: mcr.microsoft.com/devcontainers/base:ubuntu
features:
  ghcr.io/devcontainers/features/docker-outside-of-docker:1:
    version: latest
```

Use `docker-in-docker:2` instead when the container must not see the host's
images and containers.

---

## Keep the module cache between rebuilds

A named volume survives rebuilds; a bind mount would drag host permissions
along with it.

```yaml
name: cached-go
image: mcr.microsoft.com/devcontainers/go:latest
mounts:
  - type: volume
    source: go-mod-cache
    target: /go/pkg/mod
  - type: volume
    source: go-build-cache
    target: /home/vscode/.cache/go-build
containerEnv:
  GOMODCACHE: /go/pkg/mod
  GOCACHE: /home/vscode/.cache/go-build
```

---

## Pass a host token into the container

`secrets` makes the IDE prompt for the value; `remoteEnv` with `${localEnv:...}`
takes it from your shell without prompting. Neither one puts the value in the
file.

```yaml
name: registry-client
image: mcr.microsoft.com/devcontainers/base:ubuntu

secrets:
  NPM_TOKEN:
    description: Token for the private npm registry
    documentationUrl: https://docs.npmjs.com/creating-and-viewing-access-tokens

remoteEnv:
  GITHUB_TOKEN: "${localEnv:GITHUB_TOKEN}"
```

---

## Debug a process that needs ptrace

`SYS_PTRACE` plus an unconfined seccomp profile is what delve, gdb, and lldb
need to attach.

```yaml
name: debuggable
image: mcr.microsoft.com/devcontainers/go:latest
capAdd:
  - SYS_PTRACE
securityOpt:
  - seccomp=unconfined
```

---

## Multi-stage Dockerfile, stopping at the dev stage

```yaml
name: staged-build
build:
  dockerfile: Dockerfile
  context: .
  target: development
  args:
    GO_VERSION: "1.25"
  cacheFrom:
    - ghcr.io/myorg/myapp:cache
```

---

## Restart a process when files change

```yaml
name: hot-reload
image: mcr.microsoft.com/devcontainers/go:latest
postStartCommand: go run ./cmd/server
watch:
  waitFor:
    - postCreateCommand
  restart:
    - "**/*.go"
    - "**/*.templ"
```

---

## Require a machine that can actually run it

Codespaces reads `hostRequirements` when picking a machine type; locally it is a
pre-flight check rather than a hard limit.

```yaml
name: heavy
image: mcr.microsoft.com/devcontainers/base:ubuntu
hostRequirements:
  cpus: 8
  memory: 16gb
  storage: 64gb
  gpu: optional
```

---

## Run everything as root

Occasionally simplest, in a throwaway container.

```yaml
name: root-shell
image: ubuntu:24.04
remoteUser: root
containerUser: root
updateRemoteUserUID: false
overrideCommand: true
```

---

## Derive a second config from the first

`edit --output` loads one file and saves to another, leaving the original alone.

```bash
devcontainerwizard edit --config config.yaml --output config.ci.yaml
```
