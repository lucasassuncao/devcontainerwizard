# DevContainer Examples

## Preset: devops

```yaml
# DevContainer Configuration for DevOps and SRE workflows

name: devops-devcontainer

image: mcr.microsoft.com/devcontainers/base:ubuntu-24.04

workspaceFolder: /workspaces/${localWorkspaceFolderBasename}
remoteUser: vscode

containerEnv:
  KUBECONFIG: /home/vscode/.kube/config
  DOCKER_BUILDKIT: "1"

features:
  ghcr.io/devcontainers/features/git:1:
    version: latest
  ghcr.io/devcontainers/features/docker-outside-of-docker:1:
    version: latest
  ghcr.io/devcontainers/features/kubectl-helm-minikube:1:
    version: latest
    helm: latest
    minikube: none
  ghcr.io/devcontainers/features/terraform:1:
    version: latest
    tflint: latest
  ghcr.io/devcontainers/features/aws-cli:1:
    version: latest
  ghcr.io/devcontainers/features/azure-cli:1:
    version: latest
  ghcr.io/dhoeric/features/google-cloud-cli:1:
    version: latest
  ghcr.io/devcontainers/features/powershell:1:
    version: latest
  ghcr.io/sliekens/devcontainer-features/ansible:1: {}
  ghcr.io/dhoeric/features/trivy:1: {}

forwardPorts:
  - 8080
  - 9090

customizations:
  vscode:
    extensions:
      - hashicorp.terraform
      - ms-azuretools.vscode-docker
      - ms-kubernetes-tools.vscode-kubernetes-tools
      - redhat.vscode-yaml
      - amazonwebservices.aws-toolkit-vscode
      - ms-azuretools.azure-dev
      - googlecloudtools.cloudcode
      - github.vscode-github-actions
      - ms-toolsai.jupyter
      - usernamehw.errorlens

    settings:
      editor.formatOnSave: true

      "[terraform]":
        editor.defaultFormatter: hashicorp.terraform
        editor.formatOnSave: true
        editor.formatOnSaveMode: file

      terraform.experimentalFeatures.validateOnSave: true
      terraform.experimentalFeatures.prefillRequiredFields: true

      yaml.schemas:
        https://json.schemastore.org/github-workflow.json: ".github/workflows/*.yml"
        https://json.schemastore.org/kustomization.json: "kustomization.yaml"

      files.associations:
        "*.tf": terraform
        "*.tfvars": terraform
        "Dockerfile*": dockerfile

      files.trimTrailingWhitespace: true
      files.insertFinalNewline: true
```

## Preset: dockercompose

```yaml
# DevContainer Configuration
# Configuration using Docker Compose

name: my-devcontainer

# Docker Compose files
dockerComposeFile:
  - docker-compose.yml

# Service to use as dev container
service: app

# Workspace configuration
workspaceFolder: /workspace
remoteUser: vscode

# Container environment variables
containerEnv:
  NODE_ENV: development

# Port forwarding
forwardPorts:
  - 3000
  - 5432

# Features
features:
  ghcr.io/devcontainers/features/git:1: {}
  ghcr.io/devcontainers/features/github-cli:1: {}
  ghcr.io/devcontainers/features/docker-in-docker:2: {}

# Lifecycle hooks
postCreateCommand: npm install
postStartCommand: npm run dev

# VS Code customizations
customizations:
  vscode:
    extensions:
      - dbaeumer.vscode-eslint
      - esbenp.prettier-vscode
      - ms-azuretools.vscode-docker
    settings:
      editor.formatOnSave: true
```

## Preset: dockerfile

```yaml
# DevContainer Configuration
# Configuration using a custom Dockerfile

name: my-devcontainer

# Build from Dockerfile
build:
  dockerfile: Dockerfile
  context: .
  args:
    VARIANT: "20"

# Workspace configuration
workspaceFolder: /workspace
remoteUser: vscode

# Container environment variables
containerEnv:
  NODE_ENV: development
  API_URL: http://localhost:3000

# Port forwarding
forwardPorts:
  - 3000
  - 5432

portsAttributes:
  "3000":
    label: "Application"
    onAutoForward: openBrowser
  "5432":
    label: "Database"
    onAutoForward: silent

# Mounts
mounts:
  - type: bind
    source: ./data
    target: /data

# Features
features:
  ghcr.io/devcontainers/features/git:1: {}
  ghcr.io/devcontainers/features/github-cli:1: {}

# Lifecycle hooks
postCreateCommand: npm install
postStartCommand: npm run dev

# VS Code customizations
customizations:
  vscode:
    extensions:
      - dbaeumer.vscode-eslint
      - esbenp.prettier-vscode
      - ms-azuretools.vscode-docker
    settings:
      editor.formatOnSave: true
      editor.defaultFormatter: esbenp.prettier-vscode
```

## Preset: full

```yaml
# DevContainer Configuration - FULL REFERENCE
# This template shows ALL available configuration options
# Comment/uncomment sections as needed for your project

# ==============================================================================
# BASIC CONFIGURATION
# ==============================================================================

# Container name (REQUIRED)
name: full-devcontainer-example

# ==============================================================================
# CONTAINER SOURCE - Choose ONE of the following:
# ==============================================================================

# Option 1: Use existing Docker image
image: mcr.microsoft.com/devcontainers/typescript-node:20

# Option 2: Build from Dockerfile (comment 'image' above if using this)
# build:
#   # Dockerfile path (REQUIRED if using build)
#   dockerfile: Dockerfile
#
#   # Build context (REQUIRED if using build)
#   context: .
#
#   # Build arguments
#   args:
#     VARIANT: "20"
#     NODE_VERSION: "20.10.0"
#     INSTALL_ZSH: "true"
#
#   # Target stage for multi-stage builds
#   target: development
#
#   # Cache from images
#   cacheFrom:
#     - node:20-alpine
#     - myregistry.com/myimage:cache
#

# Option 3: Docker Compose (comment 'image' above if using this)
# dockerComposeFile:
#   - docker-compose.yml
#   - docker-compose.dev.yml
#
# # Service to use as dev container (REQUIRED with dockerComposeFile)
# service: app
#
# # Extra services to start alongside it. Defaults to every service in the file.
# runServices:
#   - db
#   - redis
#
# # workspaceFolder is REQUIRED with dockerComposeFile - a Compose config has no
# # automatic source mount to fall back on. Set it to the target of a volume
# # mount in docker-compose.yml.
#
# # With Compose, shutdownAction takes stopCompose instead of stopContainer.
# shutdownAction: stopCompose

# ==============================================================================
# WORKSPACE CONFIGURATION
# ==============================================================================

# Workspace folder path inside container
workspaceFolder: /workspace

# Workspace mount configuration (advanced)
# workspaceMount: source=${localWorkspaceFolder},target=/workspace,type=bind,consistency=cached

# Remote user (default user in container)
remoteUser: vscode

# User the container itself starts as, entrypoint included.
# Defaults to the user baked into the image.
containerUser: vscode

# On Linux, sync remoteUser's UID/GID with the local user so bind-mounted
# files keep the right ownership. Default: true
updateRemoteUserUID: true

# How to probe for user environment variables
# Options: none | loginShell | loginInteractiveShell | interactiveShell
userEnvProbe: loginInteractiveShell

# ==============================================================================
# ENVIRONMENT VARIABLES
# ==============================================================================

# Container environment variables (available inside container)
containerEnv:
  NODE_ENV: development
  API_URL: http://localhost:3000
  DATABASE_URL: postgresql://postgres:postgres@db:5432/myapp
  REDIS_URL: redis://redis:6379
  LOG_LEVEL: debug
  CONTAINER_WORKSPACE_FOLDER: ${containerWorkspaceFolder}

# Remote environment variables (for SSH and remote connections)
remoteEnv:
  PATH: "${containerEnv:PATH}:/usr/local/go/bin:${containerWorkspaceFolder}/bin"
  DISPLAY: "${localEnv:DISPLAY}"

# Local environment variables (not exported to JSON, used for processing only)
# Resolved from the host OS at conversion time; reference them in containerEnv
# or remoteEnv as ${localEnv:KEY}. Never written to devcontainer.json.
# localEnv:
#   MY_LOCAL_SECRET: ${env:MY_LOCAL_SECRET}
#   HOST_PROJECT_PATH: ${env:HOME}

# ==============================================================================
# PORT FORWARDING
# ==============================================================================

# Ports to forward from container to host
forwardPorts:
  - 3000 # Application
  - 5432 # PostgreSQL
  - 6379 # Redis
  - 8080 # Alternative web server
  - 9229 # Node.js debugger

# Port-specific attributes
portsAttributes:
  "3000":
    label: "Application"
    onAutoForward: openBrowser # Options: notify | openBrowser | openBrowserOnce | openPreview | silent | ignore
    protocol: http # Options: http | https

  "5432":
    label: "PostgreSQL Database"
    onAutoForward: silent

  "6379":
    label: "Redis Cache"
    onAutoForward: silent

  "8080":
    label: "Admin Panel"
    onAutoForward: notify
    protocol: http

  "9229":
    label: "Node Debugger"
    onAutoForward: silent

# Default attributes for ports not specified above
otherPortsAttributes:
  onAutoForward: ignore

# ==============================================================================
# MOUNTS & VOLUMES
# ==============================================================================

mounts:
  # Bind mount - links host directory to container
  - type: bind
    source: ./data
    target: /data

  # Read-only bind mount. The Mount object form has no "readonly" property in
  # the devcontainer specification, so read-only mounts use the raw Docker
  # --mount string form, which mounts accepts alongside objects.
  - source=./configs,target=/etc/app/config,type=bind,readonly

  # Named volume - persists data between container rebuilds
  - type: volume
    source: node_modules
    target: ${containerWorkspaceFolder}/node_modules

  # Another volume for PostgreSQL data
  - type: volume
    source: postgres-data
    target: /var/lib/postgresql/data

# ==============================================================================
# DOCKER RUN OPTIONS
# ==============================================================================

# Additional arguments for 'docker run'
runArgs:
  - --add-host=host.docker.internal:host-gateway
  - --network=my-custom-network

# Legacy port publishing, like 'docker run -p'. Unlike forwardPorts it needs a
# rebuild to change, so prefer forwardPorts above.
# appPort:
#   - 3000
#   - "8000:8010"

# Replace the image's default command with a sleep so the container stays up
# for the IDE to attach to.
# Default: true for image/Dockerfile, false for Docker Compose
overrideCommand: true

# ==============================================================================
# CONTAINER CAPABILITIES & SECURITY
# ==============================================================================

# Run init process inside container
init: true

# Run container in privileged mode (be careful!)
privileged: false

# Linux capabilities to add
capAdd:
  - SYS_PTRACE # For debugging
  - NET_ADMIN # For network configuration

# Security options
securityOpt:
  - seccomp=unconfined
  - apparmor=unconfined

# Devices to expose to container
devices:
  - /dev/fuse:/dev/fuse
  - /dev/net/tun:/dev/net/tun

# ==============================================================================
# HOST REQUIREMENTS
# ==============================================================================

# Minimum host hardware. Codespaces picks a machine type from this; locally it
# is a pre-flight check.
hostRequirements:
  # At least 1.
  cpus: 4
  # Whole number with an optional LOWERCASE unit: kb, mb, gb, tb.
  # "8GB" is rejected by the specification; "8gb" is correct.
  memory: 8gb
  storage: 32gb
  # true | false | "optional" | an object with cores and memory
  # gpu: optional

# ==============================================================================
# FEATURES
# ==============================================================================

# Order to install features (override default)
overrideFeatureInstallOrder:
  - ghcr.io/devcontainers/features/common-utils
  - ghcr.io/devcontainers/features/git
  - ghcr.io/devcontainers/features/node

# Features to install (from https://containers.dev/features)
features:
  # Common utilities
  ghcr.io/devcontainers/features/common-utils:2:
    installZsh: true
    installOhMyZsh: true
    upgradePackages: true
    username: vscode

  # Git
  ghcr.io/devcontainers/features/git:1:
    version: latest
    ppa: false

  # GitHub CLI
  ghcr.io/devcontainers/features/github-cli:1:
    version: latest

  # Docker-in-Docker
  ghcr.io/devcontainers/features/docker-in-docker:2:
    version: latest
    moby: true
    dockerDashComposeVersion: v2

  # Node.js (if not in base image)
  # ghcr.io/devcontainers/features/node:1:
  #   version: "20"
  #   nodeGypDependencies: true
  #   nvmVersion: "latest"

  # Python
  # ghcr.io/devcontainers/features/python:1:
  #   version: "3.11"
  #   installTools: true

  # Go
  # ghcr.io/devcontainers/features/go:1:
  #   version: "1.21"

  # AWS CLI
  # ghcr.io/devcontainers/features/aws-cli:1:
  #   version: latest

  # Azure CLI
  # ghcr.io/devcontainers/features/azure-cli:1:
  #   version: latest

  # kubectl
  # ghcr.io/devcontainers/features/kubectl-helm-minikube:1:
  #   version: latest
  #   helm: latest
  #   minikube: latest

  # Terraform
  # ghcr.io/devcontainers/features/terraform:1:
  #   version: latest

# ==============================================================================
# LIFECYCLE HOOKS
# ==============================================================================

# Runs on the HOST before the container is created - the only lifecycle hook
# that does. Use it to prepare files the build needs.
# initializeCommand: ./scripts/prepare-host.sh

# Command runs after container is created (runs once)
onCreateCommand: |
  echo "🚀 Container created!"
  echo "Installing global tools..."
  npm install -g pnpm typescript tsx nodemon
  echo "Setting up git configuration..."
  git config --global init.defaultBranch main

# Runs after onCreate, and again whenever the workspace content is refreshed.
# This is the stage waitFor waits on by default.
updateContentCommand: pnpm install --frozen-lockfile

# Command runs after onCreate (runs once)
postCreateCommand: |
  echo "📦 Installing dependencies..."
  pnpm install
  echo "🔧 Running database migrations..."
  pnpm run db:migrate
  echo "🌱 Seeding database..."
  pnpm run db:seed
  echo "✅ Setup complete!"

# Command runs every time container starts
postStartCommand: |
  echo "🟢 Container started!"
  echo "Starting development server..."
  pnpm run dev

# Command runs after attaching to container
postAttachCommand: |
  echo ""
  echo "╔════════════════════════════════════════╗"
  echo "║   Welcome to DevContainer! 🎉         ║"
  echo "╚════════════════════════════════════════╝"
  echo ""
  echo "📝 Available commands:"
  echo "  - pnpm run dev      : Start dev server"
  echo "  - pnpm run test     : Run tests"
  echo "  - pnpm run build    : Build project"
  echo ""
  echo "🔗 Services:"
  echo "  - App: http://localhost:3000"
  echo "  - DB:  postgresql://localhost:5432"
  echo ""

# ==============================================================================
# FILE WATCHING
# ==============================================================================

watch:
  # Wait for these processes/files before starting
  waitFor:
    - postCreateCommand
    - postStartCommand

  # Restart container when these change
  restart:
    - .devcontainer/devcontainer.json
    - .devcontainer/Dockerfile
    - docker-compose.yml
    - package.json

# ==============================================================================
# CUSTOMIZATIONS
# ==============================================================================

customizations:
  # ============================================================================
  # VS CODE
  # ============================================================================
  vscode:
    # Extensions to install
    extensions:
      # ---- Language Support ----
      - dbaeumer.vscode-eslint
      - esbenp.prettier-vscode
      - ms-python.python
      - ms-python.vscode-pylance
      - golang.go
      - rust-lang.rust-analyzer

      # ---- Git ----
      - eamodio.gitlens
      - mhutchie.git-graph
      - donjayamanne.githistory

      # ---- Docker ----
      - ms-azuretools.vscode-docker
      - ms-kubernetes-tools.vscode-kubernetes-tools

      # ---- Database ----
      - mtxr.sqltools
      - mtxr.sqltools-driver-pg
      - cweijan.vscode-mysql-client2

      # ---- Testing ----
      - orta.vscode-jest
      - hbenl.vscode-test-explorer
      - kavod-io.vscode-jest-test-adapter

      # ---- API Development ----
      - humao.rest-client
      - 42crunch.vscode-openapi

      # ---- Utilities ----
      - usernamehw.errorlens
      - christian-kohler.path-intellisense
      - visualstudioexptteam.vscodeintellicode
      - wayou.vscode-todo-highlight
      - gruntfuggly.todo-tree
      - streetsidesoftware.code-spell-checker
      - editorconfig.editorconfig

      # ---- Productivity ----
      - github.copilot
      - github.copilot-chat
      - tabnine.tabnine-vscode

      # ---- Theme & UI ----
      - pkief.material-icon-theme
      - github.github-vscode-theme

    # VS Code settings
    settings:
      # ---- Editor ----
      editor.formatOnSave: true
      editor.defaultFormatter: esbenp.prettier-vscode
      editor.codeActionsOnSave:
        source.fixAll.eslint: "explicit"
        source.organizeImports: "explicit"
      editor.tabSize: 2
      editor.insertSpaces: true
      editor.rulers: [80, 120]
      editor.wordWrap: "on"
      editor.minimap.enabled: true
      editor.suggestSelection: "first"
      editor.bracketPairColorization.enabled: true
      editor.guides.bracketPairs: true
      editor.inlineSuggest.enabled: true

      # ---- Files ----
      files.autoSave: "afterDelay"
      files.autoSaveDelay: 1000
      files.trimTrailingWhitespace: true
      files.insertFinalNewline: true
      files.exclude:
        "**/.git": true
        "**/node_modules": true
        "**/.next": true
        "**/dist": true
        "**/build": true
        "**/.turbo": true
        "**/__pycache__": true
        "**/.pytest_cache": true

      # ---- Search ----
      search.exclude:
        "**/node_modules": true
        "**/dist": true
        "**/build": true
        "**/.git": true
        "**/coverage": true

      # ---- Terminal ----
      terminal.integrated.defaultProfile.linux: "zsh"
      terminal.integrated.fontFamily: "MesloLGS NF"
      terminal.integrated.fontSize: 14
      terminal.integrated.cursorStyle: "line"
      terminal.integrated.cursorBlinking: true

      # ---- Git ----
      git.enableSmartCommit: true
      git.confirmSync: false
      git.autofetch: true
      gitlens.defaultDateFormat: "YYYY-MM-DD HH:mm"

      # ---- Language Specific ----
      "[javascript]":
        editor.defaultFormatter: "esbenp.prettier-vscode"
      "[typescript]":
        editor.defaultFormatter: "esbenp.prettier-vscode"
      "[json]":
        editor.defaultFormatter: "esbenp.prettier-vscode"
      "[jsonc]":
        editor.defaultFormatter: "esbenp.prettier-vscode"
      "[html]":
        editor.defaultFormatter: "esbenp.prettier-vscode"
      "[css]":
        editor.defaultFormatter: "esbenp.prettier-vscode"
      "[python]":
        editor.defaultFormatter: "ms-python.black-formatter"
      "[go]":
        editor.defaultFormatter: "golang.go"

      # ---- ESLint ----
      eslint.validate:
        - javascript
        - typescript
        - javascriptreact
        - typescriptreact

      # ---- Prettier ----
      prettier.singleQuote: true
      prettier.trailingComma: "es5"
      prettier.semi: true
      prettier.tabWidth: 2

      # ---- TypeScript ----
      typescript.updateImportsOnFileMove.enabled: "always"
      typescript.suggest.autoImports: true
      javascript.suggest.autoImports: true

      # ---- Python ----
      python.linting.enabled: true
      python.linting.pylintEnabled: true
      python.formatting.provider: "black"
      python.testing.pytestEnabled: true

      # ---- Go ----
      go.useLanguageServer: true
      go.lintTool: "golangci-lint"
      go.formatTool: "goimports"

      # ---- Other ----
      todo-tree.general.tags: ["TODO", "FIXME", "HACK", "NOTE"]
      errorLens.enabledDiagnosticLevels: ["error", "warning"]

    # Remote user override (if different from container remoteUser)
    # remoteUser: vscode

  # ============================================================================
  # GITHUB CODESPACES
  # ============================================================================
  codespaces:
    settings:
      editor.formatOnSave: true
      terminal.integrated.defaultProfile.linux: "zsh"

    extensions:
      - dbaeumer.vscode-eslint
      - esbenp.prettier-vscode
      - github.copilot

  # ============================================================================
  # JETBRAINS IDEs
  # ============================================================================
  jetbrains:
    plugins:
      - com.intellij.plugins.github
      - org.jetbrains.plugins.go
      - com.jetbrains.python


# ==============================================================================
# SECRETS
# ==============================================================================

secrets:
  # NPM token for private packages
  NPM_TOKEN:
    description: "NPM authentication token for private registry"

  # GitHub token
  GITHUB_TOKEN:
    description: "GitHub personal access token"

  # Database password
  DB_PASSWORD:
    description: "PostgreSQL database password"

  # API keys
  API_SECRET_KEY:
    description: "Secret key for API authentication"

# ==============================================================================
# SHUTDOWN ACTION
# ==============================================================================

# Action when container is stopped
# Options: none | stopContainer | stopCompose
shutdownAction: stopContainer
```

## Preset: golang

```yaml
# DevContainer Configuration for Go Development
# Optimized setup for professional Go development with best practices

name: golang-devcontainer

# Use official Go devcontainer image
image: "mcr.microsoft.com/devcontainers/go:latest"

# Workspace configuration
workspaceFolder: "/workspaces/${localWorkspaceFolderBasename}"
remoteUser: vscode

# Dev Container Features
features:
  # Essential version control
  ghcr.io/devcontainers/features/git:1:
    version: latest
  ghcr.io/devcontainers/features/go:1:
    golangciLintVersion: latest
    version: "1.25.4"

# VS Code customizations
customizations:
  vscode:
    extensions:
      # Go language support
      - golang.go

    settings:
      # Go settings
      go.useLanguageServer: true
      go.lintTool: golangci-lint
      go.lintOnSave: package
      go.formatTool: goimports
      go.testFlags:
        - -v
        - -race
      go.coverOnSave: true
      go.toolsManagement.autoUpdate: true
      
      # Delve debugger settings
      go.delveConfig:
        apiVersion: 2
        showGlobalVariables: true
      
      # Editor settings
      editor.formatOnSave: true
      editor.codeActionsOnSave:
        source.organizeImports: explicit
      
      # Files settings
      files.eol: "\n"
      files.trimTrailingWhitespace: true
      files.insertFinalNewline: true
```

## Preset: image

```yaml
# DevContainer Configuration
# This is a minimal configuration using an existing Docker image

name: my-devcontainer
image: mcr.microsoft.com/devcontainers/base:ubuntu

# Workspace configuration
workspaceFolder: /workspace
remoteUser: vscode

# Container environment variables
containerEnv:
  NODE_ENV: development

# Port forwarding
forwardPorts:
  - 3000

# VS Code customizations
customizations:
  vscode:
    extensions:
      - dbaeumer.vscode-eslint
      - esbenp.prettier-vscode
    settings:
      editor.formatOnSave: true
```

## Preset: python

```yaml
# DevContainer Configuration for Python Development

name: python-devcontainer

image: mcr.microsoft.com/devcontainers/python:3.12

workspaceFolder: /workspaces/${localWorkspaceFolderBasename}
remoteUser: vscode

containerEnv:
  PYTHONUNBUFFERED: "1"
  PYTHONDONTWRITEBYTECODE: "1"
  PIP_NO_CACHE_DIR: "1"

features:
  ghcr.io/devcontainers/features/git:1:
    version: latest
  ghcr.io/devcontainers/features/python:1:
    version: "3.12"
    installTools: true

forwardPorts:
  - 8000

postCreateCommand: pip install --upgrade pip && pip install -r requirements.txt 2>/dev/null || true

customizations:
  vscode:
    extensions:
      - ms-python.python
      - ms-python.vscode-pylance
      - ms-toolsai.jupyter
      - ms-python.black-formatter
      - ms-python.isort
      - ms-python.mypy-type-checker
      - charliermarsh.ruff
      - usernamehw.errorlens

    settings:
      editor.formatOnSave: true
      editor.defaultFormatter: ms-python.black-formatter

      python.defaultInterpreterPath: /usr/local/bin/python
      python.testing.pytestEnabled: true
      python.testing.unittestEnabled: false
      python.analysis.typeCheckingMode: basic

      "[python]":
        editor.defaultFormatter: ms-python.black-formatter
        editor.codeActionsOnSave:
          source.organizeImports: explicit

      isort.args:
        - --profile
        - black

      files.trimTrailingWhitespace: true
      files.insertFinalNewline: true
      files.eol: "\n"
```

## Preset: rust

```yaml
# DevContainer Configuration for Rust Development

name: rust-devcontainer

image: mcr.microsoft.com/devcontainers/rust:latest

workspaceFolder: /workspaces/${localWorkspaceFolderBasename}
remoteUser: vscode

containerEnv:
  RUST_BACKTRACE: "1"
  CARGO_HOME: /usr/local/cargo

features:
  ghcr.io/devcontainers/features/git:1:
    version: latest
  ghcr.io/devcontainers/features/rust:1:
    version: stable
    profile: default

forwardPorts:
  - 8080

postCreateCommand: rustup component add clippy rustfmt

customizations:
  vscode:
    extensions:
      - rust-lang.rust-analyzer
      - vadimcn.vscode-lldb
      - tamasfe.even-better-toml
      - fill-labs.dependi
      - usernamehw.errorlens

    settings:
      editor.formatOnSave: true
      editor.defaultFormatter: rust-lang.rust-analyzer

      rust-analyzer.checkOnSave.command: clippy
      rust-analyzer.inlayHints.enabled: true
      rust-analyzer.cargo.buildScripts.enable: true

      files.trimTrailingWhitespace: true
      files.insertFinalNewline: true
      files.eol: "\n"
```

